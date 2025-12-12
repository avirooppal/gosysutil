package monitor

import (
	"github.com/avirooppal/gosysutil/cpu"
	"github.com/avirooppal/gosysutil/disk"
	"github.com/avirooppal/gosysutil/memory"
	"github.com/avirooppal/gosysutil/network"
)

// SystemStats aggregates all system statistics
type SystemStats struct {
	CPU     *cpu.CPUStats
	Memory  *memory.MemoryStats
	Disks   []disk.DiskStats
	Network []network.NetworkStats
}

// GetSystemStats collects all available system statistics.
// It returns partial results even if some collectors fail, collecting errors in the process if needed,
// but here we simply return the successful parts or error if critical.
func GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{}
	var err error

	stats.CPU, err = cpu.GetCPU()
	if err != nil {
		// Decide if we want to fail hard or partial. 
        // For a monitor, partial might be better, but let's just return error for now to be safe.
        // Actually, let's log/ignore for robustness? 
        // "Plug and play" usually expects valid data. 
        // If /proc/stat is missing, something is very wrong.
		return nil, err
	}

	stats.Memory, err = memory.GetMemory()
	if err != nil {
		return nil, err
	}

	stats.Disks, err = disk.GetDisk()
	if err != nil {
        // Disk stats might fail permissions or missing file on some containers?
        // But user said "no docker etc for now".
		return nil, err
	}

	stats.Network, err = network.GetNetwork()
	if err != nil {
		return nil, err
	}

	return stats, nil
}
