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
		{PID: 1, Name: "System", State: "Running", Cmdline: "System", RSS: 100 * 1024 * 1024, Utime: 1000, Stime: 500},
		{PID: 4, Name: "smss.exe", State: "Running", Cmdline: "smss.exe", RSS: 50 * 1024 * 1024, Utime: 500, Stime: 200},
		{PID: 8, Name: "csrss.exe", State: "Running", Cmdline: "csrss.exe", RSS: 80 * 1024 * 1024, Utime: 800, Stime: 300},
	}, nil
}

// GetTopByCPU returns mock top N processes by CPU usage for Windows
func GetTopByCPU(n int) ([]Process, error) {
	procs, _ := GetProcesses()
	if n > len(procs) {
		n = len(procs)
	}
	return procs[:n], nil
}

// GetTopByMemory returns mock top N processes by memory usage for Windows
func GetTopByMemory(n int) ([]Process, error) {
	procs, _ := GetProcesses()
	if n > len(procs) {
		n = len(procs)
	}
	return procs[:n], nil
}
