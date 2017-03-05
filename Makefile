SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

.DEFAULT_GOAL := dev
BINARY=rome
BUILD_TIME=`date +%FT%T%z`

LDFLAGS=-ldflags "-X github.com/jwhitcraft/rome/cmd.Version=${VERSION} -X github.com/jwhitcraft/rome/cmd.BuildTime=${BUILD_TIME} -s -w"

check-env:
ifndef VERSION
	$(error VERSION is undefined)
endif

GCFLAGS=-gcflags "-N -l"


build = GOOS=$(1) GOARCH=$(2) go build ${LDFLAGS} ${GCFLAGS} -o packages/$(1)-$(2)$(3)
rename = cp packages/$(1)-$(2)$(3) public/${BINARY}-$(4)-$(5)$(3)

release: check-env clean windows darwin linux

dev: $(SOURCES)
	go build ${LDFLAGS} ${GCFLAGS} -o ${BINARY} main.go

.PHONY: clean
clean:
	if [ -f ./${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -d ./packages ] ; then rm ./packages/* ; fi

test:
	go test -v `glide novendor`

.PHONY: aqueduct
aqueduct:
	protoc -I aqueduct/ aqueduct/aqueduct.proto --go_out=plugins=grpc:aqueduct

##### LINUX BUILDS #####
linux: packages/linux_amd64.tar.gz

packages/linux_amd64.tar.gz: $(sources)
	$(call build,linux,amd64,)
	$(call rename,linux,amd64,,Linux,x86_64)

packages/linux_arm.tar.gz: $(sources)
	$(call build,linux,arm,)

packages/linux_arm64.tar.gz: $(sources)
	$(call build,linux,arm64,)

##### DARWIN (MAC) BUILDS #####
darwin: packages/darwin_amd64.tar.gz

packages/darwin_amd64.tar.gz: $(sources)
	$(call build,darwin,amd64,)
	$(call rename,darwin,amd64,,Darwin,x86_64)

##### WINDOWS BUILDS #####
windows: packages/windows_amd64.zip

packages/windows_amd64.zip: $(sources)
	$(call build,windows,amd64,.exe)
	$(call rename,windows,amd64,.exe,Windows,x86_64)
