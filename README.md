# QLogger

A simple logging client.

v 1.0.1

## To set up locally

Install go and git, then install deps with ```go mod download```

Create .env file according to the example

Run ```go run cmd/main.go``` for development.

!! Important! make sure you're running from the root of the project, not the cmd folder.

For production, first run ```go build -o out ./cmd/``` on mac/linux, or ```go build -o out .\cmd\``` on windows.

Then run ```./out/``` for mac/linux or ```out.exe``` on windows.

## TODO

- write & store multiple websocket connections
- restructure repo
- page feature (50 a page)
- maybe use something other than json for speed?
- CheckOrigin in logger.go
- better auth handling
- better error handling
- cron jobs to clean old logs

## cleaning

- context
- ws on frontend
