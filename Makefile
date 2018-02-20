#SHELL = /bin/bash
#PREFIX = /usr/local/bin
#GOPATH = $(shell pwd)
#GOBIN = $(GOPATH)/bin
#TARGET1 = shlancd
#TARGET2 = shlanc
#CONF = /etc/shlanc
#
#
## These will be provided to the target
#VERSION := 1.0.0
#BUILD := `git rev-parse HEAD`
#
## Use linker flags to provide version/build settings to the target
#LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
#
#.PHONY: all install uninstall
#
#all:
#	mkdir -p $(GOBIN)
#	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get .
#	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(LDFLAGS) $(TARGET1) $(TARGET2)
#
#install:
#	install $(GOBIN)/$(TARGET1) $(PREFIX)
#	install $(GOBIN)/$(TARGET2) $(PREFIX)
#
#	@mkdir $(CONF) -p
#	install ./config.json $(CONF)/config.json
#
#uninstall:
#	rm -rf $(PREFIX)/$(TARGET1)
#	rm -rf $(PREFIX)/$(TARGET2)
#	rm -rf $(CONF)/config.json


TARGET = shlancd
BIN	= $(DESTDIR)/usr/bin
CONF = /etc/shlancd

# These will be provided to the target
VERSION := 1.0.0
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

.PHONY: install uninstall

install:
	install ./bin/$(TARGET) $(BIN)/$(TARGET)

	@mkdir $(CONF) -p
	install ./config.json $(CONF)/config.json

uninstall:
	rm -rf $(BIN)/$(TARGET)
	rm -rf $(CONF)/config.json
