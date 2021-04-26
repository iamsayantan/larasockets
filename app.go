package larasockets

import "github.com/iamsayantan/larasockets/config"

// ApplicationManager defines methods to manage all the different applications.
type ApplicationManager interface {
	// All returns all the applications running in our system.
	All() []*Application

	// FindByApplicationId returns an application instance with the given id.
	FindById(id string) *Application

	// FindByApplicationKey returns an application instance by the apps key.
	FindByKey(key string) *Application
}

// Application represents an individual application instance. The server can hold multiple
// apps, and each app is an logically isolated and serves one client.
type Application struct {
	id                   string
	appKey               string
	appSecret            string
	appName              string
	host                 string
	path                 string
	capacity             int
	clientMessageEnabled bool
}

func (app *Application) Id() string {
	return app.id
}

func (app *Application) Key() string {
	return app.appKey
}

func (app *Application) Secret() string {
	return app.appSecret
}

func (app *Application) Name() string {
	return app.appName
}

func (app *Application) Host() string {
	return app.host
}

func (app *Application) Path() string {
	return app.path
}

func (app *Application) Capacity() int {
	return app.capacity
}

func (app *Application) ClientMessageEnabled() bool {
	return app.clientMessageEnabled
}

func (app *Application) SetName(name string) {
	if name == "" {
		return
	}

	app.appName = name
}

func (app *Application) EnableClientMessages() {
	app.clientMessageEnabled = true
}

func (app *Application) SetCapacity(capacity int) {
	if capacity == 0 {
		return
	}

	app.capacity = capacity
}

// NewApplication returns a new application instance.
func NewApplication(appConfig config.AppConfig) *Application {
	app := &Application{
		id:        appConfig.ID,
		appKey:    appConfig.Key,
		appSecret: appConfig.Secret,
		appName:   appConfig.Name,
	}

	if appConfig.Capacity > 0 {
		app.SetCapacity(appConfig.Capacity)
	}

	return app
}
