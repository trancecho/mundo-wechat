package server

import (
	"fmt"
	"sync"
)

var (
	Modules   = make(map[string]ModuleInfo)
	modulesMu sync.RWMutex
)

type Module interface {
	GetModuleInfo() ModuleInfo

	Init()

	PostInit()

	Serve(server *Server)

	Start(server *Server)

	Stop(server *Server, wg *sync.WaitGroup)
}

func RegisterModule(instance Module) {
	mod := instance.GetModuleInfo()
	if mod.Instance == nil {
		panic("missing ModuleInfo.Instance")
	}

	modulesMu.Lock()
	defer modulesMu.Unlock()

	if _, ok := Modules[mod.ID.String()]; ok {
		panic(fmt.Sprintf("module already registered: %s", mod.ID))
	}
	Modules[mod.ID.String()] = mod
}

// GetModule - 获取一个已注册的 Module 的 ModuleInfo
func GetModule(id ModuleID) (ModuleInfo, error) {
	modulesMu.Lock()
	defer modulesMu.Unlock()
	m, ok := Modules[id.String()]
	if !ok {
		return ModuleInfo{}, fmt.Errorf("module not registered: %s", id)
	}
	return m, nil
}
