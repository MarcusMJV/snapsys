# SnapSys

SnapSys is a Linux terminal utility written in Go that captures and logs system resource snapshots at configurable intervals.  
It is designed for lightweight performance monitoring, diagnostics, and learning purposes.

---

## Features

- Interval-based CPU usage snapshotting
- JSON output logging
- Configurable duration and intervals
- Modular Go architecture using Cobra CLI
- Designed specifically for Linux systems (reads from `/proc`)
- Simple and extensible structure

---

## Example Usage

```bash
snapsys snapshot --interval 5s --duration 1m --output cpu.json
