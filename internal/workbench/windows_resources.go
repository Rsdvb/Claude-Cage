package workbench

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sylvinus/agent-vm/internal/cage"
	"github.com/sylvinus/agent-vm/internal/guestos"
	"github.com/sylvinus/agent-vm/internal/host"
)

func windowsResourceMaximum() (int, int) {
	shareMemory, shareCPUs := cage.ResourceLimits()
	managerMemory, managerCPUs := guestos.WindowsVMResourceCeiling(host.MemGiB(), host.CPUs())
	return min(shareMemory, managerMemory), min(shareCPUs, managerCPUs)
}

func windowsAllocationLimits(allocation cage.ResourceAllocation, maxMemory, maxCPUs int) (int, int, error) {
	if maxMemory < 4 || maxCPUs < 2 {
		return 0, 0, errors.New("主机可分配的资源不足，Windows 11 客机至少需要 4 GiB 内存与 2 个逻辑核心")
	}
	if allocation.Automatic {
		return min(maxMemory, 6), min(maxCPUs, 4), nil
	}
	if allocation.MemoryGiB < 4 || allocation.MemoryGiB > maxMemory || allocation.CPUs < 2 || allocation.CPUs > maxCPUs {
		return 0, 0, fmt.Errorf("Windows 客机内存须为 4–%d GiB，CPU 须为 2–%d 个逻辑核心", maxMemory, maxCPUs)
	}
	return allocation.MemoryGiB, allocation.CPUs, nil
}

func (s *Server) configuredWindowsResources() (int, int, error) {
	allocation, err := cage.LoadResourceAllocation(s.configPath)
	if err != nil {
		return 0, 0, err
	}
	maxMemory, maxCPUs := windowsResourceMaximum()
	return windowsAllocationLimits(allocation, maxMemory, maxCPUs)
}

func (s *Server) windowsResourceSample(ctx context.Context, state guestos.WindowsVMResourceState) (map[string]any, error) {
	cmd, err := guestos.WindowsRelayCommand(ctx, s.profile, "resources")
	if err != nil {
		return nil, err
	}
	cmd.Stdin = bytes.NewBufferString("{\"action\":\"sample\"}\n")
	cmd.Stderr = io.Discard
	var buffer bytes.Buffer
	cmd.Stdout = &boundedWindowsStatusWriter{buffer: &buffer, remaining: 16384}
	if err = cmd.Run(); err != nil {
		return nil, errors.New("Windows 客机资源观测尚未连接")
	}
	var observation map[string]any
	decoder := json.NewDecoder(&buffer)
	if decoder.Decode(&observation) != nil || decoder.Decode(new(any)) != io.EOF || observation["version"] != float64(1) || observation["cpus"] != float64(state.CPUs) {
		return nil, errors.New("Windows 客机资源观测与配置不一致")
	}
	guest, ok := observation["guest"].(map[string]any)
	if !ok {
		return nil, errors.New("Windows 客机内存观测无效")
	}
	total, tok := guest["total_bytes"].(float64)
	available, aok := guest["available_bytes"].(float64)
	if !tok || !aok || total < 2<<30 || available < 0 || available > total {
		return nil, errors.New("Windows 客机内存观测无效")
	}
	maximum := state.MemoryStartup
	if state.DynamicMemory {
		maximum = state.MemoryMaximum
	}
	observation["maximum_bytes"] = maximum
	observation["target_bytes"] = state.MemoryAssigned
	return observation, nil
}
func (s *Server) windowsResourceStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	state, err := guestos.WindowsVMResources(ctx, s.profile, 0, 0)
	if err != nil {
		requestError(w, err)
		return
	}
	allocation, err := cage.LoadResourceAllocation(s.configPath)
	if err != nil {
		requestError(w, err)
		return
	}
	maxMemory, maxCPUs := windowsResourceMaximum()
	recommendedMemory, recommendedCPUs, _ := windowsAllocationLimits(cage.ResourceAllocation{Automatic: true}, maxMemory, maxCPUs)
	memory, cpus, err := windowsAllocationLimits(allocation, maxMemory, maxCPUs)
	allocationError := ""
	if err != nil {
		allocationError = err.Error()
		if !allocation.Automatic {
			memory, cpus = allocation.MemoryGiB, allocation.CPUs
		}
	}
	maximum := state.MemoryStartup
	if state.DynamicMemory {
		maximum = state.MemoryMaximum
	}
	running := state.State == "Running"
	var observation map[string]any
	if running {
		if sample, sampleErr := s.windowsResourceSample(ctx, state); sampleErr == nil {
			observation = sample
		} else if allocationError == "" {
			allocationError = sampleErr.Error()
		}
	}
	jsonReply(w, map[string]any{"runtime": "hyperv", "running": running, "automatic": allocation.Automatic, "memory_gib": memory, "cpus": cpus, "min_memory_gib": 4, "min_cpus": 2, "max_memory_gib": maxMemory, "max_cpus": maxCPUs, "resource_error": allocationError, "recommended_memory_bytes": uint64(recommendedMemory) << 30, "recommended_cpus": recommendedCPUs, "configured_memory_bytes": maximum, "configured_cpus": state.CPUs, "dynamic_memory": state.DynamicMemory, "resources_verified": allocationError == "", "restart_required": running && err == nil && !state.MatchesAllocation(memory, cpus), "observation": observation, "unavailable_observations": []string{"oom_kills", "memory_stall_percent", "swap_total_bytes", "swap_free_bytes", "host_pressure"}})
}
func (s *Server) applyWindowsResources(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Allocation *cage.ResourceAllocation `json:"allocation"`
	}
	if !readJSON(w, r, &input) {
		return
	}
	if input.Allocation != nil {
		maxMemory, maxCPUs := windowsResourceMaximum()
		_, _, err := windowsAllocationLimits(*input.Allocation, maxMemory, maxCPUs)
		if err != nil {
			requestError(w, err)
			return
		}
		if !s.op.TryLock() {
			http.Error(w, "另一个操作仍在进行，请稍候", http.StatusConflict)
			return
		}
		err = writePrivate(cage.ResourceAllocationPath(s.configPath), input.Allocation)
		s.op.Unlock()
		if err != nil {
			requestError(w, err)
			return
		}
		jsonReply(w, map[string]bool{"ok": true})
		return
	}
	s.launch(w, "resources", func() error {
		memory, cpus, err := s.configuredWindowsResources()
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(s.operationContext(), 5*time.Minute)
		defer cancel()
		state, err := guestos.WindowsVMResources(ctx, s.profile, 0, 0)
		if err != nil {
			return err
		}
		if state.State == "Running" {
			if err = s.stopWindowsGuest(); err != nil {
				return err
			}
		}
		if _, err = guestos.WindowsVMResources(ctx, s.profile, memory, cpus); err != nil {
			return err
		}
		return s.startWindowsGuest()
	})
}
