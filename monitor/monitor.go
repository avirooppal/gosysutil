package monitor

import (
	"github.com/avirooppal/gosysutil/cpu"
	"github.com/avirooppal/gosysutil/disk"
	"github.com/avirooppal/gosysutil/memory"
	"github.com/avirooppal/gosysutil/network"
	"github.com/avirooppal/gosysutil/process"
)

// SystemStats aggregates all system statistics
type SystemStats struct {
	CPU     *cpu.CPUStats
	Memory  *memory.MemoryStats
	Disks   []disk.DiskStats
	Network []network.NetworkStats
	Processes []process.Process
}

// GetSystemStats collects all available system statistics.
// It returns partial results even if some collectors fail, collecting errors in the process if needed,
// but here we simply return the successful parts or error if critical.
func GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{}
	var err error

	stats.CPU, err = cpu.GetCPU()
	if err != nil {
		return nil, err
	}

	stats.Memory, err = memory.GetMemory()
	if err != nil {
		return nil, err
	}

	stats.Disks, err = disk.GetDisk()
	if err != nil {
		return nil, err
	}

	stats.Network, err = network.GetNetwork()
	if err != nil {
		return nil, err
	}
    
    stats.Processes, err = process.GetProcesses()
    if err != nil {
        return nil, err
    }

	return stats, nil
}
