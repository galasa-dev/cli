all: bin/galasactl



bin/galasactl: ./cmd/galasactl/main.go ./pkg/api/api.go ./pkg/cmd/root.go ./pkg/cmd/runs.go ./pkg/cmd/runsAssemble.go ./pkg/utils/javaProperties.go  ./pkg/utils/testStream.go  ./pkg/utils/testCatalog.go  ./pkg/utils/testSelection.go  ./pkg/utils/portfolio.go
	go build -o bin/galasactl ./cmd/galasactl


clean:
	rm -rf bin
