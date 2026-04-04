.PHONY: all build-ui build-server build-agent clean

all: build-ui build-server build-agent

build-ui:
	@echo "Building UI..."
	cd ui && npm install && npm run build
	@echo "Copying UI build to server..."
	rm -rf server/ui-dist
	cp -r ui/dist server/ui-dist

build-server: build-ui
	@echo "Building Server..."
	cd server && go build -o ../bin/edgesynth-server

build-agent:
	@echo "Building Agent..."
	cd agent && go build -o ../bin/edgesynth-agent

build-all-platforms: build-ui
	@echo "Building cross-platform binaries..."
	mkdir -p bin
	# Linux
	GOOS=linux GOARCH=amd64 cd server && go build -o ../bin/edgesynth-server-linux-amd64
	GOOS=linux GOARCH=amd64 cd agent && go build -o ../bin/edgesynth-agent-linux-amd64
	# Windows
	GOOS=windows GOARCH=amd64 cd server && go build -o ../bin/edgesynth-server-windows-amd64.exe
	GOOS=windows GOARCH=amd64 cd agent && go build -o ../bin/edgesynth-agent-windows-amd64.exe
	# macOS (Darwin)
	GOOS=darwin GOARCH=amd64 cd server && go build -o ../bin/edgesynth-server-darwin-amd64
	GOOS=darwin GOARCH=amd64 cd agent && go build -o ../bin/edgesynth-agent-darwin-amd64
	GOOS=darwin GOARCH=arm64 cd server && go build -o ../bin/edgesynth-server-darwin-arm64
	GOOS=darwin GOARCH=arm64 cd agent && go build -o ../bin/edgesynth-agent-darwin-arm64

clean:
	rm -rf bin/
	rm -rf server/ui-dist/
	rm -rf ui/dist/
	rm -rf server/edgesynth.db
