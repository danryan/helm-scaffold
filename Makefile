HELM_HOME ?= $(shell helm home)
HELM_PLUGIN_DIR ?= $(HELM_HOME)/plugins/helm-scaffold
HAS_DEP := $(shell command -v dep;)
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
DIST := $(CURDIR)/_dist
LDFLAGS := "-X main.version=${VERSION}"
OUT := "scaffold" 

.PHONY: install
install: bootstrap build
	cp scaffold $(HELM_PLUGIN_DIR)
	cp plugin.yaml $(HELM_PLUGIN_DIR)

.PHONY: hook-install
hookInstall: bootstrap build

.PHONY: build
build:
	go build -o scaffold -ldflags $(LDFLAGS) ./main.go

.PHONY: dist
dist:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -o scaffold -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/helm-scaffold-linux-$(VERSION).tgz scaffold README.md LICENSE.txt plugin.yaml templates
	GOOS=darwin GOARCH=amd64 go build -o scaffold -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/helm-scaffold-macos-$(VERSION).tgz scaffold README.md LICENSE.txt plugin.yaml templates
	GOOS=windows GOARCH=amd64 go build -o scaffold.exe -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/helm-scaffold-windows-$(VERSION).tgz scaffold.exe README.md LICENSE.txt plugin.yaml templates

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
	dep ensure