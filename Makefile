# build a linux binary and host OS binary (mac in my case)
build:
	GOOS=linux GOARCH=amd64 go build -o demo-linux-amd64 .
	go build -o demo .
