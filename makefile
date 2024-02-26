BINARY_NAME := wordcombiner
OSX_BINARY := ./build/${BINARY_NAME}_osx_arm64
WIN_BINARY := ./build/${BINARY_NAME}_win_amd64.exe
ALL_SRC_GO := ./cmd/main.go

.PHONY: build run clean dev
build:
	mkdir -p ./build
	GOOS=darwin GOARCH=arm64 go build -ldflags "-w" -o ${OSX_BINARY} ${ALL_SRC_GO}
	GOOS=windows GOARCH=amd64 go build -ldflags "-w" -o ${WIN_BINARY} ${ALL_SRC_GO}
	du -sh ./build/*

run: build
	${OSX_BINARY}

dev:
	go run ${ALL_SRC_GO}

clean:
	go clean
	rm -f ${OSX_BINARY}
	rm -f ${WIN_BINARY}

fmt:
	go fmt ./...