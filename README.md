# gosysutil

Lightweight, dependency-free Go package for monitoring Linux system statistics.
It is designed to be "plug and play", offering a simple API to retrieve CPU, Memory, Disk, and Network usage.

## Features

- **CPU**: Aggregated CPU usage stats (User, System, Idle, Iowait, etc.) from `/proc/stat`.
- **Memory**: Total, Used, Free, Buffers, Cached, Swap stats from `/proc/meminfo`.
- **Disk**: I/O statistics (Reads, Writes, IO Time) for physical disks from `/proc/diskstats`.
- **Network**: Traffic statistics (RX/TX bytes, packets, drops) for network interfaces from `/proc/net/dev`.
- **Load Average**: System load averages (1m, 5m, 15m) from `/proc/loadavg`.
- **Uptime**: System uptime with human-readable formatting from `/proc/uptime`.
- **Top Processes**: Top CPU and RAM consuming processes.
- **VPS Metrics**: CPU steal time and IO wait percentages for virtualized environments.
- **Socket Stats**: TCP/UDP connection counts and memory usage from `/proc/net/sockstat`.
- **File Descriptors**: System-wide file descriptor usage from `/proc/sys/fs/file-nr`.
- **Pressure (PSI)**: CPU, Memory, and IO pressure stall information from `/proc/pressure/*`.
- **VM Stats**: Page faults, paging, swap activity, OOM kills from `/proc/vmstat`.
- **SNMP Stats**: IP/TCP/UDP packet counts, errors, retransmissions from `/proc/net/snmp`.
- **NetStat**: Extended TCP stats (syncookies, listen drops) from `/proc/net/netstat`.
- **GPU Stats**: Real-time NVIDIA GPU statistics (Utilization, Memory, Temp, Power) using `nvidia-smi`.
- **Zero Dependencies**: Uses only the Go standard library (and `nvidia-smi` for GPU).

## Usage

### Running the CLI (SysMetric Tool)

If you want to use the included terminal-based system monitor:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/avirooppal/gosysutil.git
   cd gosysutil
   ```

2. **Run directly:**
   ```bash
   go run ./cli
   ```

   **Or build and run:**
   ```bash
   go build -o sysmon ./cli
   ./sysmon
   ```

### Running the Backend API

The project includes a plug-and-play HTTP backend that exposes system metrics as JSON:

1. **Configure the environment (Optional):**
   Copy the example environment file and customize it:
   ```bash
   cp .env.example .env
   ```
   *You can set the `PORT` variable in the `.env` file.*

2. **Run the API server:**
   ```bash
   go run ./cmd/api
   ```
   *By default, the server runs on port `5001`. You can change this by setting the `PORT` environment variable or editing the `.env` file.*

3. **Endpoints:**
   - `GET /api/cpu`: CPU statistics
   - `GET /api/disk`: Disk I/O statistics
   - `GET /api/memory`: Memory usage statistics
   - `GET /api/network`: Network interface statistics
   - `GET /api/process`: Process list
   - `GET /api/all`: All-in-one system overview
   - `GET /api/loadavg`: Load average (1m, 5m, 15m)
   - `GET /api/uptime`: System uptime with formatted output
   - `GET /api/topcpu`: Top 5 CPU-consuming processes
   - `GET /api/topram`: Top 5 memory-consuming processes
   - `GET /api/steal`: IO Wait and Steal time (VPS metrics)
   - `GET /api/sockstat`: Socket statistics (TCP/UDP connections)
   - `GET /api/filenr`: File descriptor usage
   - `GET /api/pressure`: PSI (Pressure Stall Information) for CPU/Memory/IO
   - `GET /api/vmstat`: Virtual memory stats (page faults, swap, OOM)
   - `GET /api/snmp`: SNMP network stats (IP/TCP/UDP counters)
   - `GET /api/netstat`: Extended network stats (syncookies, listen drops)
   - `GET /api/gpu`: NVIDIA GPU statistics

3. **Documentation:**
   A Postman collection is available at [docs/postman_collection.json](file:///c:/Users/aviroop/Desktop/gosysutil/docs/postman_collection.json).

## Library Usage

Import the package:

```go
import "github.com/avirooppal/gosysutil/monitor"
```

### Get All Stats

```go
stats, err := monitor.GetSystemStats()
if err != nil {
    panic(err)
}

fmt.Printf("CPU User: %d\n", stats.CPU.User)
fmt.Printf("Mem Used: %d bytes\n", stats.Memory.Used)
```

### Get Specific Stats

You can also call individual functions if you only need specific metrics. Note that these are in their own packages:

```go
import (
    "github.com/avirooppal/gosysutil/cpu"
    "github.com/avirooppal/gosysutil/memory"
    "github.com/avirooppal/gosysutil/disk"
    "github.com/avirooppal/gosysutil/network"
)

// CPU
c, err := cpu.GetCPU()

// Memory
m, err := memory.GetMemory()

// Disk
d, err := disk.GetDisk()

// Network
n, err := network.GetNetwork()
```

## Structures

### CPUStats
Contains fields like `User`, `System`, `Idle`, `Iowait`, `Total`.

### MemoryStats
Contains `Total`, `Used`, `Free`, `Buffers`, `Cached`, `SwapTotal`, `SwapUsed`.
*Note: `Used` is calculated as `Total - Free - Buffers - Cached`.*

### DiskStats
Contains `Name` (e.g., "sda"), `ReadsCompleted`, `WritesCompleted`, `IoTime`, etc.

### NetworkStats
Contains `Name` (e.g., "eth0"), `RxBytes`, `TxBytes`, `RxPackets`, `TxPackets`, etc.

## Compatibility

- **Linux Only**: This package primarily relies on the `/proc` filesystem which is specific to Linux kernels.
   - *Note: GPU statistics work on both Linux and Windows as long as `nvidia-smi` is in the PATH.*
- It is designed to work on bare metal or VMs.
