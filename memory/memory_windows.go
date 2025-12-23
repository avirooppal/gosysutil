// +build windows

package memory

// MemoryStats represents memory statistics
type MemoryStats struct {
	Total     uint64
	Used      uint64
	Free      uint64
	Buffers   uint64
	Cached    uint64
	Active    uint64
	Inactive  uint64
	SwapTotal uint64
	SwapUsed  uint64
	SwapFree  uint64
}

// GetMemory returns mock memory statistics for Windows
func GetMemory() (*MemoryStats, error) {
	return &MemoryStats{
		Total: 16 * 1024 * 1024 * 1024,
		Used:  8 * 1024 * 1024 * 1024,
		Free:  8 * 1024 * 1024 * 1024,
	}, nil
}
