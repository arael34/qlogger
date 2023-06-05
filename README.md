# QLogger

A simple logging client.

v 1.0.1

## To set up locally

Install go and git, then install deps with ```go mod download```

Create .env file according to the example

Run ```go run main.go``` for development.

For production, first run ```go build -o out``` to build the binary.

Then run ```./out/``` for mac/linux or ```out.exe``` on windows.

## TODO

- page feature (50 a page)
- maybe use something other than json for speed?
- CheckOrigin in logger.go
- better auth handling
- better error handling
- cron jobs to clean old logs

## cleaning

- context
- ws on frontend
