// +build linux

package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// SockStats represents socket statistics from /proc/net/sockstat
type SockStats struct {
	SocketsUsed int `json:"sockets_used"`
	TCPInUse    int `json:"tcp_inuse"`
	TCPOrphan   int `json:"tcp_orphan"`
	TCPTimeWait int `json:"tcp_tw"`
	TCPAlloc    int `json:"tcp_alloc"`
	TCPMem      int `json:"tcp_mem"`
	UDPInUse    int `json:"udp_inuse"`
	UDPMem      int `json:"udp_mem"`
	RAWInUse    int `json:"raw_inuse"`
	FragInUse   int `json:"frag_inuse"`
	FragMem     int `json:"frag_mem"`
}

// FileNRStats represents file descriptor stats from /proc/sys/fs/file-nr
type FileNRStats struct {
	Allocated uint64  `json:"allocated"`
	Free      uint64  `json:"free"`
	Max       uint64  `json:"max"`
	UsedPct   float64 `json:"used_percent"`
}

// PressureStats represents PSI (Pressure Stall Information) from /proc/pressure/*
type PressureStats struct {
	SomeAvg10  float64 `json:"some_avg10"`
	SomeAvg60  float64 `json:"some_avg60"`
	SomeAvg300 float64 `json:"some_avg300"`
	SomeTotal  uint64  `json:"some_total"`
	FullAvg10  float64 `json:"full_avg10"`
	FullAvg60  float64 `json:"full_avg60"`
	FullAvg300 float64 `json:"full_avg300"`
	FullTotal  uint64  `json:"full_total"`
}

// GetSockStats returns socket statistics from /proc/net/sockstat
func GetSockStats() (*SockStats, error) {
	file, err := os.Open("/proc/net/sockstat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := &SockStats{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "sockets:":
			if len(fields) >= 3 {
				stats.SocketsUsed, _ = strconv.Atoi(fields[2])
			}
		case "TCP:":
			for i := 1; i < len(fields)-1; i += 2 {
				val, _ := strconv.Atoi(fields[i+1])
				switch fields[i] {
				case "inuse":
					stats.TCPInUse = val
				case "orphan":
					stats.TCPOrphan = val
				case "tw":
					stats.TCPTimeWait = val
				case "alloc":
					stats.TCPAlloc = val
				case "mem":
					stats.TCPMem = val
				}
			}
		case "UDP:":
			for i := 1; i < len(fields)-1; i += 2 {
				val, _ := strconv.Atoi(fields[i+1])
				switch fields[i] {
				case "inuse":
					stats.UDPInUse = val
				case "mem":
					stats.UDPMem = val
				}
			}
		case "RAW:":
			if len(fields) >= 3 {
				stats.RAWInUse, _ = strconv.Atoi(fields[2])
			}
		case "FRAG:":
			for i := 1; i < len(fields)-1; i += 2 {
				val, _ := strconv.Atoi(fields[i+1])
				switch fields[i] {
				case "inuse":
					stats.FragInUse = val
				case "memory":
					stats.FragMem = val
				}
			}
		}
	}

	return stats, scanner.Err()
}

// GetFileNRStats returns file descriptor statistics from /proc/sys/fs/file-nr
func GetFileNRStats() (*FileNRStats, error) {
	data, err := os.ReadFile("/proc/sys/fs/file-nr")
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return nil, fmt.Errorf("unexpected format in /proc/sys/fs/file-nr")
	}

	allocated, _ := strconv.ParseUint(fields[0], 10, 64)
	free, _ := strconv.ParseUint(fields[1], 10, 64)
	max, _ := strconv.ParseUint(fields[2], 10, 64)

	usedPct := 0.0
	if max > 0 {
		usedPct = float64(allocated-free) / float64(max) * 100
	}

	return &FileNRStats{
		Allocated: allocated,
		Free:      free,
		Max:       max,
		UsedPct:   usedPct,
	}, nil
}

// GetPressure returns PSI stats for the given resource (cpu, memory, io)
func GetPressure(resource string) (*PressureStats, error) {
	path := fmt.Sprintf("/proc/pressure/%s", resource)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := &PressureStats{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		isFull := strings.HasPrefix(fields[0], "full")

		for _, field := range fields[1:] {
			parts := strings.Split(field, "=")
			if len(parts) != 2 {
				continue
			}

			switch parts[0] {
			case "avg10":
				val, _ := strconv.ParseFloat(parts[1], 64)
				if isFull {
					stats.FullAvg10 = val
				} else {
					stats.SomeAvg10 = val
				}
			case "avg60":
				val, _ := strconv.ParseFloat(parts[1], 64)
				if isFull {
					stats.FullAvg60 = val
				} else {
					stats.SomeAvg60 = val
				}
			case "avg300":
				val, _ := strconv.ParseFloat(parts[1], 64)
				if isFull {
					stats.FullAvg300 = val
				} else {
					stats.SomeAvg300 = val
				}
			case "total":
				val, _ := strconv.ParseUint(parts[1], 10, 64)
				if isFull {
					stats.FullTotal = val
				} else {
					stats.SomeTotal = val
				}
			}
		}
	}

	return stats, scanner.Err()
}

// GetCPUPressure returns CPU pressure stats
func GetCPUPressure() (*PressureStats, error) {
	return GetPressure("cpu")
}

// GetMemoryPressure returns memory pressure stats
func GetMemoryPressure() (*PressureStats, error) {
	return GetPressure("memory")
}

// GetIOPressure returns IO pressure stats
func GetIOPressure() (*PressureStats, error) {
	return GetPressure("io")
}

// VMStats represents key metrics from /proc/vmstat
type VMStats struct {
	PgFault       uint64 `json:"pgfault"`
	PgMajFault    uint64 `json:"pgmajfault"`
	PgPgIn        uint64 `json:"pgpgin"`
	PgPgOut       uint64 `json:"pgpgout"`
	PSwpIn        uint64 `json:"pswpin"`
	PSwpOut       uint64 `json:"pswpout"`
	OOMKill       uint64 `json:"oom_kill"`
	NumaHit       uint64 `json:"numa_hit"`
	NumaMiss      uint64 `json:"numa_miss"`
}

// GetVMStats returns virtual memory statistics from /proc/vmstat
func GetVMStats() (*VMStats, error) {
	file, err := os.Open("/proc/vmstat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := &VMStats{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			continue
		}

		val, _ := strconv.ParseUint(fields[1], 10, 64)
		switch fields[0] {
		case "pgfault":
			stats.PgFault = val
		case "pgmajfault":
			stats.PgMajFault = val
		case "pgpgin":
			stats.PgPgIn = val
		case "pgpgout":
			stats.PgPgOut = val
		case "pswpin":
			stats.PSwpIn = val
		case "pswpout":
			stats.PSwpOut = val
		case "oom_kill":
			stats.OOMKill = val
		case "numa_hit":
			stats.NumaHit = val
		case "numa_miss":
			stats.NumaMiss = val
		}
	}

	return stats, scanner.Err()
}

// SNMPStats represents key metrics from /proc/net/snmp
type SNMPStats struct {
	// IP stats
	IPInReceives      uint64 `json:"ip_in_receives"`
	IPOutRequests     uint64 `json:"ip_out_requests"`
	IPInDiscards      uint64 `json:"ip_in_discards"`
	IPOutDiscards     uint64 `json:"ip_out_discards"`
	// TCP stats
	TCPActiveOpens    uint64 `json:"tcp_active_opens"`
	TCPPassiveOpens   uint64 `json:"tcp_passive_opens"`
	TCPCurrEstab      uint64 `json:"tcp_curr_estab"`
	TCPInSegs         uint64 `json:"tcp_in_segs"`
	TCPOutSegs        uint64 `json:"tcp_out_segs"`
	TCPRetransSegs    uint64 `json:"tcp_retrans_segs"`
	TCPInErrs         uint64 `json:"tcp_in_errs"`
	TCPOutRsts        uint64 `json:"tcp_out_rsts"`
	// UDP stats
	UDPInDatagrams    uint64 `json:"udp_in_datagrams"`
	UDPOutDatagrams   uint64 `json:"udp_out_datagrams"`
	UDPInErrors       uint64 `json:"udp_in_errors"`
	UDPNoPorts        uint64 `json:"udp_no_ports"`
}

// GetSNMPStats returns SNMP statistics from /proc/net/snmp
func GetSNMPStats() (*SNMPStats, error) {
	file, err := os.Open("/proc/net/snmp")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := &SNMPStats{}
	scanner := bufio.NewScanner(file)

	var headers []string
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		prefix := strings.TrimSuffix(fields[0], ":")
		
		// First line is headers, second line is values
		if headers == nil || headers[0] != prefix {
			headers = fields
			continue
		}

		// Parse values based on prefix
		for i := 1; i < len(fields) && i < len(headers); i++ {
			val, _ := strconv.ParseUint(fields[i], 10, 64)
			key := prefix + headers[i]
			
			switch key {
			case "IpInReceives":
				stats.IPInReceives = val
			case "IpOutRequests":
				stats.IPOutRequests = val
			case "IpInDiscards":
				stats.IPInDiscards = val
			case "IpOutDiscards":
				stats.IPOutDiscards = val
			case "TcpActiveOpens":
				stats.TCPActiveOpens = val
			case "TcpPassiveOpens":
				stats.TCPPassiveOpens = val
			case "TcpCurrEstab":
				stats.TCPCurrEstab = val
			case "TcpInSegs":
				stats.TCPInSegs = val
			case "TcpOutSegs":
				stats.TCPOutSegs = val
			case "TcpRetransSegs":
				stats.TCPRetransSegs = val
			case "TcpInErrs":
				stats.TCPInErrs = val
			case "TcpOutRsts":
				stats.TCPOutRsts = val
			case "UdpInDatagrams":
				stats.UDPInDatagrams = val
			case "UdpOutDatagrams":
				stats.UDPOutDatagrams = val
			case "UdpInErrors":
				stats.UDPInErrors = val
			case "UdpNoPorts":
				stats.UDPNoPorts = val
			}
		}
		headers = nil
	}

	return stats, scanner.Err()
}

// NetStatStats represents key metrics from /proc/net/netstat
type NetStatStats struct {
	// TCP Extension stats
	TCPSyncookiesSent   uint64 `json:"tcp_syncookies_sent"`
	TCPSyncookiesRecv   uint64 `json:"tcp_syncookies_recv"`
	TCPSyncookiesFailed uint64 `json:"tcp_syncookies_failed"`
	TCPListenOverflows  uint64 `json:"tcp_listen_overflows"`
	TCPListenDrops      uint64 `json:"tcp_listen_drops"`
	TCPTimeouts         uint64 `json:"tcp_timeouts"`
	// IP Extension stats
	IPInOctets          uint64 `json:"ip_in_octets"`
	IPOutOctets         uint64 `json:"ip_out_octets"`
}

// GetNetStatStats returns extended network statistics from /proc/net/netstat
func GetNetStatStats() (*NetStatStats, error) {
	file, err := os.Open("/proc/net/netstat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := &NetStatStats{}
	scanner := bufio.NewScanner(file)

	var headers []string
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		prefix := strings.TrimSuffix(fields[0], ":")

		if headers == nil || headers[0] != prefix {
			headers = fields
			continue
		}

		for i := 1; i < len(fields) && i < len(headers); i++ {
			val, _ := strconv.ParseUint(fields[i], 10, 64)
			key := headers[i]

			switch key {
			case "SyncookiesSent":
				stats.TCPSyncookiesSent = val
			case "SyncookiesRecv":
				stats.TCPSyncookiesRecv = val
			case "SyncookiesFailed":
				stats.TCPSyncookiesFailed = val
			case "ListenOverflows":
				stats.TCPListenOverflows = val
			case "ListenDrops":
				stats.TCPListenDrops = val
			case "TCPTimeouts":
				stats.TCPTimeouts = val
			case "InOctets":
				stats.IPInOctets = val
			case "OutOctets":
				stats.IPOutOctets = val
			}
		}
		headers = nil
	}

	return stats, scanner.Err()
}
