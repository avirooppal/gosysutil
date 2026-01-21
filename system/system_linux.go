// +build linux

package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

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

// GetLoadAvg returns the system load averages from /proc/loadavg
func GetLoadAvg() (*LoadAvg, error) {
	file, err := os.Open("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			return nil, fmt.Errorf("unexpected format in /proc/loadavg")
		}

		load1, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			return nil, err
		}
		load5, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return nil, err
		}
		load15, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, err
		}

		return &LoadAvg{
			Load1:  load1,
			Load5:  load5,
			Load15: load15,
		}, nil
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("failed to read /proc/loadavg")
}

// GetUptime returns the system uptime from /proc/uptime
func GetUptime() (*UptimeStats, error) {
	file, err := os.Open("/proc/uptime")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			return nil, fmt.Errorf("unexpected format in /proc/uptime")
		}

		uptime, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			return nil, err
		}
		idleTime, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return nil, err
		}

		return &UptimeStats{
			Uptime:   uptime,
			IdleTime: idleTime,
		}, nil
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("failed to read /proc/uptime")
}

// GetStealIOWait returns the CPU steal and IO wait percentages
// This is particularly useful for VPS environments where CPU steal indicates
// the hypervisor is using CPU time that was allocated to this VM
func GetStealIOWait() (*StealIOStats, error) {
	readStats := func() (iowait, steal, total uint64, err error) {
		file, err := os.Open("/proc/stat")
		if err != nil {
			return 0, 0, 0, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) > 0 && fields[0] == "cpu" {
				if len(fields) < 9 {
					return 0, 0, 0, fmt.Errorf("insufficient fields in cpu line")
				}

				var vals []uint64
				for i := 1; i < len(fields) && i <= 10; i++ {
					v, _ := strconv.ParseUint(fields[i], 10, 64)
					vals = append(vals, v)
				}

				// Fields: user, nice, system, idle, iowait, irq, softirq, steal, guest, guest_nice
				iowait = vals[4]
				if len(vals) > 7 {
					steal = vals[7]
				}

				for _, v := range vals {
					total += v
				}

				return iowait, steal, total, nil
			}
		}

		return 0, 0, 0, fmt.Errorf("cpu line not found")
	}

	iowait1, steal1, total1, err := readStats()
	if err != nil {
		return nil, err
	}

	time.Sleep(500 * time.Millisecond)

	iowait2, steal2, total2, err := readStats()
	if err != nil {
		return nil, err
	}

	totalDelta := total2 - total1
	if totalDelta == 0 {
		return &StealIOStats{}, nil
	}

	iowaitDelta := iowait2 - iowait1
	stealDelta := steal2 - steal1

	return &StealIOStats{
		StealPercent:  float64(stealDelta) / float64(totalDelta) * 100,
		IOWaitPercent: float64(iowaitDelta) / float64(totalDelta) * 100,
	}, nil
}
