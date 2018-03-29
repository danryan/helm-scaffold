HELM_HOME ?= $(shell helm home)
HELM_PLUGIN_DIR ?= $(HELM_HOME)/plugins/helm-scaffold
HAS_DEP := $(shell command -v dep;)
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
DIST := $(CURDIR)/_dist
LDFLAGS := "-X main.version=${VERSION}"
OUT := "scaffold" 

.PHONY: install
install: build
	cp $(DIST)/darwin/scaffold $(HELM_PLUGIN_DIR)
	cp plugin.yaml $(HELM_PLUGIN_DIR)
	cp -Rf templates $(HELM_PLUGIN_DIR)

.PHONY: hook-install
hook-install: bootstrap build

.PHONY: build
build:
	go build -o $(DIST)/darwin/scaffold -ldflags $(LDFLAGS) ./main.go


.PHONY: dist
dist:
	mkdir -p $(DIST)/linux
	GOOS=linux GOARCH=amd64 go build -o $(DIST)/linux/scaffold -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/helm-scaffold-linux-$(VERSION).tgz README.md LICENSE.txt plugin.yaml templates -C $(DIST)/linux/ scaffold
	mkdir -p $(DIST)/darwin
	GOOS=darwin GOARCH=amd64 go build -o $(DIST)/darwin/scaffold -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/helm-scaffold-darwin-$(VERSION).tgz README.md LICENSE.txt plugin.yaml templates -C $(DIST)/darwin/ scaffold 
	mkdir -p $(DIST)/windows
	GOOS=windows GOARCH=amd64 go build -o $(DIST)/windows/scaffold.exe -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/helm-scaffold-windows-$(VERSION).tgz README.md LICENSE.txt plugin.yaml templates -C $(DIST)/windows/ scaffold.exe

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
	dep ensure