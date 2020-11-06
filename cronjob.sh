#!/bin/bash

# crontab -e
# 0 1 * * * /home/kafeg/dev/ladacsv/cronjob.sh

rm /home/kafeg/dev/ladacsv/cron.log

echo "$(date): Start collecting" > /home/kafeg/dev/ladacsv/cron.log 2>&1

cd /home/kafeg/dev/ladacsv/

chmod a+x ./collect.sh
chmod a+x ./push.sh

./collect.sh lol >> /home/kafeg/dev/ladacsv/cron.log 2>&1

sleep 5

./push.sh >> /home/kafeg/dev/ladacsv/cron.log

echo "$(date): Finish collecting" >> /home/kafeg/dev/ladacsv/cron.log 2>&1
