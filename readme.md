# MySQL dumper and uploader to AWS S3

[![CircleCI](https://circleci.com/gh/NikolayOskin/sql-dumper.svg?style=svg)](https://circleci.com/gh/NikolayOskin/sql-dumper)


Simple command-line interface for MySQL dumps and uploading to AWS S3. Run it in console or using cron.  
Dump file name example: `2020-01-15_1580218829.sql` where `1580218829` is the unix timestamp.

#### How to run:  
 
```
cp config.example.json config.json  
go build  
```


You can exec the command with "-config" flag: `./aws-dumper -config=/path/to/config.json`  
If executed without flags then "./config.json" will be used by default.  

The program uses `mysqldump` utility to dump your database so it must be installed. It's probably already installed if you use mysql on your server(machine).

## Cron  
Because dumper executes system command (mysqldump), you may run into issues where crontab can't find it.  
The issue can be resolved by adding path variables to the `crontab -e` file.  
```
PATH=$PATH:/usr/local/bin:/usr/bin:/bin

# once a day at 00:00

# +---------------- minute (0 - 59)
# | +------------- hour (0 - 23)
# | | +---------- day of month (1 - 31)
# | | | +------- month (1 - 12)
# | | | | +---- day of week (0 - 6) (Sunday=0)
# | | | | |
  0 0 * * * /backup/aws-dumper -config=/path/to/the/config.json >> /dev/null 2>&1
```
You must set your "-config" flag with the absolute path to your config file, 
otherwise you will get `open ./config.json: no such file or directory` error.  

Instead of `/dev/null` you can set your log file location.
```
0 0 * * * /backup/aws-dumper -config=/path/to/the/config.json >> /var/log/aws-dumper.log 2>&1
```

## Config
#### dumpsToKeep  
"dumpsToKeep" is the amount of latest dumps which will be kept in AWS S3 Bucket.  
Let's say if you are running dumper once a day
by cronjob and "dumpsToKeep" is 10, then only 10 latest dumps will be kept and all older dump files will be deleted each time the program is executed.

If the "dumpsToKeep" is 0, then no files will be deleted.  

#### aws_region  
"aws_region" is the region code of AWS services such as S3. You can find the list of all region names and codes here: https://docs.aws.amazon.com/general/latest/gr/rande.html