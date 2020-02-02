package env

import (
	"os"
	"fmt"
	"strings"
	"bytes"
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

func escapeValue(value *string) (*string, error) {
	if value == nil {
		return nil, nil
	}
	if !valueNeedsQuotes(*value) {
		return value, nil
	}

	escaped := bytes.NewBufferString("\"")
	for len(*value) > 0 {
		i := strings.IndexAny(*value, "\"\r\n")
		if i < 0 {
			i = len(*value)
		}

		if _, err := fmt.Fprint(escaped, (*value)[:i]); err != nil {
			return cloneString(""), err
		}
		value = cloneString((*value)[i:])

		if len(*value) > 0 {
			var err error
			switch (*value)[0] {
			case '"':
				_, err = fmt.Fprint(escaped, `\"`)
			case '\n', '\r':
				_, err = fmt.Fprint(escaped, "\n")
			}
			value = cloneString((*value)[1:])
			if err != nil {
				return cloneString(""), err
			}
		}
	}
	return cloneString(escaped.String() + "\""), nil
}

func valueNeedsQuotes(value string) bool {
	return strings.ContainsAny(value, " =\\\"\r\n")
}

func cloneString(str string) *string {
	clone := str
	return &clone
}
