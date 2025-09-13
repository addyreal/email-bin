.PHONY: all install

BINARY=email-bin
export PATH := /usr/local/go/bin:$(PATH)

all:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY)

install: all
	sudo cp $(BINARY) /usr/local/bin
