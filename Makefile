export GO111MODULE=on
export GOVCS=*:git

TARGET   = poleno
NAME     = poleno
MAIN     = ./cmd/main.go
BIN      = ./bin
GO       = go
TIMEOUT  = 15
COVEROUT = cover.out

all: build

lint:
	golangci-lint run

run:
	go run $(MAIN)

dev-run:
	go run $(MAIN) | poleno

cover:
	$(GO) tool cover -html=$(COVEROUT)
	-rm -f $(COVEROUT)

test:
	$(GO) test -timeout $(TIMEOUT)s -cover -coverprofile=$(COVEROUT) ./pkg/...

build:
	$(GO) build -o ./bin/$(TARGET) $(MAIN)

build-linux:
	CGO_ENABLED=0 GOOS=linux $(GO) build -a -installsuffix cgo -o $(BIN)/$(TARGET) $(MAIN)

dockerise: build-linux
	docker build -t $(NAME) .
	docker run --rm -h $(shell hostname) $(NAME)

rpm: clean build-linux
	./scripts/build_rpm.sh

generate mocks:
	$(GO) generate ./pkg/...

clean:
	-rm -rf $(BIN)
	-rm -rf ./dist
	-rm -f $(COVEROUT)

.PHONY: build
