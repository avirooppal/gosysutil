// +build linux

package network

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// NetworkStats represents network interface statistics from /proc/net/dev
type NetworkStats struct {
	Name        string
	RxBytes     uint64
	RxPackets   uint64
	RxErrors    uint64
	RxDropped   uint64
	TxBytes     uint64
	TxPackets   uint64
	TxErrors    uint64
	TxDropped   uint64
}

// GetNetwork returns network statistics for all interfaces in /proc/net/dev
func GetNetwork() ([]NetworkStats, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var netStats []NetworkStats
	scanner := bufio.NewScanner(file)
    
    // Skip first 2 lines (header)
    if scanner.Scan() { // Line 1
        if scanner.Scan() { // Line 2
             // Good, filtered headers
        }
    }

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
        if len(parts) != 2 {
            continue
        }
        
        name := strings.TrimSpace(parts[0])
        fields := strings.Fields(parts[1])
        
        if len(fields) < 16 {
            // Need at least up to tx_compressed usually, or just check basic 8 for RX/TX
            // Standard is 16 fields
            if len(fields) < 8 {
                 continue
            }
        }
        
        stats := NetworkStats{
            Name: name,
        }
        
        // RX
        stats.RxBytes, _ = strconv.ParseUint(fields[0], 10, 64)
        stats.RxPackets, _ = strconv.ParseUint(fields[1], 10, 64)
        stats.RxErrors, _ = strconv.ParseUint(fields[2], 10, 64)
        stats.RxDropped, _ = strconv.ParseUint(fields[3], 10, 64)
        // 4,5,6,7 are fifo, frame, compressed, multicast
        
        // TX starts at index 8
        if len(fields) >= 16 {
             stats.TxBytes, _ = strconv.ParseUint(fields[8], 10, 64)
             stats.TxPackets, _ = strconv.ParseUint(fields[9], 10, 64)
             stats.TxErrors, _ = strconv.ParseUint(fields[10], 10, 64)
             stats.TxDropped, _ = strconv.ParseUint(fields[11], 10, 64)
        }

		netStats = append(netStats, stats)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return netStats, nil
}
