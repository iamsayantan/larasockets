package config

import "errors"

type LarasocketsConfig struct {
	Apps   []AppConfig
	Server ServerConfig
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
	Port      string
	TLS       bool
	TLSConfig TlsConfig
}

func (s ServerConfig) validate() error {
	if s.TLS {
		return s.TLSConfig.validate()
	}

	return nil
}

type TlsConfig struct {
	KeyPath         string
	CertificatePath string
}

func (t TlsConfig) validate() error {
	if t.KeyPath == "" {
		return errors.New("certificate key path can not be empty")
	}

	if t.CertificatePath == "" {
		return errors.New("certificate path can not be empty")
	}

	return nil
}
