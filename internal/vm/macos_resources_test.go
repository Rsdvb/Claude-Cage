package vm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sylvinus/agent-vm/internal/cage"
	"github.com/sylvinus/agent-vm/internal/host"
)

func TestMacOSResourceRequirementsPreserveExistingFloor(t *testing.T) {
	directory := t.TempDir()
	config := filepath.Join(directory, "lima.yaml")
	if err := os.WriteFile(config, []byte("os: Darwin\nmemory: 8GiB\ncpus: 4\n"), 0600); err != nil {
		t.Fatal(err)
	}
	requirements, err := ReadMacOSResourceRequirements(directory)
	if err != nil || requirements.MemoryGiB != 8 || requirements.CPUs != 4 || requirements.Source != "existing-configuration" {
		t.Fatalf("existing guest floor lost: %+v %v", requirements, err)
	}
	if err = saveMacOSResourceRequirements(directory, requirements); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(config, []byte("os: Darwin\nmemory: 12GiB\ncpus: 6\n"), 0600); err != nil {
		t.Fatal(err)
	}
	after, err := ReadMacOSResourceRequirements(directory)
	if err != nil || after != requirements {
		t.Fatalf("minimum drifted with allocation: %+v %v", after, err)
	}
	if err = os.WriteFile(filepath.Join(directory, macOSRequirementsFile), []byte(`{"version":1,"minimum_memory_gib":0}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = ReadMacOSResourceRequirements(directory); err == nil {
		t.Fatal("corrupt minimum silently accepted")
	}
}

func TestMacOSResourceRequirementsUseRestoreMinimumAndHostShare(t *testing.T) {
	oldMemory, oldCPUs := host.MemGiB, host.CPUs
	t.Cleanup(func() { host.MemGiB, host.CPUs = oldMemory, oldCPUs })
	host.MemGiB = func() int { return 24 }
	host.CPUs = func() int { return 10 }
	t.Setenv("AGENT_VM_HOST_SHARE", "2")
	requirements := MacOSResourceRequirements{Version: 1, MemoryGiB: 4, CPUs: 2, Source: "restore-image"}
	for _, allocation := range []cage.ResourceAllocation{{MemoryGiB: 2, CPUs: 2}, {MemoryGiB: 4, CPUs: 1}, {MemoryGiB: 13, CPUs: 2}, {MemoryGiB: 4, CPUs: 6}} {
		if _, _, err := requirements.Limits(allocation); err == nil {
			t.Fatalf("accepted incompatible allocation %+v", allocation)
		}
	}
	if memory, cpus, err := requirements.Limits(cage.ResourceAllocation{Automatic: true}); err != nil || memory != 6 || cpus != 4 {
		t.Fatalf("automatic=%d/%d %v", memory, cpus, err)
	}
	requirements.MemoryGiB, requirements.CPUs = 8, 5
	if memory, cpus, err := requirements.Limits(cage.ResourceAllocation{Automatic: true}); err != nil || memory != 8 || cpus != 5 {
		t.Fatalf("automatic with higher restore minimum=%d/%d %v", memory, cpus, err)
	}
	requirements.MemoryGiB, requirements.CPUs = 4, 2
	if memory, cpus, err := requirements.Limits(cage.ResourceAllocation{MemoryGiB: 4, CPUs: 2}); err != nil || memory != 4 || cpus != 2 {
		t.Fatalf("manual=%d/%d %v", memory, cpus, err)
	}
}
