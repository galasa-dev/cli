#
# Licensed Materials - Property of IBM
#
# (c) Copyright IBM Corp. 2021.
#

all: bin/galasactl-linux-amd64 bin/galasactl-windows-amd64.exe bin/galasactl-darwin-amd64 bin/galasactl-darwin-arm64 bin/galasactl-linux-s390x

bin/galasactl-linux-amd64 : ./cmd/galasactl/*.go ./pkg/api/*.go ./pkg/cmd/*.go ./pkg/utils/*.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/galasactl-linux-amd64 ./cmd/galasactl

bin/galasactl-windows-amd64.exe : ./cmd/galasactl/*.go ./pkg/api/*.go ./pkg/cmd/*.go ./pkg/utils/*.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/galasactl-windows-amd64.exe ./cmd/galasactl

bin/galasactl-darwin-amd64 : ./cmd/galasactl/*.go ./pkg/api/*.go ./pkg/cmd/*.go ./pkg/utils/*.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/galasactl-darwin-amd64 ./cmd/galasactl

bin/galasactl-darwin-arm64 : ./cmd/galasactl/*.go ./pkg/api/*.go ./pkg/cmd/*.go ./pkg/utils/*.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/galasactl-darwin-arm64 ./cmd/galasactl

bin/galasactl-linux-s390x : ./cmd/galasactl/*.go ./pkg/api/*.go ./pkg/cmd/*.go ./pkg/utils/*.go
	CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build -o bin/galasactl-linux-s390x ./cmd/galasactl

clean:
	rm -rf bin
