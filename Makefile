MODULE := github.com/RamXX/nd
BIN    := nd
PREFIX := $(HOME)/.local/bin

.PHONY: build test install clean vet

build:
	go build -o $(BIN) .

test:
	go test ./...

vet:
	go vet ./...

install: build
	mkdir -p $(PREFIX)
	cp $(BIN) $(PREFIX)/$(BIN)

clean:
	rm -f $(BIN)
