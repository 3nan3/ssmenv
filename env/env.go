package env

import (
	"os"
)

type Env struct {
	envs map[string]*string
}

func New() *Env {
	envs := map[string]*string{}
	return &Env{envs}
}

func (env *Env) GetEnvs() map[string]*string {
	return env.envs
}

func (env *Env) GetEnv(name string) *string {
	return env.envs[name]
}

func (env *Env) PutEnv(name string, value *string) {
	if value == nil {
		env.envs[name] = nil
		return
	}
	clone := *value
	env.envs[name] = &clone
}

func (env *Env) ApplyEnv() {
	for name, value := range env.envs {
		os.Setenv(name, *value)
	}
}
