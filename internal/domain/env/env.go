package env

import (
	"os"
	"strconv"
)

type Env struct {
	AppEnv  string
	AppPort int
}

func New() (*Env, error) {
	port := 8080
	if p := os.Getenv("APP_PORT"); p != "" {
		if val, err := strconv.Atoi(p); err == nil {
			port = val
		}
	}

	appEnv := "development"
	if e := os.Getenv("APP_ENV"); e != "" {
		appEnv = e
	}

	return &Env{
		AppEnv:  appEnv,
		AppPort: port,
	}, nil
}
