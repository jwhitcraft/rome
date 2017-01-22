SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=rome

BUILD_TIME=`date +%FT%T%z`

LDFLAGS=-ldflags "-X github.com/jwhitcraft/rome/cmd.Version=${VERSION} -X github.com/jwhitcraft/rome/cmd.BuildTime=${BUILD_TIME}"

check-env:
ifndef VERSION
	$(error VERSION is undefined)
endif

.DEFAULT_GOAL: all

build = GOOS=$(1) GOARCH=$(2) go build ${LDFLAGS} -o packages/$(1)-$(2)$(3)

.PHONY: all windows darwin linux clean

all: check-env windows darwin linux

.PHONY: dev

dev: $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} main.go

.PHONY: clean
clean:
	if [ -f ./${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -d ./packages ] ; then rm ./packages/* ; fi

##### LINUX BUILDS #####
linux: packages/linux_arm.tar.gz packages/linux_arm64.tar.gz packages/linux_amd64.tar.gz

packages/linux_amd64.tar.gz: $(sources)
	$(call build,linux,amd64,)

packages/linux_arm.tar.gz: $(sources)
	$(call build,linux,arm,)

packages/linux_arm64.tar.gz: $(sources)
	$(call build,linux,arm64,)

##### DARWIN (MAC) BUILDS #####
darwin: packages/darwin_amd64.tar.gz

packages/darwin_amd64.tar.gz: $(sources)
	$(call build,darwin,amd64,)

##### WINDOWS BUILDS #####
windows: packages/windows_amd64.zip

packages/windows_amd64.zip: $(sources)
	$(call build,windows,amd64,.exe)