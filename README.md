# beermate
A Slack bot to return Beer Advocate search results

# Build

`go build .`

# Docker build

`docker build -t beermate:0.1 --build-arg TOKEN=<slack-bot-token>`

# Cross compile for linux

`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .`

# Run

`docker run beermate`
