package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/avirooppal/gosysutil/cpu"
	"github.com/avirooppal/gosysutil/disk"
	"github.com/avirooppal/gosysutil/memory"
	"github.com/avirooppal/gosysutil/network"
	"github.com/avirooppal/gosysutil/process"
)

// HandleCPU returns CPU statistics with detailed usage breakdown
func HandleCPU(w http.ResponseWriter, r *http.Request) {
	usage, err := cpu.GetCPUUsage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"total_usage":  fmt.Sprintf("%.2f%%", usage.TotalPercent),
		"user_usage":   fmt.Sprintf("%.2f%%", usage.UserPercent),
		"system_usage": fmt.Sprintf("%.2f%%", usage.SystemPercent),
		"idle_usage":   fmt.Sprintf("%.2f%%", usage.IdlePercent),
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

	type ReadableDisk struct {
		Name   string `json:"name"`
		Reads  uint64 `json:"reads"`
		Writes uint64 `json:"writes"`
	}

	var readable []ReadableDisk
	for _, d := range stats {
		readable = append(readable, ReadableDisk{
			Name:   d.Name,
			Reads:  d.ReadsCompleted,
			Writes: d.WritesCompleted,
		})
	}

	respondWithJSON(w, http.StatusOK, readable)
}

// HandleMemory returns Memory statistics with human-readable sizes and percentage
func HandleMemory(w http.ResponseWriter, r *http.Request) {
	stats, err := memory.GetMemory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	percentUsed := 0.0
	if stats.Total > 0 {
		percentUsed = float64(stats.Total-stats.Free) / float64(stats.Total) * 100
	}

	response := map[string]interface{}{
		"total_memory":   formatBytes(stats.Total),
		"used_memory":    formatBytes(stats.Used),
		"free_memory":    formatBytes(stats.Free),
		"percent_used":   fmt.Sprintf("%.2f%%", percentUsed),
		"swap_total":     formatBytes(stats.SwapTotal),
		"swap_used":      formatBytes(stats.SwapUsed),
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
		Name      string `json:"interface"`
		Rx        string `json:"received"`
		Tx        string `json:"sent"`
		RxPackets uint64 `json:"rx_packets"`
		TxPackets uint64 `json:"tx_packets"`
	}

	var readable []ReadableNetwork
	for _, s := range stats {
		readable = append(readable, ReadableNetwork{
			Name:      s.Name,
			Rx:        formatBytes(s.RxBytes),
			Tx:        formatBytes(s.TxBytes),
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
		PID     int    `json:"pid"`
		Name    string `json:"name"`
		Memory  string `json:"memory_usage"`
		Cmdline string `json:"command"`
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

// HandleAll returns a concise summary of all system statistics
func HandleAll(w http.ResponseWriter, r *http.Request) {
	cpuUsg, _ := cpu.GetCPUUsage()
	memStats, _ := memory.GetMemory()
	
	percentUsed := 0.0
	if memStats.Total > 0 {
		percentUsed = float64(memStats.Total-memStats.Free) / float64(memStats.Total) * 100
	}

	response := map[string]interface{}{
		"cpu": map[string]interface{}{
			"usage": fmt.Sprintf("%.2f%%", cpuUsg.TotalPercent),
			"user":  fmt.Sprintf("%.2f%%", cpuUsg.UserPercent),
			"sys":   fmt.Sprintf("%.2f%%", cpuUsg.SystemPercent),
		},
		"memory": map[string]interface{}{
			"total": formatBytes(memStats.Total),
			"used":  formatBytes(memStats.Used),
			"usage": fmt.Sprintf("%.2f%%", percentUsed),
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
