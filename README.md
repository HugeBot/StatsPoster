# HUGE - StatsPoster
Simple script for post server count from database to diferent discord list websites.

## Requirements
- A postgresql database.
- A schema called discord.
- A table called discord.shard_guilds.
- A [credentials.conf](/credentials.conf.example) file with all your credentials


### Database schema ans table
```sql
CREATE SCHEMA discord AUTHORIZATION your_database_user;

CREATE TABLE shard_stats (
	id int4 NOT NULL,
	guild_count int4 NOT NULL,
	CONSTRAINT shard_stats_pkey PRIMARY KEY (id)
);
```

### How to build and execute the script
##### With Golang
```shell
go run main
```
##### Windows
```shell
go build
./stats-poster.exe
```
##### Linux
```shell
go build
chmod +x main
./stats-poster
```

### Schedule with cronjob on Linux
```shell
crontab -e
```

```shell
@hourly cd /path/to/binary/main && ./main
```
or
```shell
0 * * * * cd /path/to/binary/main && ./main
```
