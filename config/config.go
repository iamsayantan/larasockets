package config

import "errors"

type LarasocketsConfig struct {
	Apps     []AppConfig
	Server   ServerConfig
	Database DatabaseConfig
}

func (c LarasocketsConfig) Validate() error {
	if len(c.Apps) == 0 {
		return errors.New("at least one application needs to be configured")
	}

	for _, app := range c.Apps {
		if err := app.validate(); err != nil {
			return err
		}
	}

	if err := c.Server.validate(); err != nil {
		return err
	}

	if err := c.Database.validate(); err != nil {
		return err
	}

	return nil
}

type AppConfig struct {
	ID                   string
	Name                 string
	Key                  string
	Secret               string
	Capacity             int
	EnableClientMessages bool
	EnableStatistics     bool
	AllowedOrigins       []string
}

func (a *AppConfig) validate() error {
	if a.ID == "" {
		return errors.New("app id can not be empty")
	}

	if a.Key == "" {
		return errors.New("application key can not be empty")
	}

	if a.Secret == "" {
		return errors.New("application secret can not be empty")
	}

	return nil
}

type ServerConfig struct {
	Port        string
	TLS         bool
	Key         string
	Certificate string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Db       string
}

func (s ServerConfig) validate() error {
	if !s.TLS {
		return nil
	}

	if s.Key == "" {
		return errors.New("key path can not be empty")
	}

	if s.Certificate == "" {
		return errors.New("certificate can not be empty")
	}

	return nil
}

func (d DatabaseConfig) validate() error {
	if d.Host == "" {
		return errors.New("database host is required")
	}

	if d.Username == "" {
		return errors.New("database username is required")
	}

	if d.Db == "" {
		return errors.New("database name is required")
	}

	return nil
}
