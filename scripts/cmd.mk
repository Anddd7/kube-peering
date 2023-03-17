# required parameters
# COMMAND

ROOT_PATH 		= $(shell git rev-parse --show-toplevel 2>/dev/null)
BUILD_BINARY 	= $(ROOT_PATH)/bin/$(COMMAND)

run:
	go run . $(args)

build:
	go build -o $(BUILD_BINARY) ./

test:
	go test ./...

clean:
	rm $(BUILD_BINARY)

install: build
	mv $(BUILD_BINARY) /usr/local/bin/$(COMMAND) 

uninstall:
	rm /usr/local/bin/$(COMMAND)