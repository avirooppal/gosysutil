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

// HandleCPU returns CPU statistics
func HandleCPU(w http.ResponseWriter, r *http.Request) {
	stats, err := cpu.GetCPU()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, stats)
}

// HandleDisk returns Disk statistics
func HandleDisk(w http.ResponseWriter, r *http.Request) {
	stats, err := disk.GetDisk()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, stats)
}

// HandleMemory returns Memory statistics
func HandleMemory(w http.ResponseWriter, r *http.Request) {
	stats, err := memory.GetMemory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, stats)
}

// HandleNetwork returns Network statistics
func HandleNetwork(w http.ResponseWriter, r *http.Request) {
	stats, err := network.GetNetwork()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, stats)
}

// HandleProcess returns Process statistics
func HandleProcess(w http.ResponseWriter, r *http.Request) {
	stats, err := process.GetProcesses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, stats)
}

// HandleAll returns all system statistics
func HandleAll(w http.ResponseWriter, r *http.Request) {
	stats, err := monitor.GetSystemStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, stats)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
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
