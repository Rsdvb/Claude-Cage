package workbench

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/go-units"
	"github.com/sylvinus/agent-vm/internal/cage"
	"github.com/sylvinus/agent-vm/internal/containerenv"
	"github.com/sylvinus/agent-vm/internal/guestos"
	"go.yaml.in/yaml/v3"
)

func (s *Server) instanceDir() string {
	if s.windowsGuest != nil {
		// The profile owns storage paths, while the server retains the bound
		// VM/user identity across a migration. Relay calls load this file too.
		if current, err := guestos.LoadWindowsProfile(s.profile); err == nil && current != nil && current.VMID == s.windowsGuest.VMID && current.GUIUserSID == s.windowsGuest.GUIUserSID {
			return current.InstanceDirectory
		}
		return s.windowsGuest.InstanceDirectory
	}
	if s.macosGuest != nil {
		return filepath.Join(s.profile, "state", "lima", s.macosGuest.Instance)
	}
	if containerenv.Configured(s.profile) {
		return filepath.Join(s.profile, "state", "container", "environment")
	}
	return filepath.Join(s.profile, "state", "lima", cage.InstanceName(s.profile, s.environmentDir()))
}

func (s *Server) resourceStatus(w http.ResponseWriter, r *http.Request) {
	if s.windowsGuest != nil {
		s.windowsResourceStatus(w, r)
		return
	}
	if s.macosGuest != nil {
		s.macOSResourceStatus(w, r)
		return
	}
	if containerenv.Configured(s.profile) {
		environment, err := containerenv.Load(s.profile)
		if err != nil {
			requestError(w, err)
			return
		}
		state, err := environment.State(r.Context())
		if err != nil {
			requestError(w, err)
			return
		}
		config := environment.Config()
		jsonReply(w, map[string]any{"running": state.Running, "automatic": false, "runtime": "container",
			"recommended_memory_bytes": config.MemoryBytes, "recommended_cpus": config.CPUs,
			"configured_memory_bytes": state.MemoryBytes, "configured_cpus": state.CPUs,
			"resources_verified": state.ResourcesVerified, "resource_error": state.ResourceError, "pids_limit": state.PIDs,
			"restart_required": false, "observation": nil})
		return
	}
	a, err := cage.LoadResourceAllocation(s.configPath)
	if err != nil {
		requestError(w, err)
		return
	}
	memory, cpus, err := a.Limits()
	allocationError := ""
	if err != nil {
		allocationError = err.Error()
		memory, cpus = a.MemoryGiB, a.CPUs
	}
	maxMemory, maxCPUs := cage.ResourceLimits()
	recommendedMemory, recommendedCPUs, _ := (cage.ResourceAllocation{Automatic: true}).Limits()
	running := s.vmRunning()
	var configured struct {
		Memory string `yaml:"memory"`
		CPUs   int    `yaml:"cpus"`
	}
	if b, err := os.ReadFile(filepath.Join(s.instanceDir(), "lima.yaml")); err == nil {
		_ = yaml.Unmarshal(b, &configured)
	}
	maximum, _ := units.RAMInBytes(configured.Memory)
	var observation *cage.ResourceStatus
	status, err := cage.ReadResourceStatus(s.instanceDir())
	if running && err == nil {
		observation = &status
	}
	jsonReply(w, map[string]any{
		"running": running, "automatic": a.Automatic,
		"memory_gib": memory, "cpus": cpus,
		"max_memory_gib": maxMemory, "max_cpus": maxCPUs,
		"resource_error":           allocationError,
		"recommended_memory_bytes": uint64(recommendedMemory) << 30, "recommended_cpus": recommendedCPUs,
		"configured_memory_bytes": maximum, "configured_cpus": configured.CPUs,
		"restart_required": running && (maximum != int64(memory)<<30 || configured.CPUs != cpus || observation == nil),
		"observation":      observation,
	})
}

func (s *Server) applyResources(w http.ResponseWriter, r *http.Request) {
	if s.windowsGuest != nil {
		s.applyWindowsResources(w, r)
		return
	}
	var input struct {
		Allocation *cage.ResourceAllocation `json:"allocation"`
	}
	if !readJSON(w, r, &input) {
		return
	}
	if containerenv.Configured(s.profile) {
		http.Error(w, "容器资源由服务器部署配置管理", http.StatusBadRequest)
		return
	}
	if input.Allocation != nil {
		if _, _, err := s.resourceAllocationLimits(*input.Allocation); err != nil {
			requestError(w, err)
			return
		}
		if !s.op.TryLock() {
			http.Error(w, "另一个操作仍在进行，请稍候", http.StatusConflict)
			return
		}
		err := writePrivate(cage.ResourceAllocationPath(s.configPath), input.Allocation)
		s.op.Unlock()
		if err != nil {
			requestError(w, err)
			return
		}
		jsonReply(w, map[string]bool{"ok": true})
		return
	}
	s.launch(w, "resources", func() error {
		allocation, err := cage.LoadResourceAllocation(s.configPath)
		if err != nil {
			return err
		}
		if _, _, err := s.resourceAllocationLimits(allocation); err != nil {
			return err
		}
		if s.vmRunning() {
			if err := s.stopEnvironment(); err != nil {
				return err
			}
		}
		if err := s.operationContext().Err(); err != nil {
			return err
		}
		return s.startEnvironment(false)
	})
}

type networkCheckError struct{ error }

func (e networkCheckError) Unwrap() error { return e.error }

func (s *Server) checkNetworkIfWanted() error {
	s.mu.Lock()
	paused := s.networkPaused
	s.mu.Unlock()
	if paused || s.macSetupPending() {
		return nil
	}
	var err error
	if s.macosGuest != nil {
		err = s.checkMacOSNetwork()
	} else {
		err = s.runSafetyCheck()
	}
	if err != nil {
		return networkCheckError{err}
	}
	return nil
}
