// +build linux

package memory

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// MemoryStats represents memory statistics from /proc/meminfo
type MemoryStats struct {
	Total     uint64
	Used      uint64
	Free      uint64
	Buffers   uint64
	Cached    uint64
	Active    uint64
	Inactive  uint64
	SwapTotal uint64
	SwapUsed  uint64
	SwapFree  uint64
}

// GetMemory returns memory statistics from /proc/meminfo
func GetMemory() (*MemoryStats, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := &MemoryStats{}
	scanner := bufio.NewScanner(file)
	
    // We strictly need units to be in kB as per standard, or handle parsing logic.
    // Standard /proc/meminfo is in kB.
    
    for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
        if len(parts) < 2 {
            continue
        }
        
        // Remove trailing colon from key
        key := strings.TrimSuffix(parts[0], ":")
        valString := parts[1]
        
        val, err := strconv.ParseUint(valString, 10, 64)
        if err != nil {
            continue 
        }
        
        // Convert kB to Bytes
        valBytes := val * 1024

		switch key {
		case "MemTotal":
			stats.Total = valBytes
		case "MemFree":
			stats.Free = valBytes
		case "Buffers":
			stats.Buffers = valBytes
		case "Cached":
			stats.Cached = valBytes
		case "Active":
			stats.Active = valBytes
		case "Inactive":
			stats.Inactive = valBytes
		case "SwapTotal":
			stats.SwapTotal = valBytes
		case "SwapFree":
			stats.SwapFree = valBytes
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

    stats.Used = stats.Total - stats.Free - stats.Buffers - stats.Cached
    stats.SwapUsed = stats.SwapTotal - stats.SwapFree

	return stats, nil
}
