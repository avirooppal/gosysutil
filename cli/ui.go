package main

import (
	"fmt"
	"math"
    "sort"
	"strings"

	// tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	// Colors
	primaryColor = lipgloss.Color("#7D56F4") // Purple
	subColor     = lipgloss.Color("#43BF6D") // Green
	errorColor   = lipgloss.Color("#FF5F87") // Red
	grayColor    = lipgloss.Color("#626262")

	// Base Styles
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor)

	sectionStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(grayColor).
			Padding(0, 1).
			Margin(0, 1)

	labelStyle = lipgloss.NewStyle().Foreground(grayColor)
)

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error fetching stats: %v\n(Are you on Linux?)", m.err)
	}
	if m.currentStats == nil {
		return "Loading..."
	}

	// 1. Header
	header := titleStyle.Render(" SYSTEM MONITOR ")

	// 2. Content
	// Calculate CPU Usage %
	cpuUsage := 0.0
	if m.lastStats != nil && m.currentStats != nil {
		deltaTotal := float64(m.currentStats.CPU.Total - m.lastStats.CPU.Total)
		deltaIdle := float64(m.currentStats.CPU.Idle - m.lastStats.CPU.Idle)
        if deltaTotal > 0 {
		    cpuUsage = ((deltaTotal - deltaIdle) / deltaTotal) * 100
        }
	} else if m.currentStats != nil {
        // Fallback for first frame
        cpuUsage = 0.0 
    }

	// CPU View
	cpuBar := renderProgressBar(cpuUsage, 30) // Fixed width for now
	cpuSection := sectionStyle.Render(fmt.Sprintf(
		"%s\n\nUsage: %s %.1f%%",
		labelStyle.Render("CPU"),
		cpuBar,
		cpuUsage,
	))

	// Memory View
	memUsedGB := float64(m.currentStats.Memory.Used) / 1024 / 1024 / 1024
	memTotalGB := float64(m.currentStats.Memory.Total) / 1024 / 1024 / 1024
    memUsagePct := (memUsedGB / memTotalGB) * 100
    if math.IsNaN(memUsagePct) { memUsagePct = 0 }
    
    memBar := renderProgressBar(memUsagePct, 30)
	memSection := sectionStyle.Render(fmt.Sprintf(
		"%s\n\nUsage: %s %.1f%%\nFrom:  %.2f GB / %.2f GB",
		labelStyle.Render("MEMORY"),
		memBar,
        memUsagePct,
		memUsedGB,
		memTotalGB,
	))

    // Disks View (First 3)
    var diskRows []string
    for i, d := range m.currentStats.Disks {
        if i >= 3 { break }
        diskRows = append(diskRows, fmt.Sprintf("%-10s R:%d W:%d", d.Name, d.ReadsCompleted, d.WritesCompleted))
    }
    if len(diskRows) == 0 { diskRows = append(diskRows, "No disks found") }
    diskSection := sectionStyle.Render(fmt.Sprintf(
        "%s\n\n%s",
        labelStyle.Render("DISKS"),
        strings.Join(diskRows, "\n"),
    ))
    
    // Network View (First 3)
    var netRows []string
    for i, n := range m.currentStats.Network {
        if i >= 3 { break }
        // Calc speed if possible
        rxSpeed := 0.0
        txSpeed := 0.0
        if m.lastStats != nil {
             // Find matching interface in lastStats
             for _, lastN := range m.lastStats.Network {
                 if lastN.Name == n.Name {
                     rxSpeed = float64(n.RxBytes - lastN.RxBytes) // Bytes per second (since ticker is 1s)
                     txSpeed = float64(n.TxBytes - lastN.TxBytes)
                     break
                 }
             }
        }
        
        netRows = append(netRows, fmt.Sprintf("%-6s ↓ %s/s ↑ %s/s", 
            n.Name, 
            humanizeBytes(rxSpeed), 
            humanizeBytes(txSpeed),
        ))
    }
    if len(netRows) == 0 { netRows = append(netRows, "No interfaces found") }
    netSection := sectionStyle.Render(fmt.Sprintf(
        "%s\n\n%s",
        labelStyle.Render("NETWORK"),
        strings.Join(netRows, "\n"),
    ))

    // Processes View (Top 10 by RSS)
    var procRows []string
    procs := m.currentStats.Processes
    // Sort by RSS descending
    sort.Slice(procs, func(i, j int) bool {
        return procs[i].RSS > procs[j].RSS
    })
    
    // Header
    procRows = append(procRows, fmt.Sprintf("%-6s %-10s %-20s", "PID", "RSS", "CMD"))
    
    for i, p := range procs {
        if i >= 10 { break }
        // Truncate cmdline
        cmd := p.Cmdline
        if len(cmd) > 20 { cmd = cmd[:17] + "..." }
        
        procRows = append(procRows, fmt.Sprintf("%-6d %-10s %-20s", 
            p.PID, 
            humanizeBytes(float64(p.RSS)), 
            cmd,
        ))
    }
    
    procSection := sectionStyle.Render(fmt.Sprintf(
        "%s\n\n%s",
        labelStyle.Render("TOP PROCESSES (MEM)"),
        strings.Join(procRows, "\n"),
    ))

	// Layout: Top Row (CPU + Mem), Bottom Row (Disk + Net)
    // We join horizontally using lipgloss.JoinHorizontal
    topRow := lipgloss.JoinHorizontal(lipgloss.Top, cpuSection, memSection)
    bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, diskSection, netSection)

	return appStyle.Render(fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s",
		header,
		topRow,
		bottomRow,
        procSection,
		labelStyle.Render("Press 'q' or 'esc' to quit"),
	))
}

func renderProgressBar(percent float64, width int) string {
	// Simple text-based progress bar
    if percent < 0 { percent = 0 }
    if percent > 100 { percent = 100 }
    
	fullChars := int((percent / 100) * float64(width))
	emptyChars := width - fullChars
    
    // Safety check
    if fullChars < 0 { fullChars = 0 }
    if emptyChars < 0 { emptyChars = 0 }

	bar := strings.Repeat("█", fullChars) + strings.Repeat("░", emptyChars)
    
    // Colorize based on usage
    barStyle := lipgloss.NewStyle().Foreground(subColor)
    if percent > 80 {
        barStyle = barStyle.Foreground(errorColor)
    }
    
	return barStyle.Render(bar)
}

func humanizeBytes(b float64) string {
    const unit = 1024
    if b < unit {
        return fmt.Sprintf("%.0f B", b)
    }
    div, exp := int64(unit), 0
    for n := b / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", b/float64(div), "KMGTPE"[exp])
}
