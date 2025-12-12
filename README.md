# procmon

Lightweight, dependency-free Go package for monitoring Linux system statistics.
It is designed to be "plug and play", offering a simple API to retrieve CPU, Memory, Disk, and Network usage.

## Features

- **CPU**: Aggregated CPU usage stats (User, System, Idle, Iowait, etc.) from `/proc/stat`.
- **Memory**: Total, Used, Free, Buffers, Cached, Swap stats from `/proc/meminfo`.
- **Disk**: I/O statistics (Reads, Writes, IO Time) for physical disks from `/proc/diskstats`.
- **Network**: Traffic statistics (RX/TX bytes, packets, drops) for network interfaces from `/proc/net/dev`.
- **Zero Dependencies**: Uses only the Go standard library.

## Usage

Import the package:

```go
import "github.com/user/server-moni/pkg/procmon"
```

### Get All Stats

```go
stats, err := procmon.GetSystemStats()
if err != nil {
    panic(err)
}

fmt.Printf("CPU User: %d\n", stats.CPU.User)
fmt.Printf("Mem Used: %d bytes\n", stats.Memory.Used)
```

### Get Specific Stats

You can also call individual functions if you only need specific metrics:

```go
// CPU
cpu, err := procmon.GetCPU()

// Memory
mem, err := procmon.GetMemory()

// Disk
disks, err := procmon.GetDisk()

// Network
net, err := procmon.GetNetwork()
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

- **Linux Only**: This package relies on the `/proc` filesystem which is specific to Linux kernels.
- It is designed to work on bare metal or VMs.
