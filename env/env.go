package env

import (
	"os"
	"fmt"
	"io"
	"strings"
	"bytes"
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

func (env *Env) Stdout() error {
	return env.print(os.Stdout)
}

func (env *Env) print(io io.Writer) error {
	for name, value := range env.envs {
		escapedValue, err := escapeValue(value); if err != nil {
			return err
		}

		_, err = fmt.Fprintf(io, "%s=%s\n", name, escapedValue); if err != nil {
			return err
		}
	}
	return nil
}

func escapeValue(value string) (string, error) {
	if !valueNeedsQuotes(value) {
		return value, nil
	}

	escaped := bytes.NewBufferString("\"")
	for len(value) > 0 {
		i := strings.IndexAny(value, "\"\r\n")
		if i < 0 {
			i = len(value)
		}

		if _, err := fmt.Fprint(escaped, value[:i]); err != nil {
			return "", err
		}
		value = value[i:]

		if len(value) > 0 {
			var err error
			switch value[0] {
			case '"':
				_, err = fmt.Fprint(escaped, `\"`)
			case '\n', '\r':
				_, err = fmt.Fprint(escaped, "\n")
			}
			value = value[1:]
			if err != nil {
				return "", err
			}
		}
	}
	return escaped.String() + "\"", nil
}

func valueNeedsQuotes(value string) bool {
	return strings.ContainsAny(value, " =\\\"\r\n")
}
