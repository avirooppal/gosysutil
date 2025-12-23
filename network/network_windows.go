// +build windows

package network

// NetworkStats represents network interface statistics
type NetworkStats struct {
	Name        string
	RxBytes     uint64
	RxPackets   uint64
	RxErrors    uint64
	RxDropped   uint64
	TxBytes     uint64
	TxPackets   uint64
	TxErrors    uint64
	TxDropped   uint64
}

// GetNetwork returns mock network statistics for Windows
func GetNetwork() ([]NetworkStats, error) {
	return []NetworkStats{
		{Name: "Ethernet", RxBytes: 1000000, TxBytes: 500000},
	}, nil
}
