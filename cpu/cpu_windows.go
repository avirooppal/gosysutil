// +build windows

package cpu

import (
	"fmt"
	"time"
)

// CPUStats represents CPU statistics
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

// GetCPU returns mock CPU statistics for Windows
func GetCPU() (*CPUStats, error) {
	// Mock data for Windows to allow the app to run and show something
	return &CPUStats{
		User:   1000,
		System: 500,
		Idle:   8500,
		Total:  10000,
	}, nil
}

// GetCPUUsage calculates the CPU usage percentage over a 200ms interval
func GetCPUUsage() (float64, error) {
	// Mock usage for Windows
	return 12.5, nil
}
	
