package env

import (
	"fmt"
)

type Env struct {
	envs map[string]string
}

func New() *Env {
	envs := map[string]string{}
	return &Env{envs}
}

func (env *Env) GetEnv(name string) string {
	return env.envs[name]
}

func (env *Env) PutEnv(name string, value string) {
	env.envs[name] = value
}

func (env *Env) Print() {
	for name, value := range env.envs {
		fmt.Printf("%s='%s'\n", name, value)
	}
}
