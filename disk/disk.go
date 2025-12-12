package disk

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// DiskStats represents disk I/O statistics from /proc/diskstats
type DiskStats struct {
	Name            string
	ReadsCompleted  uint64 // fields 1
	ReadsMerged     uint64 // fields 2
	SectorsRead     uint64 // fields 3
	ReadTime        uint64 // fields 4
	WritesCompleted uint64 // fields 5
	WritesMerged    uint64 // fields 6
	SectorsWritten  uint64 // fields 7
	WriteTime       uint64 // fields 8
	IoInProgress    uint64 // fields 9
	IoTime          uint64 // fields 10
	WeightedIoTime  uint64 // fields 11
}

// GetDisk returns disk I/O statistics for all disks found in /proc/diskstats
func GetDisk() ([]DiskStats, error) {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var diskStats []DiskStats
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
        // Expected at least 14 fields usually, but let's be safe with 7+ to get basics
        // Major, Minor, Name, ... stats ...
		if len(fields) < 7 {
            continue
		}
        
        name := fields[2]
        // Filter loop devices and ram devices if desired, but "plug and play" might want all.
        // Usually we want physical disks like sda, nvme0n1, vda.
        // Simple filter: skip if name starts with "loop" or "ram"
        if strings.HasPrefix(name, "loop") || strings.HasPrefix(name, "ram") {
            continue
        }
        
        // Also skip partitions (like sda1) to avoid double counting if user sums
        // A common heuristic is to check if it ends in a number, but nvme0n1 is a disk...
        // Let's just include everything for now, or let user filter.
        // Actually, for "plug and play", filtering partitions is usually helpful but risky if we miss something.
        // Let's keep it raw but maybe exclude obvious partitions like sda1 if sda exists? 
        // No, standard tools usually just dump everything or filter by type.
        // Let's minimal filter: loop/ram/sr (optical)
         if strings.HasPrefix(name, "sr") {
            continue
        }

		stats := DiskStats{
            Name: name,
        }
        
        // Fields start at index 3 for the first stat
        if len(fields) > 3 { stats.ReadsCompleted, _ = strconv.ParseUint(fields[3], 10, 64) }
        if len(fields) > 4 { stats.ReadsMerged, _ = strconv.ParseUint(fields[4], 10, 64) }
        if len(fields) > 5 { stats.SectorsRead, _ = strconv.ParseUint(fields[5], 10, 64) }
        if len(fields) > 6 { stats.ReadTime, _ = strconv.ParseUint(fields[6], 10, 64) }
        if len(fields) > 7 { stats.WritesCompleted, _ = strconv.ParseUint(fields[7], 10, 64) }
        if len(fields) > 8 { stats.WritesMerged, _ = strconv.ParseUint(fields[8], 10, 64) }
        if len(fields) > 9 { stats.SectorsWritten, _ = strconv.ParseUint(fields[9], 10, 64) }
        if len(fields) > 10 { stats.WriteTime, _ = strconv.ParseUint(fields[10], 10, 64) }
        if len(fields) > 11 { stats.IoInProgress, _ = strconv.ParseUint(fields[11], 10, 64) }
        if len(fields) > 12 { stats.IoTime, _ = strconv.ParseUint(fields[12], 10, 64) }
        if len(fields) > 13 { stats.WeightedIoTime, _ = strconv.ParseUint(fields[13], 10, 64) }

		diskStats = append(diskStats, stats)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return diskStats, nil
}
