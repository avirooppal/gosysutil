package api

import (
	"encoding/json"
	"net/http"

	"github.com/avirooppal/gosysutil/cpu"
	"github.com/avirooppal/gosysutil/disk"
	"github.com/avirooppal/gosysutil/memory"
	"github.com/avirooppal/gosysutil/monitor"
	"github.com/avirooppal/gosysutil/network"
	"github.com/avirooppal/gosysutil/process"
)

// HandleCPU returns CPU statistics with usage percentage
func HandleCPU(w http.ResponseWriter, r *http.Request) {
	rawStats, err := cpu.GetCPU()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usage, err := cpu.GetCPUUsage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"usage_percent": fmt.Sprintf("%.2f%%", usage),
		"raw_stats":     rawStats,
	}
	respondWithJSON(w, http.StatusOK, response)
}

// HandleDisk returns Disk statistics with human-readable sizes
func HandleDisk(w http.ResponseWriter, r *http.Request) {
	stats, err := disk.GetDisk()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}

// HandleMemory returns Memory statistics with human-readable sizes
func HandleMemory(w http.ResponseWriter, r *http.Request) {
	stats, err := memory.GetMemory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"total":    formatBytes(stats.Total),
		"used":     formatBytes(stats.Used),
		"free":     formatBytes(stats.Free),
		"buffers":  formatBytes(stats.Buffers),
		"cached":   formatBytes(stats.Cached),
		"raw_stats": stats,
	}
	respondWithJSON(w, http.StatusOK, response)
}

// HandleNetwork returns Network statistics with human-readable sizes
func HandleNetwork(w http.ResponseWriter, r *http.Request) {
	stats, err := network.GetNetwork()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type ReadableNetwork struct {
		Name      string
		RxBytes   string
		TxBytes   string
		RxPackets uint64
		TxPackets uint64
	}

	var readable []ReadableNetwork
	for _, s := range stats {
		readable = append(readable, ReadableNetwork{
			Name:      s.Name,
			RxBytes:   formatBytes(s.RxBytes),
			TxBytes:   formatBytes(s.TxBytes),
			RxPackets: s.RxPackets,
			TxPackets: s.TxPackets,
		})
	}

	respondWithJSON(w, http.StatusOK, readable)
}

// HandleProcess returns Process statistics with human-readable RSS
func HandleProcess(w http.ResponseWriter, r *http.Request) {
	stats, err := process.GetProcesses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type ReadableProcess struct {
		PID     int
		Name    string
		Memory  string
		Cmdline string
	}

	var readable []ReadableProcess
	for _, p := range stats {
		readable = append(readable, ReadableProcess{
			PID:     p.PID,
			Name:    p.Name,
			Memory:  formatBytes(p.RSS),
			Cmdline: p.Cmdline,
		})
	}

	respondWithJSON(w, http.StatusOK, readable)
}

// HandleAll returns all system statistics in a readable format
func HandleAll(w http.ResponseWriter, r *http.Request) {
	// For "All", we'll just return the transformed versions of each sub-metric
	// However, monitor.GetSystemStats() returns raw.
	// Let's just collect them manually here for simplicity to reuse handlers or logic.
	
	cpuUsage, _ := cpu.GetCPUUsage()
	memStats, _ := memory.GetMemory()
	
	response := map[string]interface{}{
		"cpu": map[string]interface{}{
			"usage": fmt.Sprintf("%.2f%%", cpuUsage),
		},
		"memory": map[string]interface{}{
			"total": formatBytes(memStats.Total),
			"used":  formatBytes(memStats.Used),
			"free":  formatBytes(memStats.Free),
		},
	}

	respondWithJSON(w, http.StatusOK, response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.MarshalIndent(payload, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// RegisterRoutes registers the API routes to the given multiplexer
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/cpu", HandleCPU)
	mux.HandleFunc("/api/disk", HandleDisk)
	mux.HandleFunc("/api/memory", HandleMemory)
	mux.HandleFunc("/api/network", HandleNetwork)
	mux.HandleFunc("/api/process", HandleProcess)
	mux.HandleFunc("/api/all", HandleAll)
}
