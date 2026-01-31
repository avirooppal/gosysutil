package gpu

import (
	"encoding/csv"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// GPUStats represents statistics for a single GPU
type GPUStats struct {
	Index       int     `json:"index"`
	Name        string  `json:"name"`
	UUID        string  `json:"uuid"`
	UtilGPU     float64 `json:"util_gpu_percent"`
	UtilMemory  float64 `json:"util_memory_percent"`
	MemoryTotal uint64  `json:"memory_total_mb"` // keeping as MB usually from smi
	MemoryFree  uint64  `json:"memory_free_mb"`
	MemoryUsed  uint64  `json:"memory_used_mb"`
	Temperature float64 `json:"temperature_c"`
	PowerDraw   float64 `json:"power_draw_w"`
	PowerLimit  float64 `json:"power_limit_w"`
}

// GetGPUInfo returns statistics for all detected NVIDIA GPUs
func GetGPUInfo() ([]GPUStats, error) {
	// Check if nvidia-smi is available
	_, err := exec.LookPath("nvidia-smi")
	if err != nil {
		return nil, fmt.Errorf("nvidia-smi not found: %v", err)
	}

	// Run nvidia-smi query
	// We use CSV format for easier parsing
	// Query fields: index, name, uuid, utilization.gpu, utilization.memory, memory.total, memory.free, memory.used, temperature.gpu, power.draw, power.limit
	cmd := exec.Command("nvidia-smi",
		"--query-gpu=index,name,uuid,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used,temperature.gpu,power.draw,power.limit",
		"--format=csv,noheader,nounits")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute nvidia-smi: %v", err)
	}

	r := csv.NewReader(strings.NewReader(string(output)))
	r.TrimLeadingSpace = true
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse nvidia-smi output: %v", err)
	}

	var stats []GPUStats
	for _, record := range records {
		if len(record) < 11 {
			continue
		}

		s := GPUStats{
			Name: record[1],
			UUID: record[2],
		}

		s.Index, _ = strconv.Atoi(record[0])
		s.UtilGPU, _ = strconv.ParseFloat(record[3], 64)
		s.UtilMemory, _ = strconv.ParseFloat(record[4], 64)
		s.MemoryTotal, _ = strconv.ParseUint(record[5], 10, 64)
		s.MemoryFree, _ = strconv.ParseUint(record[6], 10, 64)
		s.MemoryUsed, _ = strconv.ParseUint(record[7], 10, 64)
		s.Temperature, _ = strconv.ParseFloat(record[8], 64)
		s.PowerDraw, _ = strconv.ParseFloat(record[9], 64)
		s.PowerLimit, _ = strconv.ParseFloat(record[10], 64)

		stats = append(stats, s)
	}

	return stats, nil
}
