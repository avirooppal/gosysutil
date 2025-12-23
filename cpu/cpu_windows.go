// +build windows

package cpu

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

// CPUUsage represents more detailed CPU usage percentages
type CPUUsage struct {
	TotalPercent  float64 `json:"total_percent"`
	UserPercent   float64 `json:"user_percent"`
	SystemPercent float64 `json:"system_percent"`
	IdlePercent   float64 `json:"idle_percent"`
}

// GetCPUUsage calculates mortality CPU usage statistics for Windows
func GetCPUUsage() (*CPUUsage, error) {
	// Mock usage for Windows
	return &CPUUsage{
		TotalPercent:  12.5,
		UserPercent:   8.0,
		SystemPercent: 4.5,
		IdlePercent:   87.5,
	}, nil
}
	
