// +build windows

package system

// LoadAvg represents system load averages
type LoadAvg struct {
	Load1  float64 `json:"load_1m"`
	Load5  float64 `json:"load_5m"`
	Load15 float64 `json:"load_15m"`
}

// UptimeStats represents system uptime information
type UptimeStats struct {
	Uptime   float64 `json:"uptime_seconds"`
	IdleTime float64 `json:"idle_time_seconds"`
}

// StealIOStats represents CPU steal and IO wait times (VPS specific)
type StealIOStats struct {
	StealPercent  float64 `json:"steal_percent"`
	IOWaitPercent float64 `json:"iowait_percent"`
}

// GetLoadAvg returns mock load average statistics for Windows
func GetLoadAvg() (*LoadAvg, error) {
	return &LoadAvg{
		Load1:  0.5,
		Load5:  0.3,
		Load15: 0.2,
	}, nil
}

// GetUptime returns mock uptime statistics for Windows
func GetUptime() (*UptimeStats, error) {
	return &UptimeStats{
		Uptime:   86400.0, // 1 day in seconds
		IdleTime: 43200.0, // Half day idle
	}, nil
}

// GetStealIOWait returns mock steal/iowait statistics for Windows
// On Windows these metrics are not applicable (no hypervisor steal time)
func GetStealIOWait() (*StealIOStats, error) {
	return &StealIOStats{
		StealPercent:  0.0,
		IOWaitPercent: 0.0,
	}, nil
}
