SHELL = /bin/bash
GOPATH = $(shell pwd)
GOBIN = $(GOPATH)/bin
TARGET1 = shlancd
TARGET2 = shlanc
BIN	= $(DESTDIR)/usr/bin
CONF = $(DESTDIR)/etc/shlanc


## These will be provided to the target
VERSION=0.25-0ubuntu1
BUILD = $(shell date +%s)

## Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

export PATH := $(PATH):/usr/local/go/bin

.PHONY: all install uninstall clean ppa source source_us

all:
	@echo "##################################"
	@echo "# 	Compile binaries"
	@echo "##################################"
	@echo $(LDFLAGS)

	mkdir -p $(GOBIN)
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(LDFLAGS) $(TARGET1) $(TARGET2)

	@echo " *** Done ***"
	@echo ""


install: all

	mkdir -p $(BIN)

	@echo " *** Install server"
	install $(GOBIN)/$(TARGET1) $(BIN)

	@echo " *** Install client"
	install $(GOBIN)/$(TARGET2) $(BIN)

	@echo " *** Install config"
	mkdir -p $(CONF)
	install ./config.json $(CONF)/config.json


uninstall:
	rm -rf $(BIN)/$(TARGET1)
	rm -rf $(BIN)/$(TARGET2)
	rm -rf $(CONF)/config.json


configure_test:
	@echo "##############################"
	@echo "# Compile TEST build"
	@echo "##############################"

	$(MAKE) clean
	./configure_ppa.sh $(VERSION) 1

	cd build/tmp; debuild  -S -us -uc

	@echo " *** Build(TEST) is compiled ***"
	@echo ""


configure:
	@echo "###############################"
	@echo "# Compile build"
	@echo "###############################"

	$(MAKE) clean
	./configure_ppa.sh $(VERSION) 0

	### build package (https://help.launchpad.net/Packaging/PPA/BuildingASourcePackage)
	cd build/tmp; debuild -S -sa

	@echo " *** Build is compiled ***"
	@echo ""


ppa: build_test configure
	@echo " ==> Uploading to PPA..."
	dput -d ppa:onm/shlanc $(shell ls build/*.changes)
	$(MAKE) clean



build_test: configure_test
	@echo "###############################"
	@echo "#**** 	TESTING BUILD	 ****#"
	@echo "###############################"

	@echo " ==> Unpacking..."
	cd build/; dpkg-source -x *.dsc

	@echo " ==> Testing..."
	cd build/shlanc-$(VERSION); dh_auto_test;

	@echo " ==> Compiling..."
	cd build/shlanc-$(VERSION); dh_auto_build -a

	@echo " *** Build(TEST) is OK ***"
	@echo ""


clean:
	@rm -rf build/
