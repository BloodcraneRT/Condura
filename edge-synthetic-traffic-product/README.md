# EdgeSynth

EdgeSynth is a polished, multi-platform product for synthetic traffic generation and monitoring. Think of it as a modernized, distributed iperf / Cisco TRex / Kentik Synthetics system for edge networks.

It consists of two main components:
- **Control Plane Server**: A Go-based backend (with an embedded SQLite database) that manages agents, schedules tests, and serves a React/Tailwind frontend.
- **Edge Agent**: A lightweight Go binary deployed on edge nodes (Linux, Windows, macOS). It connects to the server, pulls down synthetic task configurations (Ping, HTTP, TCP, UDP), runs them, and pushes results back.

## Features
- **Multi-platform Agents**: Run the agent on Linux, Windows, or macOS.
- **Synthetic Testing**: HTTP latency checks, TCP connection tests, and UDP payload tests.
- **Centralized Dashboard**: Beautiful React/Tailwind UI to monitor agent status and view test results.
- **Easy Deployment**: Single binary for the server (with embedded UI) and single binary for the agents.

## Architecture
1. The **Server** starts up, initializes the database, and serves the API + UI.
2. An **Agent** starts up, registers with the Server via `/api/v1/agents/register`, and begins sending heartbeats.
3. The Agent periodically polls `/api/v1/tasks` for work.
4. When tasks are present, the Agent runs the synthetic network test (e.g., HTTP request) and POSTs the latency/success metrics to `/api/v1/results`.

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

1. **Start the Server**
```bash
cd server
./../bin/edgesynth-server --port 8080
```
Open your browser to `http://localhost:8080` to view the dashboard.

2. **Start an Agent**
In another terminal:
```bash
cd agent
./../bin/edgesynth-agent --server http://localhost:8080
```
