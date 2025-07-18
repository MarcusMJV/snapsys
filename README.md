# SnapSys

![Go](https://img.shields.io/badge/Go-1.23+-blue?logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/Platform-Linux-success?logo=linux&logoColor=white)

SnapSys is a lightweight, purpose-built CLI tool for capturing CPU, memory, and disk usage at fixed intervals, outputting the results in clean, structured JSONL format. Designed for system administrators, developers, DevOps engineers, and performance testers, SnapSys is ideal for monitoring system behavior, debugging performance issues, or capturing lightweight metrics during CI/CD pipelines or local build processes. Whether you're analyzing resource usage in virtual machines, testing workloads in Docker containers, or gathering evidence for incident reports, SnapSys fits seamlessly into your workflow.

> ⚠️ SnapSys is currently developed and tested for Linux systems only.

## Key Features

* Track **CPU, Memory, and Disk** usage at defined intervals
* Output in **JSONL** format for easy processing and analysis
* Minimal dependencies – just Go and your shell
* Includes automatic fallback output path if none is provided

## Installation (via `go install`)

To install the latest version:

```bash
go install github.com/MarcusMJV/snapsys.git@latest
```

Ensure `~/go/bin` is in your `$PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Once installed, use it like:

```bash
snapsys snapshot
```
## Download a Precompiled Binary

Visit the Releases Page

Download the binary called snapsys

In your terminal run
```bash
chmod +x snapsys
```
Run it:
```bash
./snapsys snapshot
```
To simplify usage:
```bash
sudo mv snapsys /usr/local/bin/snapsys
snapsys snapshot
```

---

## Usage

### Basic Command

```bash
snapsys snapshot
```

This will collect system stats every 3 seconds for 30 seconds, and write to a default path like:

```bash
./snapsys_runs/snaprun_2025-06-16_13-00-00.jsonl
```

### Command Flags

| Flag         | Default            | Description                                     |
| ------------ | ------------------ | ----------------------------------------------- |
| `--duration` | `30s`              | Total time to run the snapshot                  |
| `--interval` | `3s`               | Time between each snapshot (minimum 1s)         |
| `--output`   | *(auto-generated)* | Output file path (.jsonl required if specified) |

> ⚠️ If the interval is set below 1s, SnapSys will default it to `1s` and warn the user.

> ⚠️ If you use `--output`, the file **must** end in `.jsonl`.

---

## Use Cases

SnapSys is ideal for:

* Lightweight benchmarking in Linux VMs or containers
* Capturing performance stats during builds or deployments
* Logging system behavior for incident reports
* Analyzing trends across time during load tests

---

## Example Output (JSONL)

Each snapshot is written as a line in JSON:

```json
{"timestamp":"2025-06-16T13:00:00Z","cpu":{"usage_pct":12.8, ...},"memory":{"usage_pct":31.9, ...},"disks":{"/":{...}}}
```

This is great for:

* Ingesting into tools like Elasticsearch
* Plotting in Python, Grafana, etc.
* Please see example.jsonl file for full example 

---

## Future Plans

* Sub-second snapshot support using concurrent snapshot routines
* Additional system metrics:
    - Uptime
    - Hostname and kernel version
    - Network interface stats (bytes in/out, errors, drops)
    - System load average

* Snapshot tagging and summaries
* Ability to push .jsonl data to remote HTTP APIs
* Support for compressed output (.jsonl.gz)
* Interactive shell mode for quick inspection and dev workflows
* CSV export format option

---

PRs welcome!
