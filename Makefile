SHELL = /bin/bash
BIN	= $(DESTDIR)/usr/bin
GOPATH = $(shell pwd)
GOBIN = $(GOPATH)/bin
TARGET1 = shlancd
TARGET2 = shlanc
CONF = /etc/shlanc


## These will be provided to the target
VERSION := 1.0.0
BUILD = $(shell date +%s)
#
## Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

.PHONY: all install uninstall

all:
	mkdir -p $(GOBIN)
	@echo $(LDFLAGS)
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(LDFLAGS) $(TARGET1) $(TARGET2)

install:
	install $(GOBIN)/$(TARGET1) $(BIN)
	install $(GOBIN)/$(TARGET2) $(BIN)

	@mkdir $(CONF) -p
	install ./config.json $(CONF)/config.json

uninstall:
	rm -rf $(BIN)/$(TARGET1)
	rm -rf $(BIN)/$(TARGET2)
	rm -rf $(CONF)/config.json