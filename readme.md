## MySQL dumper and uploader to AWS S3

Simple command-line interface for MySQL dumps and uploading to AWS S3.  
Dump example: `2020-01-15_1580218829.sql` where `1580218829` is the unix timestamp.

#####How to run:
 
`cp config.example.json config.json`  
`go build`  


You can exec the command with "-config" flag: `./aws-dumper -config any.json`  
"config.json" is used by default if executed without flags

Run in console or using cron.

###Cron
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
  0 0 * * * /backup/aws-dumper

```


### Config fields
##### dumpsToKeep
"dumpsToKeep" is the amount of latest dumps which will be kept in AWS S3 Bucket.  
Let's say if you are running dumper once a day
by cronjob and "dumpsToKeep" is 10, then only 10 latest dumps will be kept and all older dump files will be deleted.

If the "dumpsToKeep" is 0, then no files will be deleted.  

