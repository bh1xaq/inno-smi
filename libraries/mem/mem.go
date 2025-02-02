package mem

import (
	"bufio"
	"errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const KERNEL_DEBUG_INNO_MEMINFO = "/sys/kernel/debug/inno/mem_info"
const KERNEL_DEBUG_INNO_MEMUSAGE = "/sys/kernel/debug/inno/mem_usage"

type MemInfo struct {
	PID       int
	Name      string
	AllocSize int64
}
type MemUsage struct {
	HeapUsage     string
	TotalSize     int64
	FreeSize      int64
	MaxBucketSize int64
	MinBucketSize int64
	BucketCount   int64
}

func parseMemInfo() ([]MemInfo, error) {
	file, err := os.Open(KERNEL_DEBUG_INNO_MEMINFO)
	if err != nil {
		return nil, errors.New("无法获取信息，请确认使用 root 权限，或 sudo 运行")
	}
	defer file.Close()

	var infos []MemInfo
	scanner := bufio.NewScanner(file)
	memInfoRegex := regexp.MustCompile(`Pid\s+(\d+)\s+Name\s+(.+?)\s+Alloc size\s+(\d+) \[K\]`)

	for scanner.Scan() {
		line := scanner.Text()
		if matches := memInfoRegex.FindStringSubmatch(line); matches != nil {
			pid, _ := strconv.Atoi(matches[1])
			allocSize, _ := strconv.ParseInt(matches[3], 10, 64)
			infos = append(infos, MemInfo{
				PID:       pid,
				Name:      strings.TrimSpace(matches[2]),
				AllocSize: allocSize,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return infos, nil
}

func parseMemUsage() ([]MemUsage, error) {
	file, err := os.Open(KERNEL_DEBUG_INNO_MEMUSAGE)
	if err != nil {
		return nil, errors.New("无法获取信息，请确认使用 root 权限，或 sudo 运行")
	}
	defer file.Close()

	var usages []MemUsage
	scanner := bufio.NewScanner(file)
	firstLineSkipped := false
	for scanner.Scan() {
		line := scanner.Text()
		if !firstLineSkipped {
			firstLineSkipped = true
			continue
		}
		fields := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(line), -1)
		if len(fields) != 6 {
			continue
		}
		totalSize, _ := strconv.ParseInt(fields[1], 10, 64)
		freeSize, _ := strconv.ParseInt(fields[2], 10, 64)
		maxBucketSize, _ := strconv.ParseInt(fields[3], 10, 64)
		minBucketSize, _ := strconv.ParseInt(fields[4], 10, 64)
		bucketCount, _ := strconv.ParseInt(fields[5], 10, 64)
		usages = append(usages, MemUsage{
			HeapUsage:     fields[0],
			TotalSize:     totalSize,
			FreeSize:      freeSize,
			MaxBucketSize: maxBucketSize,
			MinBucketSize: minBucketSize,
			BucketCount:   bucketCount,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return usages, nil
}

func GetInfoToTable() string {
	memInfo, err := parseMemInfo()
	if err != nil {
		return "ERROR：无法获取 MemInfo 信息（ " + err.Error() + " )"
	}
	t := table.Table{}
	t.SetTitle("Inno GPU Memory Info")
	t.AppendHeader(table.Row{"PID", "Name", "Alloc size[K]"})
	for _, info := range memInfo {
		t.AppendRow(table.Row{info.PID, info.Name, info.AllocSize})
	}
	return t.Render()
}

func GetUsageToTable() string {
	usages, err := parseMemUsage()
	if err != nil {
		return "ERROR：无法获取 MemUsage 信息（ " + err.Error() + " )"
	}
	t := table.Table{}
	t.SetTitle("Inno GPU Memory Usage")
	t.AppendHeader(table.Row{"HEAP_USAGE", "TOTAL_SIZE", "FREE_SIZE", "MAX_BUCKET_SIZE", "MIN_BUCKET_SIZE", "BUCKET_COUNT"})
	for _, usage := range usages {
		t.AppendRow(table.Row{usage.HeapUsage, usage.TotalSize, usage.FreeSize, usage.MaxBucketSize, usage.MinBucketSize, usage.BucketCount})
	}
	return t.Render()
}
