#!/bin/bash

# to test:
# - python -m SimpleHTTPServer
# - open http://127.0.0.1:8000/csvlada.html

#ls -alh
export GOPATH=~/.gopath
go get -v
go build

SESSIONID=$1

if [ -z $SESSIONID ]; then
    echo "Please set SESSIONID as first arg (Copy PHPSESSID cookie from http://sklad.lada-direct.ru/v2/cars/vesta/sw-cross/prices.html). Exit..."
    exit 1
fi

COLLECT_TS=`date --iso-8601=seconds`

# Vesta
./ladacsv -session=$SESSIONID -model=vesta/sedan -out=vesta_sedan.csv
./ladacsv -session=$SESSIONID -model=vesta/cross -out=vesta_cross.csv
./ladacsv -session=$SESSIONID -model=vesta/sw -out=vesta_sw.csv
./ladacsv -session=$SESSIONID -model=vesta/sw-cross -out=vesta_sw_cross.csv
./ladacsv -session=$SESSIONID -model=vesta/cng -out=vesta_cng.csv
./ladacsv -session=$SESSIONID -model=vesta/sport -out=vesta_sport.csv

# Granta
./ladacsv -session=$SESSIONID -model=granta/sedan -out=granta_sedan.csv
./ladacsv -session=$SESSIONID -model=granta/liftback -out=granta_liftback.csv
./ladacsv -session=$SESSIONID -model=granta/hatchback -out=grant_hatchback.csv
./ladacsv -session=$SESSIONID -model=granta/universal -out=granta_universal.csv
./ladacsv -session=$SESSIONID -model=granta/cross -out=granta_cross.csv
./ladacsv -session=$SESSIONID -model=granta/school -out=granta_school.csv
./ladacsv -session=$SESSIONID -model=granta/drive-active -out=granta_drive_active.csv

# XRAY
./ladacsv -session=$SESSIONID -model=xray/hatchback -out=xray_hatchback.csv
./ladacsv -session=$SESSIONID -model=xray/cross -out=xray_cross.csv

# Largus
./ladacsv -session=$SESSIONID -model=largus/universal -out=largus_universal.csv
./ladacsv -session=$SESSIONID -model=largus/cross -out=largus_cross.csv
./ladacsv -session=$SESSIONID -model=largus/furgon -out=largus_furgon.csv
./ladacsv -session=$SESSIONID -model=largus/furgon-cng -out=largus_furgon_cng.csv

# Niva
./ladacsv -session=$SESSIONID -model=niva/niva -out=niva_standard.csv
./ladacsv -session=$SESSIONID -model=niva/off-road -out=niva_offroad.csv

# 4x4
./ladacsv -session=$SESSIONID -model=4x4/3dv -out=n4x4_3doors.csv
./ladacsv -session=$SESSIONID -model=4x4/5dv -out=n4x4_5doors.csv
./ladacsv -session=$SESSIONID -model=4x4/urban -out=n4x4_3doors_urban.csv
./ladacsv -session=$SESSIONID -model=4x4/urban5 -out=n4x4_5doors_urban.csv

rm -f csvlada.html
cp csvlada.html.template index.html

sed -i "s/UPDATE_TS/${COLLECT_TS}/g" index.html


