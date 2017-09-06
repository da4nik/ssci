BINARY=ssci
VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`
REPO=da4nik/$(BINARY)
LDFLAGS=-ldflags "-linkmode external -s -w -extldflags -static -X github.com/$(REPO)/main.version=${VERSION} -X github.com/$(REPO)/main.buildTime=${BUILD_TIME}"

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

.PHONY: build run install clean image test
.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	# glide install --skip-test
	go build ${LDFLAGS} -o ${BINARY} ${BINARY}.go

build: $(BINARY)

run:
	@go run ${BINARY}.go

install:
	@go install ${LDFLAGS} ./...

image:
	docker build -t ${BINARY} .

test:
	@go test

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
