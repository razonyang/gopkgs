package core

import (
	"path"
	"plugin"
)

// PluginManager is a manager that manage plugins.
type PluginManager struct {
	// The location of plugins.
	Path string
}

// NewPluginManager returns a plugins manager with the given path.
func NewPluginManager(path string) *PluginManager {
	return &PluginManager{
		Path: path,
	}
}

// Open opens a Go plugin.
func (pm *PluginManager) Open(name string) (*plugin.Plugin, error) {
	return plugin.Open(path.Join(pm.Path, name))
}

// Lookup searchs for a symbol named symName in plugin named pluginName.
func (pm *PluginManager) Lookup(pluginName, symName string) (plugin.Symbol, error) {
	p, err := pm.Open(pluginName)
	if err != nil {
		return nil, err
	}
	return p.Lookup(symName)
}
