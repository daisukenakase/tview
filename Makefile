BINARY := tview
PREFIX ?= /usr/local
INSTALL_DIR := $(PREFIX)/bin

.PHONY: build install uninstall clean build-arm build-arm64 dist

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BINARY) .

# Raspberry Pi (32bit OS, Pi Zero/1/2/3/4)
build-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -trimpath -ldflags="-s -w" -o $(BINARY)-armv6 .

# Raspberry Pi (64bit OS, Pi 3/4/5)
build-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o $(BINARY)-arm64 .

dist: build-arm build-arm64

install: build
	install -Dm755 $(BINARY) $(INSTALL_DIR)/$(BINARY)

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY)

clean:
	rm -f $(BINARY) $(BINARY)-armv6 $(BINARY)-arm64
