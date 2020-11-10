# ladacsv
Parser for all available cars from all dillers from the http://sklad.lada-direct.ru/

To setup auto-collect:
- cd ~/
- git clone https://github.com/kafeg/ladacsv.git
- cd ladacsv
- chmod a+x *.sh
- crontab -e
- - 0 1 * * * /home/kafeg/dev/ladacsv/cronjob.sh

So at the 01:00 AM every day collecting will starts. Collect time ~3-4 hours in very lite mode to not DDoS lada site, grabber makes requests one by one without parallelism

Currently I have cron job, which collect everything at night and publish all results every morning there: https://kafeg.github.io/ladacsv

PS: Project is just for experiments with GoLang and does not targeting to abuse data or like that.

PPS: All grabbed data from the public sources, available on the official Lada site.
