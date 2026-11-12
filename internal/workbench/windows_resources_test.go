package workbench

import (
	"path/filepath"
	"testing"

	"github.com/sylvinus/agent-vm/internal/cage"
	"github.com/sylvinus/agent-vm/internal/guestos"
	"github.com/sylvinus/agent-vm/internal/host"
)

func TestWindowsAllocationLimits(t *testing.T) {
	for _, tc := range []struct{ hostGiB, hostCPUs, wantGiB, wantCPUs int }{
		{1024, 256, 6, 4}, {32, 8, 6, 4}, {6, 2, 4, 2},
	} {
		maxGiB, maxCPUs := guestos.WindowsVMResourceCeiling(tc.hostGiB, tc.hostCPUs)
		mem, cpus, err := windowsAllocationLimits(cage.ResourceAllocation{Automatic: true}, maxGiB, maxCPUs)
		if err != nil || mem != tc.wantGiB || cpus != tc.wantCPUs {
			t.Fatalf("automatic allocation = %d GiB, %d CPUs, %v", mem, cpus, err)
		}
		for _, pair := range [][2]int{{2, 2}, {4, 1}, {maxGiB + 1, 2}, {4, maxCPUs + 1}, {0, 0}} {
			if _, _, err := windowsAllocationLimits(cage.ResourceAllocation{MemoryGiB: pair[0], CPUs: pair[1]}, maxGiB, maxCPUs); err == nil {
				t.Fatalf("accepted out-of-range manual allocation %v", pair)
			}
		}
		if _, _, err := windowsAllocationLimits(cage.ResourceAllocation{MemoryGiB: 4, CPUs: 2}, maxGiB, maxCPUs); err != nil {
			t.Fatal(err)
		}
	}
	if _, _, err := windowsAllocationLimits(cage.ResourceAllocation{Automatic: true}, 3, 2); err == nil {
		t.Fatal("accepted automatic allocation below the Windows 11 minimum")
	}
}

func TestWindowsColdStartUsesManagerAndHostShareLimits(t *testing.T) {
	originalMemory, originalCPUs := host.MemGiB, host.CPUs
	t.Cleanup(func() { host.MemGiB, host.CPUs = originalMemory, originalCPUs })
	host.MemGiB = func() int { return 128 }
	host.CPUs = func() int { return 256 }
	s := &Server{configPath: filepath.Join(t.TempDir(), "config.json")}
	for _, tc := range []struct {
		share      string
		wantMemory int
	}{{"2", 64}, {"1", 96}} {
		t.Setenv("AGENT_VM_HOST_SHARE", tc.share)
		memory, cpus, err := s.configuredWindowsResources()
		if err != nil || memory != min(tc.wantMemory, 6) || cpus != 4 {
			t.Fatalf("share %s: got %d GiB / %d CPUs, %v", tc.share, memory, cpus, err)
		}
	}
	if err := writePrivate(cage.ResourceAllocationPath(s.configPath), cage.ResourceAllocation{MemoryGiB: 32, CPUs: 128}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.configuredWindowsResources(); err == nil {
		t.Fatal("legacy allocation above manager cap was accepted")
	}
}

func TestWindowsDynamicMemoryRequiresExpectedStartupAndMinimum(t *testing.T) {
	state := guestos.WindowsVMResourceState{
		DynamicMemory: true, MemoryMinimum: 4 << 30, MemoryStartup: 4 << 30,
		MemoryMaximum: 16 << 30, CPUs: 8,
	}
	if !state.MatchesAllocation(16, 8) {
		t.Fatal("current configuration requires an unnecessary restart")
	}
	state.MemoryMinimum = 2 << 30
	if state.MatchesAllocation(16, 8) {
		t.Fatal("ignored the previous minimum memory")
	}
	state.MemoryMinimum, state.MemoryStartup = 4<<30, 8<<30
	if state.MatchesAllocation(16, 8) {
		t.Fatal("ignored the previous startup memory")
	}
}
