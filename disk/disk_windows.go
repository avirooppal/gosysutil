// +build windows

package disk

// DiskStats represents disk I/O statistics
type DiskStats struct {
	Name            string
	ReadsCompleted  uint64
	ReadsMerged     uint64
	SectorsRead     uint64
	ReadTime        uint64
	WritesCompleted uint64
	WritesMerged    uint64
	SectorsWritten  uint64
	WriteTime       uint64
	IoInProgress    uint64
	IoTime          uint64
	WeightedIoTime  uint64
}

// GetDisk returns mock disk statistics for Windows
func GetDisk() ([]DiskStats, error) {
	return []DiskStats{
		{Name: "C:", ReadsCompleted: 5000, WritesCompleted: 2000},
	}, nil
}
