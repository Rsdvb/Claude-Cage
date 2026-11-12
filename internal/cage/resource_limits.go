package cage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ResourceAllocation is separate from network configuration: changing a VM's
// capacity does not change its network authorization fingerprint.
type ResourceAllocation struct {
	Automatic bool `json:"automatic"`
	MemoryGiB int  `json:"memory_gib"`
	CPUs      int  `json:"cpus"`
}

func ResourceAllocationPath(configPath string) string {
	return filepath.Join(filepath.Dir(configPath), "resource-limits.json")
}

func (a ResourceAllocation) Limits() (memory, cpus int, err error) {
	memory, cpus = ResourceLimits()
	if a.Automatic {
		return min(memory, 4), min(cpus, 4), nil
	}
	if a.MemoryGiB < 2 || a.MemoryGiB > memory || a.CPUs < 1 || a.CPUs > cpus {
		return 0, 0, fmt.Errorf("内存上限须为 2–%d GiB，CPU 须为 1–%d 个逻辑核心", memory, cpus)
	}
	return a.MemoryGiB, a.CPUs, nil
}

func LoadResourceAllocation(configPath string) (ResourceAllocation, error) {
	a := ResourceAllocation{Automatic: true}
	b, err := os.ReadFile(ResourceAllocationPath(configPath))
	if errors.Is(err, os.ErrNotExist) {
		return a, nil
	}
	if err != nil {
		return a, err
	}
	if err = json.Unmarshal(b, &a); err != nil {
		return a, err
	}
	return a, nil
}

func ConfiguredResourceLimits(configPath string) (int, int, error) {
	a, err := LoadResourceAllocation(configPath)
	if err != nil {
		return 0, 0, err
	}
	return a.Limits()
}
