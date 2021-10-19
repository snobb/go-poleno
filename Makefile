export GO111MODULE=on
export GOVCS=*:git

TARGET   = poleno
MAIN     = ./cmd/main.go
BIN      = ./bin
COVEROUT = cover.out

all: build

lint:
	golangci-lint run

cover:
	go tool cover -html=$(COVEROUT)
	-rm -f $(COVEROUT)

test:
	go test -timeout $(TIMEOUT)s -cover -coverprofile=$(COVEROUT) ./pkg/...

build:
	go build -o $(BIN)/$(TARGET) $(MAIN)

build-linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BIN)/$(TARGET) $(MAIN)

clean:
	-rm -rf $(BIN)
	-rm -f $(COVEROUT)

.PHONY: build
