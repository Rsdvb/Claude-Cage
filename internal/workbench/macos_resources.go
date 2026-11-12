package workbench

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"al.essio.dev/pkg/shellescape"
	"github.com/docker/go-units"
	"github.com/sylvinus/agent-vm/internal/cage"
	"github.com/sylvinus/agent-vm/internal/vm"
	"go.yaml.in/yaml/v3"
)

//go:embed macos_resources.py
var macOSResourceProbe string

type macOSResourceSample struct {
	Version    int       `json:"version"`
	Total      uint64    `json:"total_bytes"`
	Available  uint64    `json:"available_bytes"`
	Free       uint64    `json:"free_bytes"`
	Inactive   uint64    `json:"inactive_bytes"`
	Wired      uint64    `json:"wired_bytes"`
	Compressed uint64    `json:"compressed_bytes"`
	SwapTotal  *uint64   `json:"swap_total_bytes,omitempty"`
	SwapFree   *uint64   `json:"swap_free_bytes,omitempty"`
	CPUs       int       `json:"cpus"`
	CPUTicks   [4]uint32 `json:"cpu_ticks"`
	Boot       int64     `json:"boot_seconds"`
	Pressure   *uint32   `json:"pressure_level,omitempty"`
}

type macOSResourceObserver struct {
	mu         sync.Mutex
	previous   macOSResourceSample
	observedAt time.Time
}

func parseMacOSResourceSample(data []byte, memory uint64, cpus int) (macOSResourceSample, error) {
	var sample macOSResourceSample
	if json.Unmarshal(data, &sample) != nil || sample.Version != 1 || sample.Total != memory || sample.CPUs != cpus || sample.Boot <= 0 ||
		sample.Available > sample.Total || sample.Free > sample.Total || sample.Inactive > sample.Total-sample.Free ||
		sample.Available != sample.Free+sample.Inactive || sample.Wired > sample.Total || sample.Compressed > sample.Total ||
		(sample.SwapTotal == nil) != (sample.SwapFree == nil) || sample.SwapTotal != nil && *sample.SwapFree > *sample.SwapTotal {
		return sample, errors.New("macOS 客机返回的资源用量不完整或与启动配置不一致")
	}
	return sample, nil
}

func (observer *macOSResourceObserver) observation(sample macOSResourceSample, now time.Time) map[string]any {
	var percent *float64
	previous := observer.previous
	if previous.Boot == sample.Boot && previous.CPUs == sample.CPUs && now.Sub(observer.observedAt) <= 10*time.Second {
		var total, idle uint64
		valid := true
		for index, value := range sample.CPUTicks {
			if value < previous.CPUTicks[index] {
				valid = false
				break
			}
			delta := uint64(value - previous.CPUTicks[index])
			total += delta
			if index == 2 {
				idle = delta
			}
		}
		if valid && total > 0 {
			value := 100 * float64(total-idle) / float64(total)
			percent = &value
		}
	}
	observer.previous, observer.observedAt = sample, now
	result := map[string]any{"version": 1, "observed_at": now, "cpus": sample.CPUs, "cpu_percent": percent, "guest": sample}
	if sample.Pressure != nil {
		result["pressure"] = *sample.Pressure > 1
	}
	return result
}

func (s *Server) readMacOSResources(ctx context.Context, memory uint64, cpus int) (map[string]any, error) {
	observer := &s.macResources
	observer.mu.Lock()
	defer observer.mu.Unlock()
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()
	command := "exec " + shellescape.QuoteCommand([]string{macCLIInstall + "/bin/python3", "-I", "-c", macOSResourceProbe})
	data, err := s.guestControl(ctx, command, "macos-resources")
	if err != nil {
		return nil, errors.New("暂时无法读取 macOS 客机资源用量，请确认客机已启动并完成组件安装")
	}
	sample, err := parseMacOSResourceSample(data, memory, cpus)
	if err != nil {
		return nil, err
	}
	return observer.observation(sample, time.Now()), nil
}

func (s *Server) resourceAllocationLimits(allocation cage.ResourceAllocation) (int, int, error) {
	if s.macosGuest == nil {
		return allocation.Limits()
	}
	requirements, err := vm.ReadMacOSResourceRequirements(s.instanceDir())
	if err != nil {
		return 0, 0, err
	}
	return requirements.Limits(allocation)
}

func (s *Server) macOSResourceStatus(w http.ResponseWriter, r *http.Request) {
	requirements, err := vm.ReadMacOSResourceRequirements(s.instanceDir())
	if err != nil {
		requestError(w, err)
		return
	}
	allocation, err := cage.LoadResourceAllocation(s.configPath)
	if err != nil {
		requestError(w, err)
		return
	}
	memory, cpus, allocationErr := requirements.Limits(allocation)
	if allocationErr != nil && !allocation.Automatic {
		memory, cpus = allocation.MemoryGiB, allocation.CPUs
	}
	recommendedMemory, recommendedCPUs, _ := requirements.Limits(cage.ResourceAllocation{Automatic: true})
	var configured struct {
		Memory string `yaml:"memory"`
		CPUs   int    `yaml:"cpus"`
	}
	data, err := os.ReadFile(filepath.Join(s.instanceDir(), "lima.yaml"))
	if err != nil || yaml.Unmarshal(data, &configured) != nil {
		requestError(w, errors.New("无法读取 macOS 客机启动配置"))
		return
	}
	maximum, err := units.RAMInBytes(configured.Memory)
	if err != nil || maximum < 1<<30 || configured.CPUs < 1 {
		requestError(w, errors.New("macOS 客机启动资源配置无效"))
		return
	}
	running := s.vmRunning()
	var observation map[string]any
	if running {
		observation, err = s.readMacOSResources(r.Context(), uint64(maximum), configured.CPUs)
		allocationErr = errors.Join(allocationErr, err)
	}
	resourceError := ""
	if allocationErr != nil {
		resourceError = allocationErr.Error()
	}
	maxMemory, maxCPUs := cage.ResourceLimits()
	jsonReply(w, map[string]any{"runtime": "vz-macos", "memory_policy": "fixed", "automatic_reclaim": false,
		"running": running, "automatic": allocation.Automatic, "memory_gib": memory, "cpus": cpus,
		"min_memory_gib": requirements.MemoryGiB, "min_cpus": requirements.CPUs, "minimum_requirements": requirements,
		"max_memory_gib": maxMemory, "max_cpus": maxCPUs, "recommended_memory_bytes": uint64(recommendedMemory) << 30, "recommended_cpus": recommendedCPUs,
		"configured_memory_bytes": maximum, "configured_cpus": configured.CPUs, "resource_error": resourceError,
		"restart_required": running && (maximum != int64(memory)<<30 || configured.CPUs != cpus), "observation": observation,
		"unavailable_observations": []string{"balloon_bytes", "oom_kills", "memory_stall_percent", "host_pressure"}})
}
