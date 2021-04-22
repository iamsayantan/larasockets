package app_managers

import (
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/config"
)

type configApplicationManager struct {
	apps map[string]*larasockets.Application
}

// NewConfigManager returns an ApplicationManager instance that is managed in memory.
// The apps are loaded from configuration file.
func NewConfigManager(appsConfig []config.AppConfig) larasockets.ApplicationManager {
	apps := make(map[string]*larasockets.Application, 0)
	for _, configApp := range appsConfig {
		app := larasockets.NewApplication(configApp)
		apps[app.Id()] = app
	}

	return &configApplicationManager{apps: apps}
}

func (c *configApplicationManager) All() []*larasockets.Application {
	apps := make([]*larasockets.Application, 0)
	for _, app := range c.apps {
		apps = append(apps, app)
	}

	return apps
}

func (c *configApplicationManager) FindById(id string) *larasockets.Application {
	app, ok := c.apps[id]
	if !ok {
		return nil
	}

	return app
}

func (c *configApplicationManager) FindByKey(key string) *larasockets.Application {
	for _, app := range c.apps {
		if app.Key() == key {
			return app
		}
	}

	return nil
}
