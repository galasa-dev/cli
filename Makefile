#
# Copyright contributors to the Galasa project
#
# SPDX-License-Identifier: EPL-2.0
#
all: tests galasactl gendocs-galasactl

galasactl: \
	bin/galasactl-linux-arm64 \
	bin/galasactl-linux-x86_64 \
	bin/galasactl-windows-x86_64.exe \
	bin/galasactl-darwin-x86_64 \
	bin/galasactl-darwin-arm64 \
	bin/galasactl-linux-s390x

# 'gendocs-galasactl' is a command-line tool which generates documentation about the galasactl tool.
# When executed, the .md produced contain up-to-date information on tool syntax.
gendocs-galasactl: bin/gendocs-galasactl-darwin-arm64 bin/gendocs-galasactl-linux-arm64 bin/gendocs-galasactl-darwin-x86_64 bin/gendocs-galasactl-linux-x86_64

tests: galasactl-source build/coverage.txt build/coverage.html

build/coverage.out : galasactl-source
	mkdir -p build
	go test -v -cover -coverprofile=build/coverage.out -coverpkg pkg/api,./pkg/auth,./pkg/cmd,./pkg/runsformatter,./pkg/errors,./pkg/launcher,./pkg/utils,./pkg/runs,./pkg/properties,./pkg/propertiesformatter,./pkg/resources ./pkg/...

build/coverage.html : build/coverage.out
	go tool cover -html=build/coverage.out -o build/coverage.html

build/coverage.txt : build/coverage.out
	go tool cover -func=build/coverage.out > build/coverage.txt
	cat build/coverage.txt

galasactl-source : \
	./cmd/galasactl/*.go \
	./pkg/api/*.go \
	./pkg/auth/*.go \
	./pkg/runsformatter/*.go \
	./pkg/cmd/*.go \
	./pkg/utils/*.go \
	./pkg/runs/*.go \
	./pkg/launcher/*.go \
	./pkg/files/*.go \
	./pkg/images/*.go \
	./pkg/props/*.go \
	./pkg/properties/*.go \
	./pkg/propertiesformatter/*.go \
	./pkg/resources/*.go \
	./pkg/spi/*.go \
	./pkg/tokensformatter/*.go \
	embedded_info

# The build process
embedded_info : \
	pkg/embedded/templates/version/build.properties


pkg/embedded/templates/version :
	mkdir -p $@

# Build a properties file containing versions of things.
# Then the galasactl can embed the data and read it at run-time.
pkg/embedded/templates/version/build.properties : VERSION pkg/embedded/templates/version Makefile build.gradle
	echo "# Property file generated at build-time" > $@
	# Turn the contents of VERSION file into a properties file value.
	cat VERSION | sed "s/^/galasactl.version = /1" >> $@ ; echo "" >> $@
	# Add the `galasa.boot.jar.version` property based on the build.gradle value.
	cat build.gradle | grep "def galasaBootJarVersion" | cut -f2 -d\' | sed "s/^/galasa.boot.jar.version = /" >> $@
	# Add the `galasa.framework.version` property based on the build.gradle value.
	cat build.gradle | grep "def galasaFrameworkVersion" | cut -f2 -d\' | sed "s/^/galasa.framework.version = /" >> $@
	# Add the `galasactl.rest.api.version` property based on the build/dependencies/openapi.yaml value.
	echo "" >> $@
	echo "# version of the rest api that is compiled and the client is expecting from the ecosystem." >> $@
	cat build/dependencies/openapi.yaml | grep "version :" | cut -f2 -d'"' | sed "s/^/galasactl.rest.api.version = /" >> $@

bin/galasactl-linux-x86_64 : galasactl-source
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/galasactl-linux-x86_64 ./cmd/galasactl

bin/galasactl-windows-x86_64.exe : galasactl-source
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/galasactl-windows-x86_64.exe ./cmd/galasactl

bin/galasactl-darwin-x86_64 : galasactl-source
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/galasactl-darwin-x86_64 ./cmd/galasactl

bin/galasactl-darwin-arm64 : galasactl-source
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/galasactl-darwin-arm64 ./cmd/galasactl

bin/galasactl-linux-arm64 : galasactl-source
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/galasactl-linux-arm64 ./cmd/galasactl

bin/galasactl-linux-s390x : galasactl-source
	CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build -o bin/galasactl-linux-s390x ./cmd/galasactl

bin/gendocs-galasactl-darwin-arm64 : galasactl-source
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/gendocs-galasactl-darwin-arm64 ./cmd/gendocs-galasactl

bin/gendocs-galasactl-linux-arm64 : galasactl-source
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/gendocs-galasactl-linux-arm64 ./cmd/gendocs-galasactl

bin/gendocs-galasactl-linux-x86_64 : galasactl-source
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/gendocs-galasactl-linux-x86_64 ./cmd/gendocs-galasactl

bin/gendocs-galasactl-darwin-x86_64 : galasactl-source
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/gendocs-galasactl-darwin-x86_64 ./cmd/gendocs-galasactl

clean:
	rm -fr bin/galasactl*
	rm -fr build/*
	rm -fr build/coverage*
	rm -fr pkg/embedded/templates/version

