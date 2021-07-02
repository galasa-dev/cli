all: bin/galasactl



bin/galasactl: ./cmd/galasactl/main.go ./pkg/cmd/root.go ./pkg/cmd/runs.go ./pkg/cmd/runsAssemble.go
	go build -o bin/galasactl ./cmd/galasactl


clean:
	rm -rf bin
