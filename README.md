# ladacsv
Parser of official Lada site to get all available cars

Clone repository to the local machine

Open http://sklad.lada-direct.ru/v2/cars/vesta/sw-cross/prices.html and copy cokkie PHPSESSID from the browser inspector

Run the following command on Ubuntu >= 18.04 with installed Golang:
```./collect.sh <COPIED VALUE OF PHPSESSID>```

Wait for collecting all models, a looong time ;-)

Commit with ```git commit -am "Updated models `date --iso-8601=seconds`"```

Check results on the https://kafeg.github.io/ladacsv
