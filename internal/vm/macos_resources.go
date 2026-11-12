package vm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/go-units"
	"github.com/sylvinus/agent-vm/internal/cage"
	"go.yaml.in/yaml/v3"
)

// MacOSResourceRequirements belong to the installed guest, independent of the
// user's current allocation. Older guests retain their known bootable floor.
type MacOSResourceRequirements struct {
	Version   int    `json:"version"`
	MemoryGiB int    `json:"minimum_memory_gib"`
	CPUs      int    `json:"minimum_cpus"`
	Source    string `json:"source"`
}

const macOSRequirementsFile = "cage-macos-resource-requirements.json"

func (requirements MacOSResourceRequirements) Limits(allocation cage.ResourceAllocation) (int, int, error) {
	memory, cpus, err := allocation.Limits()
	if allocation.Automatic {
		maxMemory, maxCPUs := cage.ResourceLimits()
		memory = min(maxMemory, max(6, requirements.MemoryGiB))
		cpus = min(maxCPUs, max(4, requirements.CPUs))
	}
	if err != nil {
		return 0, 0, err
	}
	if memory < requirements.MemoryGiB || cpus < requirements.CPUs {
		return 0, 0, fmt.Errorf("这台 macOS 客机至少需要 %d GiB 内存与 %d 个 CPU 逻辑核心，请调整资源配置", requirements.MemoryGiB, requirements.CPUs)
	}
	return memory, cpus, nil
}

func ReadMacOSResourceRequirements(instanceDir string) (MacOSResourceRequirements, error) {
	var requirements MacOSResourceRequirements
	data, err := os.ReadFile(filepath.Join(instanceDir, macOSRequirementsFile))
	if err == nil {
		if json.Unmarshal(data, &requirements) != nil || requirements.Version != 1 || requirements.MemoryGiB < 1 || requirements.CPUs < 1 ||
			(requirements.Source != "restore-image" && requirements.Source != "existing-configuration") {
			return requirements, errors.New("macOS 最低资源要求记录无效，请检查原环境配置")
		}
		if requirements.Source == "restore-image" {
			return requirements, nil
		}
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return requirements, errors.New("无法读取 macOS 最低资源要求记录")
	}
	data, err = os.ReadFile(filepath.Join(instanceDir, "lima.yaml"))
	if err != nil {
		return requirements, errors.New("无法读取这台 macOS 客机的启动配置")
	}
	var config struct {
		OS     string `yaml:"os"`
		Memory string `yaml:"memory"`
		CPUs   int    `yaml:"cpus"`
		Images []struct {
			Location string `yaml:"location"`
		} `yaml:"images"`
	}
	if yaml.Unmarshal(data, &config) != nil || config.OS != "Darwin" {
		return requirements, errors.New("macOS 客机启动配置无效")
	}
	candidates := []string{filepath.Join(instanceDir, "image.ipsw")}
	if len(config.Images) == 1 && filepath.IsAbs(config.Images[0].Location) {
		candidates = append(candidates, config.Images[0].Location)
	}
	for _, image := range candidates {
		if info, err := os.Stat(image); err == nil && info.Mode().IsRegular() {
			if requirements, err := macOSRestoreResourceRequirements(image); err == nil {
				return requirements, nil
			}
		}
	}
	if requirements.Version == 1 {
		return requirements, nil
	}
	memory, err := units.RAMInBytes(config.Memory)
	if err != nil || memory < 1<<30 || config.CPUs < 1 {
		return requirements, errors.New("无法确定已有 macOS 客机的可启动资源配置")
	}
	return MacOSResourceRequirements{Version: 1, MemoryGiB: int((memory-1)/(1<<30) + 1), CPUs: config.CPUs, Source: "existing-configuration"}, nil
}

func saveMacOSResourceRequirements(instanceDir string, requirements MacOSResourceRequirements) error {
	data, err := json.Marshal(requirements)
	if err != nil {
		return err
	}
	if current, err := os.ReadFile(filepath.Join(instanceDir, macOSRequirementsFile)); err == nil && bytes.Equal(current, data) {
		return nil
	}
	file, err := os.CreateTemp(instanceDir, ".macos-resources-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	_, err = file.Write(data)
	if err = errors.Join(err, file.Close()); err != nil {
		return err
	}
	return os.Rename(file.Name(), filepath.Join(instanceDir, macOSRequirementsFile))
}
