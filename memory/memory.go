package memory

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// MemoryStats represents memory statistics from /proc/meminfo
type MemoryStats struct {
	Total     uint64 // Total usable RAM (i.e. physical RAM minus a few reserved bits and the kernel binary code)
	Used      uint64 // Calculated as Total - Free - Buffers - Cached
	Free      uint64 // The sum of LowFree+HighFree
	Buffers   uint64 // Relatively temporary storage for raw disk blocks
	Cached    uint64 // In-memory cache for files read from the disk (the pagecache)
	Active    uint64 // Memory that has been used more recently and usually not reclaimed unless absolutely necessary
	Inactive  uint64 // Memory which has been less recently used. It is more eligible to be reclaimed for other purposes
	SwapTotal uint64 // Total amount of swap space available
	SwapUsed  uint64 // Memory which has been evicted from RAM, and is temporarily on the disk
	SwapFree  uint64 // Amount of swap space that is currently unused
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
