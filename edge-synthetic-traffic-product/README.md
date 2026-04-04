<div align="center">
  <h1>🚀 EdgeSynth</h1>
  <p><b>A modern, standalone synthetic traffic generator and network observability platform.</b></p>
  <img src="architecture.svg" alt="Architecture Diagram" width="800"/>
</div>

---

## What is EdgeSynth?

EdgeSynth is a polished, multi-platform product for synthetic traffic generation and monitoring. Think of it as a modernized, distributed **iperf**, **Cisco TRex**, or **Kentik Synthetics** system built specifically to run natively as a single executable on edge networks.

Instead of managing complex server-agent relationships, **EdgeSynth is self-contained**. You drop the binary on any machine, and it instantly serves a beautiful React/Tailwind dashboard while seamlessly spinning up Go routines to test, measure, and record synthetic traffic metrics against local or global targets.

## Why EdgeSynth?

* **No Deployment Headaches:** It's distributed as a **single standalone binary**. Pure Go. Embedded React UI. No Node.js, Python, or complex runtimes required on the host.
* **Modern Observability:** Built with a clean, high-contrast UI (inspired by Airbnb and ThousandEyes design systems). It calculates Time To First Byte (TTFB), DNS resolution speeds, network path hops, and connection timings.
* **Massive Data Scale:** Whether you need to simulate a 1KB JSON payload or heavily saturate a link with a **10GB data download**, EdgeSynth handles it safely via memory-optimized stream discarding. Results are dumped into **ClickHouse** (`MergeTree` engine) for blistering fast analytical reads.

<br/>
<div align="center">
  <img src="lifecycle.svg" alt="Task Execution Lifecycle" width="800"/>
</div>
<br/>

## Key Features

| Feature | Description |
| :--- | :--- |
| 🗄️ **ClickHouse Backend** | Uses ClickHouse `MergeTree` engines for highly efficient timeseries inserts and aggregations. |
| 🌐 **HTTP & DNS Tracing** | Detailed lifecycle metrics (TTFB, DNS time, Connection overhead) mimicking enterprise tools. |
| 📡 **Network Path Tracing** | Native OS-level Traceroute executions to capture routing hops and path statuses directly from the edge. |
| 📦 **Gigabyte Payloads** | Dedicated HTTP endpoints to generate and absorb arbitrary sizes of synthetic data up to 5GB per test. |
| 🕵️‍♂️ **PCAP Replay** | Replay network traffic directly via `tcpreplay` over specific interfaces (e.g. `eth0`) based on `CESNET/FlowTest`. |
| ✨ **Embedded UI** | Zero-config modern dashboard served directly from the Go binary via `go:embed`. |

## Architecture Overview
1. The **EdgeSynth** application starts up, connects to your ClickHouse instance, and serves the API + Web UI.
2. The UI lets you configure new "Sources" (e.g. your production API, Github, or a 5GB loopback download).
3. You schedule tasks (Ping, TCP, Web Load, Traceroute, PCAP Replay).
4. An internal Go routine routinely polls the database for scheduled synthetic tasks and executes them asynchronously.
5. It natively executes the network tests and stores latency, success state, and throughput (Bytes Downloaded/Uploaded) back into ClickHouse.

---

## Building from Source

A `Makefile` is provided to seamlessly build the React UI and compile the Go binaries across platforms.

### Build All
```bash
make all
```

### Build Cross-Platform Release Binaries
```bash
make build-all-platforms
```

Binaries will be output in the `bin/` directory.

---

## Running EdgeSynth

### 1. Start ClickHouse
EdgeSynth requires ClickHouse to store its metrics. You can run it locally via Docker:
```bash
docker run -d --name clickhouse-server --ulimit nofile=262144:262144 -p 8123:8123 -p 9000:9000 clickhouse/clickhouse-server
```

### 2. Start the Application
By default, the binary will look for ClickHouse on `localhost:9000`.
```bash
./bin/edgesynth-linux-amd64 --port 8080
```

Open your browser to `http://localhost:8080` to view the dashboard and schedule tasks!

*(Optional)* You can override the database URL using the `DATABASE_URL` environment variable or the `--db` flag:
```bash
DATABASE_URL="clickhouse://default:@host:9000/default" ./bin/edgesynth-linux-amd64
```
