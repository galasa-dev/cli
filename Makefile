#
# Licensed Materials - Property of IBM
#
# (c) Copyright IBM Corp. 2021.
#

all: tests galasactl gendocs-galasactl

galasactl: bin/galasactl-linux-amd64 bin/galasactl-windows-amd64.exe bin/galasactl-darwin-amd64 bin/galasactl-darwin-arm64 bin/galasactl-linux-s390x

# 'gendocs-galasactl' is a command-line tool which generates documentation about the galasactl tool.
# When executed, the .md produced contain up-to-date information on tool syntax.
gendocs-galasactl: bin/gendocs-galasactl-darwin-arm64 bin/gendocs-galasactl-darwin-amd64 bin/gendocs-galasactl-linux-amd64

tests: galasactl-source
	mkdir -p build
	go test -v -cover -coverprofile=build/coverage.out -coverpkg ./pkg/cmd,./pkg/errors,.pkg/launcher,./pkg/utils,./pkg/runs ./pkg/...
	go tool cover -html=build/coverage.out -o build/coverage.html
	go tool cover -func=build/coverage.out > build/coverage.txt
	cat build/coverage.txt

galasactl-source : ./cmd/galasactl/*.go ./pkg/api/*.go ./pkg/cmd/*.go ./pkg/utils/*.go ./pkg/runs/*.go ./pkg/launcher/*.go

# when the gradle stuff works, we can rely on this jar being here: ./pkg/embedded/templates/galasahome/lib/*.jar 

bin/galasactl-linux-amd64 : galasactl-source 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/galasactl-linux-amd64 ./cmd/galasactl

bin/galasactl-windows-amd64.exe : galasactl-source
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/galasactl-windows-amd64.exe ./cmd/galasactl

bin/galasactl-darwin-amd64 : galasactl-source
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/galasactl-darwin-amd64 ./cmd/galasactl

bin/galasactl-darwin-arm64 : galasactl-source
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/galasactl-darwin-arm64 ./cmd/galasactl	

bin/galasactl-linux-s390x : galasactl-source
	CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build -o bin/galasactl-linux-s390x ./cmd/galasactl

bin/gendocs-galasactl-darwin-arm64 : galasactl-source
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/gendocs-galasactl-darwin-arm64 ./cmd/gendocs-galasactl

bin/gendocs-galasactl-linux-amd64 : galasactl-source
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/gendocs-galasactl-linux-amd64 ./cmd/gendocs-galasactl

bin/gendocs-galasactl-darwin-amd64 : galasactl-source 
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/gendocs-galasactl-darwin-amd64 ./cmd/gendocs-galasactl

clean:
	rm -rf bin

