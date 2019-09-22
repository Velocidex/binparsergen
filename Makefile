all:
	go build -o ./binparsegen cmd/*.go

install: all
	mv ./binparsegen ${GOPATH}/bin/
