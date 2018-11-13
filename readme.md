# Golf

Deployed at [golf.jacoduplessis.co.za](https://golf.jacoduplessis.co.za).

## Data update interval

New leaderboard data is fetched every 60 seconds. News content
is fetched upon request.

Please do not scrape the news page. Use [the library](https://github.com/jacoduplessis/twitterparse)
instead.

## Installing

```
go get -u github.com/jacoduplessis/golf
```

## Building

```
go build
```

The binary is named `golf`.

