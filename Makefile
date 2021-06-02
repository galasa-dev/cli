all: bin/galasactl



bin/galasactl:
	go build -o bin/galasactl ./cmd/galasactl


clean:
	rm -rf bin
