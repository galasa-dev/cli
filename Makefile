all: bin/galasactl



bin/galasactl: ./cmd/galasactl/main.go
	go build -o bin/galasactl ./cmd/galasactl


clean:
	rm -rf bin
