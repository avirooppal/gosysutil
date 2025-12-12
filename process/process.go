package process

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Process represents a single process and its metrics
type Process struct {
	PID     int
	PPID    int
	Name    string
	State   string
	RSS     uint64 // Resident Set Size in bytes
	Utime   uint64 // User time
	Stime   uint64 // System time
	Cmdline string
}

// GetProcesses returns a list of all running processes
func GetProcesses() ([]Process, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	var processes []Process
	for _, name := range names {
		// We only care about numeric directories (PIDs)
		if _, err := strconv.Atoi(name); err != nil {
			continue
		}

		p, err := parseProcess(name)
		if err != nil {
			// Process might have ended between directory list and file read
			continue
		}
		processes = append(processes, *p)
	}

	return processes, nil
}

func parseProcess(pidStr string) (*Process, error) {
	pid, _ := strconv.Atoi(pidStr)
	statPath := filepath.Join("/proc", pidStr, "stat")
	
	contents, err := ioutil.ReadFile(statPath)
	if err != nil {
		return nil, err
	}
	
	// Format is complex because Name is in parentheses and can contain spaces.
	// Example: 123 (process name) S ...
	data := string(contents)
	
	// Find the parenthesis for name
	lParen := strings.Index(data, "(")
	rParen := strings.LastIndex(data, ")")
	if lParen == -1 || rParen == -1 || rParen < lParen {
		return nil, fmt.Errorf("bad format")
	}

	name := data[lParen+1 : rParen]
	
	// Fields after the name
	rest := data[rParen+2:]
	fields := strings.Fields(rest)
	
	// Fields in `rest` start from index 2 (State) relative to the whole line
	// Check `man proc` for /proc/[pid]/stat indices
	// relative to `rest` (which starts at state):
	// 0: state (char)
	// 1: ppid (int)
	// 2: pgrp
	// ...
	// 11: utime
	// 12: stime
	// ...
	// 21: rss (pages)
	
	if len(fields) < 22 {
		return nil, fmt.Errorf("not enough fields")
	}
	
	ppid, _ := strconv.Atoi(fields[1])
	state := fields[0]
	
	utime, _ := strconv.ParseUint(fields[11], 10, 64)
	stime, _ := strconv.ParseUint(fields[12], 10, 64)
	
	rssPages, _ := strconv.ParseUint(fields[21], 10, 64)
	rss := rssPages * uint64(os.Getpagesize())

	// Read cmdline for full details
	cmdPath := filepath.Join("/proc", pidStr, "cmdline")
	cmdContent, _ := ioutil.ReadFile(cmdPath)
	cmdline := strings.ReplaceAll(string(cmdContent), "\x00", " ")
	cmdline = strings.TrimSpace(cmdline)
    
    if cmdline == "" {
        cmdline = name // fallback
    }

	return &Process{
		PID:     pid,
		PPID:    ppid,
		Name:    name,
		State:   state,
		RSS:     rss,
		Utime:   utime,
		Stime:   stime,
		Cmdline: cmdline,
	}, nil
}
