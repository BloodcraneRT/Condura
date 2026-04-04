# EdgeSynth

EdgeSynth is a polished, multi-platform product for synthetic traffic generation and monitoring. Think of it as a modernized, distributed iperf / Cisco TRex / Kentik Synthetics system for edge networks.

It is distributed as a **single standalone binary** that includes the web dashboard and internal task execution.

## Features
- **Standalone Binary**: Run it anywhere (Linux, Windows, macOS) with zero external dependencies. Pure Go SQLite is embedded.
- **Synthetic Testing**: HTTP latency checks, TCP connection tests, UDP payload tests, and heavy Data Uploads / Downloads (up to GBs).
- **Target Sources**: Target predefined open APIs (like JSONPlaceholder, PokeAPI), loopback data endpoints for load testing, or add your own custom targets.
- **Modern Dashboard**: Beautiful React/Tailwind UI to schedule tests and view results in real-time.

## Architecture
1. The **EdgeSynth** application starts up, initializes its local database, and serves the API + Web UI.
2. An internal background runner routinely polls the database for scheduled synthetic tasks.
3. It natively executes the network tests and stores latency, success state, and throughput (Bytes Downloaded/Uploaded) back into the database.

## Building
A `Makefile` is provided to build the system across platforms.

### Build All
```bash
make all
```

### Build Cross-Platform Release Binaries
```bash
make build-all-platforms
```

Binaries will be output in the `bin/` directory.

## Running

1. **Start EdgeSynth**
```bash
./bin/edgesynth-linux-amd64 --port 8080
```
Open your browser to `http://localhost:8080` to view the dashboard and schedule tasks.
