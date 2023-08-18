CMD_PKG=./cmd/httpserver
BINARY_NAME=httpserver

all: build
.PHONY: all

build: clean
	go build -o $(BINARY_NAME) $(CMD_PKG)
.PHONY: build

clean:
	rm -f ./$(BINARY_NAME)
.PHONY: clean
