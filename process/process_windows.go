// +build windows

package process

// Process represents a single process
type Process struct {
	PID     int
	PPID    int
	Name    string
	State   string
	RSS     uint64
	Utime   uint64
	Stime   uint64
	Cmdline string
}

// GetProcesses returns mock process statistics for Windows
func GetProcesses() ([]Process, error) {
	return []Process{
		{PID: 1, Name: "System", State: "Running", Cmdline: "System"},
	}, nil
}
