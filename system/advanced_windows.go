// +build windows

package system

// SockStats represents socket statistics
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

// FileNRStats represents file descriptor stats
type FileNRStats struct {
	Allocated uint64  `json:"allocated"`
	Free      uint64  `json:"free"`
	Max       uint64  `json:"max"`
	UsedPct   float64 `json:"used_percent"`
}

// PressureStats represents PSI (Pressure Stall Information)
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

// GetSockStats returns mock socket statistics for Windows
func GetSockStats() (*SockStats, error) {
	return &SockStats{
		SocketsUsed: 150,
		TCPInUse:    45,
		TCPOrphan:   0,
		TCPTimeWait: 12,
		TCPAlloc:    50,
		TCPMem:      8,
		UDPInUse:    10,
		UDPMem:      2,
		RAWInUse:    0,
		FragInUse:   0,
		FragMem:     0,
	}, nil
}

// GetFileNRStats returns mock file descriptor statistics for Windows
func GetFileNRStats() (*FileNRStats, error) {
	return &FileNRStats{
		Allocated: 1024,
		Free:      512,
		Max:       65536,
		UsedPct:   0.78,
	}, nil
}

// GetCPUPressure returns mock CPU pressure stats for Windows
func GetCPUPressure() (*PressureStats, error) {
	return &PressureStats{
		SomeAvg10:  0.0,
		SomeAvg60:  0.0,
		SomeAvg300: 0.0,
		SomeTotal:  0,
	}, nil
}

// GetMemoryPressure returns mock memory pressure stats for Windows
func GetMemoryPressure() (*PressureStats, error) {
	return &PressureStats{
		SomeAvg10:  0.0,
		SomeAvg60:  0.0,
		SomeAvg300: 0.0,
		SomeTotal:  0,
		FullAvg10:  0.0,
		FullAvg60:  0.0,
		FullAvg300: 0.0,
		FullTotal:  0,
	}, nil
}

// GetIOPressure returns mock IO pressure stats for Windows
func GetIOPressure() (*PressureStats, error) {
	return &PressureStats{
		SomeAvg10:  0.0,
		SomeAvg60:  0.0,
		SomeAvg300: 0.0,
		SomeTotal:  0,
		FullAvg10:  0.0,
		FullAvg60:  0.0,
		FullAvg300: 0.0,
		FullTotal:  0,
	}, nil
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

// GetVMStats returns mock vmstat statistics for Windows
func GetVMStats() (*VMStats, error) {
	return &VMStats{
		PgFault:    100000,
		PgMajFault: 50,
		PgPgIn:     50000,
		PgPgOut:    30000,
		PSwpIn:     0,
		PSwpOut:    0,
		OOMKill:    0,
		NumaHit:    0,
		NumaMiss:   0,
	}, nil
}

// SNMPStats represents key metrics from /proc/net/snmp
type SNMPStats struct {
	IPInReceives      uint64 `json:"ip_in_receives"`
	IPOutRequests     uint64 `json:"ip_out_requests"`
	IPInDiscards      uint64 `json:"ip_in_discards"`
	IPOutDiscards     uint64 `json:"ip_out_discards"`
	TCPActiveOpens    uint64 `json:"tcp_active_opens"`
	TCPPassiveOpens   uint64 `json:"tcp_passive_opens"`
	TCPCurrEstab      uint64 `json:"tcp_curr_estab"`
	TCPInSegs         uint64 `json:"tcp_in_segs"`
	TCPOutSegs        uint64 `json:"tcp_out_segs"`
	TCPRetransSegs    uint64 `json:"tcp_retrans_segs"`
	TCPInErrs         uint64 `json:"tcp_in_errs"`
	TCPOutRsts        uint64 `json:"tcp_out_rsts"`
	UDPInDatagrams    uint64 `json:"udp_in_datagrams"`
	UDPOutDatagrams   uint64 `json:"udp_out_datagrams"`
	UDPInErrors       uint64 `json:"udp_in_errors"`
	UDPNoPorts        uint64 `json:"udp_no_ports"`
}

// GetSNMPStats returns mock SNMP statistics for Windows
func GetSNMPStats() (*SNMPStats, error) {
	return &SNMPStats{
		IPInReceives:    500000,
		IPOutRequests:   400000,
		IPInDiscards:    10,
		IPOutDiscards:   5,
		TCPActiveOpens:  1000,
		TCPPassiveOpens: 500,
		TCPCurrEstab:    25,
		TCPInSegs:       300000,
		TCPOutSegs:      250000,
		TCPRetransSegs:  100,
		TCPInErrs:       5,
		TCPOutRsts:      50,
		UDPInDatagrams:  50000,
		UDPOutDatagrams: 40000,
		UDPInErrors:     0,
		UDPNoPorts:      10,
	}, nil
}

// NetStatStats represents key metrics from /proc/net/netstat
type NetStatStats struct {
	TCPSyncookiesSent   uint64 `json:"tcp_syncookies_sent"`
	TCPSyncookiesRecv   uint64 `json:"tcp_syncookies_recv"`
	TCPSyncookiesFailed uint64 `json:"tcp_syncookies_failed"`
	TCPListenOverflows  uint64 `json:"tcp_listen_overflows"`
	TCPListenDrops      uint64 `json:"tcp_listen_drops"`
	TCPTimeouts         uint64 `json:"tcp_timeouts"`
	IPInOctets          uint64 `json:"ip_in_octets"`
	IPOutOctets         uint64 `json:"ip_out_octets"`
}

// GetNetStatStats returns mock netstat statistics for Windows
func GetNetStatStats() (*NetStatStats, error) {
	return &NetStatStats{
		TCPSyncookiesSent:   0,
		TCPSyncookiesRecv:   0,
		TCPSyncookiesFailed: 0,
		TCPListenOverflows:  0,
		TCPListenDrops:      0,
		TCPTimeouts:         50,
		IPInOctets:          1000000000,
		IPOutOctets:         800000000,
	}, nil
}
