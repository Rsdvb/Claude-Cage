package workbench

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/sylvinus/agent-vm/internal/cage"
	"github.com/sylvinus/agent-vm/internal/host"
)

func TestResourceAllocationSaveAndStartup(t *testing.T) {
	oldMem, oldCPU := host.MemGiB, host.CPUs
	t.Cleanup(func() { host.MemGiB, host.CPUs = oldMem, oldCPU })
	host.MemGiB = func() int { return 24 }
	host.CPUs = func() int { return 10 }
	t.Setenv("AGENT_VM_HOST_SHARE", "2")
	s := testServer(t)
	check := func(memory, cpus int) {
		t.Helper()
		m, c, err := cage.ConfiguredResourceLimits(s.configPath)
		if err != nil || m != memory || c != cpus {
			t.Fatalf("startup limits = %d/%d, %v", m, c, err)
		}
	}
	check(4, 4)
	for _, body := range []string{
		`{"allocation":{"automatic":false,"memory_gib":0,"cpus":2}}`,
		`{"allocation":{"automatic":false,"memory_gib":13,"cpus":2}}`,
		`{"allocation":{"automatic":false,"memory_gib":4,"cpus":6}}`,
		`{"allocation":{"automatic":false,"memory_gib":4.5,"cpus":2}}`,
	} {
		if w := request(s, "POST", "/api/environment/resources", s.origin, body, true); w.Code == 200 {
			t.Fatalf("accepted invalid allocation: %s", body)
		}
		check(4, 4)
	}
	before, _ := os.ReadFile(s.configPath)
	if w := request(s, "POST", "/api/environment/resources", s.origin, `{"allocation":{"automatic":false,"memory_gib":4,"cpus":2}}`, true); w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	check(4, 2)
	after, _ := os.ReadFile(s.configPath)
	if string(before) != string(after) {
		t.Fatal("resource save changed network configuration")
	}
	w := request(s, "GET", "/api/environment/resources", "", "", true)
	var status struct {
		Automatic bool `json:"automatic"`
		Memory    int  `json:"memory_gib"`
		CPUs      int  `json:"cpus"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &status); err != nil || status.Automatic || status.Memory != 4 || status.CPUs != 2 {
		t.Fatal(w.Body.String(), err)
	}
	s.op.Lock()
	w = request(s, "POST", "/api/environment/resources", s.origin, `{"allocation":{"automatic":true}}`, true)
	s.op.Unlock()
	if w.Code != 409 {
		t.Fatal("save must not race an environment operation", w.Code)
	}
	check(4, 2)
	if w = request(s, "POST", "/api/environment/resources", s.origin, `{"allocation":{"automatic":true}}`, true); w.Code != 200 {
		t.Fatal(w.Body.String())
	}
	check(4, 4)
}
