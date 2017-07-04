.PHONY: help
.DEFAULT_GOAL := help

BLDDIR = _build
BINDIR = _bin
BLDDATE=$(shell date -u +%Y%m%dT%H%M%S)
VERSION= 1
SRCS = $(wildcard *.go ./**/*.go)

BINNAME="todo"
PROJECT="joelclouddistrict"
GITPROJECT="gochecklist"
ORG_PATH=github.com/joelclouddistrict
REPO_PATH=$(ORG_PATH)/$(GITPROJECT)

LDFLAGS=" -s -X $(REPO_PATH)/version.Name=$(BINNAME) -X $(REPO_PATH)/version.Version=$(VERSION)"

export PATH := $(PWD)/_bin:$(PATH)

dep: ## Vendor go dependencies
	@echo "Vendoring dependencies"
	@go get -u github.com/FiloSottile/gvt
	@gvt rebuild

$(BLDDIR):
	mkdir ${BLDDIR} || true

$(BINDIR):
	mkdir ${BINDIR} || true

linux: $(BLDDIR) ## build linux amd64 binary
	$(shell GO15VENDOREXPERIMENT=1 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -ldflags ${LDFLAGS} -a -installsuffix cgo \
	-o ${BLDDIR}/${BINNAME}_${VERSION}_srv_linux_amd64.bin . \
	&& chmod +x ${BLDDIR}/${BINNAME}_${VERSION}_srv_linux_amd64.bin \
	)

osx: $(BLDDIR) ## build darwin amd64 binary
	$(shell GO15VENDOREXPERIMENT=1 CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
	go build -ldflags ${LDFLAGS} -a -installsuffix cgo \
	-o ${BLDDIR}/${BINNAME}_${VERSION}_srv_darwin_amd64.bin . \
	&& chmod +x ${BLDDIR}/${BINNAME}_${VERSION}_srv_darwin_amd64.bin \
	)

pb: $(BINDIR) _bin/protoc _bin/protoc-gen-go ## compile the protocol buffers files into resources
	_bin/protoc -I/usr/local/include -I. \
		-I$$PWD/vendor \
		-I$$PWD/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
		services/*.proto ;

gw:  $(BINDIR) _bin/protoc _bin/protoc-gen-go _bin/protoc-gen-grpc-gateway ## build the REST gateway resource files
	@_bin/protoc -I/usr/local/include -I. \
		-I$$PWD/vendor \
		-I$$PWD/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		services/*.proto ;

doc: $(BINDIR) _bin/protoc _bin/protoc-gen-go _bin/protoc-gen-grpc-gateway _bin/protoc-gen-swagger ## build the documentation files: swagger json
	@_bin/protoc -I/usr/local/include -I. \
		-I$$PWD/vendor \
		-I$$PWD/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:. \
		services/*.proto ;

allpb: pb gw ## build proto and gateway files

_bin/protoc:
	@./scripts/get-protoc ${BINDIR}/protoc

_bin/protoc-gen-go:
	@go build -o ${BINDIR}/protoc-gen-go $(REPO_PATH)/vendor/github.com/golang/protobuf/protoc-gen-go

_bin/protoc-gen-grpc-gateway:
	@go build -o ${BINDIR}/protoc-gen-grpc-gateway $(REPO_PATH)/vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

_bin/protoc-gen-swagger:
	@go build -o ${BINDIR}/protoc-gen-swagger $(REPO_PATH)/vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
