# Golf

Deployed at [golf.jacoduplessis.co.za](https://golf.jacoduplessis.co.za).

## Data update interval

New leaderboard data is fetched every 60 seconds. News content
is fetched upon request.

Please do not scrape the news page. Use [the library](https://github.com/jacoduplessis/twitterparse)
instead.

The HTML tables on the index page is well suited for use with Google Sheets:

```
=IMPORTHTML("https://golf.jacoduplessis.co.za/", "table", 1)
```

Some folks are using this for fantasy golf competitions. If you are interested in collaborating
to build a fantasy golf platform, please get in touch.


## Installing

```
go get -u github.com/jacoduplessis/golf
```

## Building

```
go build cmd/golf/golf.go
```

The binary is named `golf`.

## Deploy

```
rsync -zz --progress ./golf host:path
ssh root@host systemctl restart golf.service 
```