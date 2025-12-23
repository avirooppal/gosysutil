// +build linux

package cpu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// CPUStats represents CPU statistics from /proc/stat
type CPUStats struct {
	User      uint64
	Nice      uint64
	System    uint64
	Idle      uint64
	Iowait    uint64
	Irq       uint64
	Softirq   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
	Total     uint64
}

// GetCPU returns the aggregate CPU statistics from /proc/stat
func GetCPU() (*CPUStats, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "cpu" {
			return parseCPULine(fields)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("cpu line not found in /proc/stat")
}

func parseCPULine(fields []string) (*CPUStats, error) {
	if len(fields) < 8 {
		return nil, fmt.Errorf("insufficient fields in cpu line")
	}

	stats := &CPUStats{}
	var err error

	stats.User, err = strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return nil, err
	}
	stats.Nice, err = strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return nil, err
	}
	stats.System, err = strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return nil, err
	}
	stats.Idle, err = strconv.ParseUint(fields[4], 10, 64)
	if err != nil {
		return nil, err
	}
	stats.Iowait, err = strconv.ParseUint(fields[5], 10, 64)
	if err != nil {
		return nil, err
	}
	stats.Irq, err = strconv.ParseUint(fields[6], 10, 64)
	if err != nil {
		return nil, err
	}
	stats.Softirq, err = strconv.ParseUint(fields[7], 10, 64)
	if err != nil {
		return nil, err
	}
    // Newer kernels have steal, guest, guest_nice
    if len(fields) > 8 {
        stats.Steal, _ = strconv.ParseUint(fields[8], 10, 64)
    }
    if len(fields) > 9 {
        stats.Guest, _ = strconv.ParseUint(fields[9], 10, 64)
    }
    if len(fields) > 10 {
        stats.GuestNice, _ = strconv.ParseUint(fields[10], 10, 64)
    }

	stats.Total = stats.User + stats.Nice + stats.System + stats.Idle +
		stats.Iowait + stats.Irq + stats.Softirq + stats.Steal +
		stats.Guest + stats.GuestNice

	return stats, nil
}

// CPUUsage represents more detailed CPU usage percentages
type CPUUsage struct {
	TotalPercent  float64 `json:"total_percent"`
	UserPercent   float64 `json:"user_percent"`
	SystemPercent float64 `json:"system_percent"`
	IdlePercent   float64 `json:"idle_percent"`
}

// GetCPUUsage calculates the CPU usage statistics over a 500ms interval
func GetCPUUsage() (*CPUUsage, error) {
	s1, err := GetCPU()
	if err != nil {
		return nil, err
	}

	time.Sleep(500 * time.Millisecond)

	s2, err := GetCPU()
	if err != nil {
		return nil, err
	}

	totalDelta := s2.Total - s1.Total
	if totalDelta == 0 {
		return &CPUUsage{}, nil
	}

	userDelta := s2.User - s1.User
	systemDelta := s2.System - s1.System
	idleDelta := s2.Idle - s1.Idle

	return &CPUUsage{
		TotalPercent:  float64(totalDelta-idleDelta) / float64(totalDelta) * 100,
		UserPercent:   float64(userDelta) / float64(totalDelta) * 100,
		SystemPercent: float64(systemDelta) / float64(totalDelta) * 100,
		IdlePercent:   float64(idleDelta) / float64(totalDelta) * 100,
	}, nil
}
