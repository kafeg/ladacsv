package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"sort"
	"io"
)

var currRegion = ""
var currCity = ""

var optPHPSESSID = "" //take it from browser from 'http://sklad.lada-direct.ru/' from the developer console
const optTargetCarModelUrlTemplate = "http://sklad.lada-direct.ru/v2/cars/MODELNAME/prices.html"
var optTargetCarModelUrl = "" //change to your target car model
var optOutputFileName = ""
var interruptHandled = false

func readLines(file string) (lines []string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		const delim = '\n'
		line, err := r.ReadString(delim)
		if err == nil || len(line) > 0 {
			if err != nil {
				line += string(delim)
			}
			lines = append(lines, line)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return lines, nil
}

func writeLines(file string, lines []string, header string) (err error) {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	_, err = w.WriteString(header)
	if err != nil {
		return err
	}

	for _, line := range lines {
		_, err := w.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func sortCSVFile(outputFile string) {
	lines, err := readLines(outputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	header := lines[0]
	lines[0] = ""
	sort.Strings(lines)
	fmt.Printf("%v", header)
	err = writeLines(outputFile, lines, header)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	modelCodePtr := flag.String("model", "", "Model code")
	outPtr := flag.String("out", "", "Output file name")
	phpSessId := flag.String("session", "", "PHP Session Id")

	flag.Parse()

	if (*modelCodePtr == "") {
		println("Wrong model code")
		os.Exit(1)
		return
	}

	if (*outPtr == "") {
		println("Wrong file name")
		os.Exit(1)
		return
	}

	if (*phpSessId == "") {
		println("Wrong PHP sesssion Id. Open https://www.lada.ru/ and copy Cookie from the browser console")
		os.Exit(1)
		return
	}

	optPHPSESSID = *phpSessId

	// signals
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		interruptHandled = true
	}()

	collectModelInfo(*modelCodePtr, *outPtr)
	sortCSVFile(*outPtr)
}

func collectModelInfo(modelCode, outputFile string) {
	optTargetCarModelUrl = strings.Replace(optTargetCarModelUrlTemplate, "MODELNAME", modelCode, 1)
	optOutputFileName = outputFile

	// remove target file
	os.Remove(optOutputFileName)

	// open output file
	fo, err := os.OpenFile(optOutputFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// write a header
	if _, err := fo.Write([]byte("Регион №, Регион, Город №, Город, Комплектация, Цвет\n")); err != nil {
		panic(err)
	}

	// close fo on
	if err := fo.Close(); err != nil {
		panic(err)
	}

	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.Get("http://sklad.lada-direct.ru/", g.Opt.ParseFunc)
		},
		ParseFunc: parseCities,
		RobotsTxtDisabled: true,
	}).Start()
}

func parseCities(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("script").Each(func(i int, s *goquery.Selection) {
        var lines = strings.Split(s.Text(), "\n")

        var citiesJson = ""
		for i := 0; i < len(lines); i++ {
			lines[i] = strings.Trim(lines[i], "\n ")

			if (strings.HasPrefix(lines[i], "var objCity = ")) {
				lines[i] = strings.Replace(lines[i], "var objCity = ", "", 1)
				lines[i] = strings.Replace(lines[i], "}}};", "}}}", 1)

				citiesJson = lines[i]
				break
			}
		}

		regions := gjson.Parse(citiesJson)

		regions.ForEach(func(regionKey, regionValue gjson.Result) bool {
			currRegion = regionKey.String() + ", " + regionValue.Get("name").String()

			cities := regionValue.Get("city")
			cities.ForEach(func(cityKey, cityValue gjson.Result) bool {
				currCity = cityKey.String() + ", " + cityValue.String()

				geziyor.NewGeziyor(&geziyor.Options {
					StartRequestsFunc: func(g *geziyor.Geziyor) {
						req, err := client.NewRequest("GET", optTargetCarModelUrl, nil)
						if err != nil {
							log.Printf("Request creating error %v\n", err)
							return
						}

						var cookies = make(map[string]string)
						cookies["cookie_city"] = cityKey.String()
						cookies["crid"] = regionKey.String()
						cookies["ccid"] = cityKey.String()
						cookies["PHPSESSID"] = optPHPSESSID //open browser and copy from the developer console

						var cookiesStr = ""
						for key, value := range cookies {
							cookiesStr += key + "=" + value + ";"
						}

						req.Header.Set("Cookie", cookiesStr)
						req.Header.Set("User-Agent", "Yandex")

						g.Do(req, g.Opt.ParseFunc)
					},
					ParseFunc: parsePrices,
					RobotsTxtDisabled: true,
				}).Start()

				return !interruptHandled // keep iterating
			})

			return !interruptHandled // keep iterating
		})
	})
}

func parsePrices(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("div.kompl").Each(func(i int, s *goquery.Selection) {
		var komplName = strings.Replace(s.Find("p.kompl_name").Text(), "\n", " ", 10)
		komplName = strings.Replace(komplName, "  ", " ", 100)
		komplName = strings.Replace(komplName, "  ", " ", 100)
		komplName = strings.Replace(komplName, "  ", " ", 100)
		komplName = strings.Replace(komplName, "  ", " ", 100)
		komplName = strings.Replace(komplName, "  ", " ", 100)
		komplName = strings.Replace(komplName, "  ", " ", 100)
		komplName = strings.Replace(komplName, ",", "", 100)

		// open output file
		fo, err := os.OpenFile(optOutputFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		// close fo on exit and check for its returned error
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()

		s.Find("p.has_dealer span.color_dealer").Each(func(i int, s *goquery.Selection) {
			item, _ := s.Attr("title")

			var line = currRegion + ", " + currCity + ", " + komplName + ", " + item
			println(line)

			// write a data line
			if _, err := fo.Write([]byte(line + "\n")); err != nil {
				panic(err)
			}
		})
	})
}