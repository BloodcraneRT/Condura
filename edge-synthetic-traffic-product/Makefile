.PHONY: all build-ui build-server clean

all: build-ui build-server

build-ui:
	@echo "Building UI..."
	cd ui && npm install && npm run build
	@echo "Copying UI build to server..."
	rm -rf server/ui-dist
	cp -r ui/dist server/ui-dist

build-server: build-ui
	@echo "Building EdgeSynth..."
	cd server && go build -o ../bin/edgesynth

build-all-platforms: build-ui
	@echo "Building cross-platform binaries..."
	mkdir -p bin
	# Linux
	GOOS=linux GOARCH=amd64 cd server && go build -o ../bin/edgesynth-linux-amd64
	# Windows
	GOOS=windows GOARCH=amd64 cd server && go build -o ../bin/edgesynth-windows-amd64.exe
	# macOS (Darwin)
	GOOS=darwin GOARCH=amd64 cd server && go build -o ../bin/edgesynth-darwin-amd64
	GOOS=darwin GOARCH=arm64 cd server && go build -o ../bin/edgesynth-darwin-arm64

clean:
	rm -rf bin/
	rm -rf server/ui-dist/
	rm -rf ui/dist/
	rm -rf server/edgesynth.db
