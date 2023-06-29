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

- delete old logs (with a time)

- page feature (50 a page)
- maybe use something other than json for speed?
- better auth handling
- better error handling

- cron job to clean reallyyyy old logs

## cleaning

- context
