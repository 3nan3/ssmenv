package env

import (
	"os"
	"fmt"
	"io"
	"strings"
	"bytes"
	"sort"
)

type Env struct {
	envs map[string]*string
}

func New() *Env {
	envs := map[string]*string{}
	return &Env{envs}
}

func PrintDiff(oldenv *Env, newenv *Env, diff string) {
	for _, name := range newenv.sortedName() {
		newv := newenv.GetEnv(name)
		oldv := oldenv.GetEnv(name)
		if newv != oldv {
			fmt.Printf("- key: %s\n", name)
			if diff == "all" {
				fmt.Printf("  old_value: %s\n  new_value: %s\n", forPrintValue(oldv), forPrintValue(newv))
			}
		}
	}
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

func (env *Env) Stdout() error {
	return env.print(os.Stdout)
}

func (env *Env) ApplyEnv() {
	for name, value := range env.envs {
		os.Setenv(name, *value)
	}
}

func (env *Env) print(io io.Writer) error {
	for _, name := range env.sortedName() {
		escapedValue, err := escapeValue(env.GetEnv(name)); if err != nil {
			return err
		}

		_, err = fmt.Fprintf(io, "%s=%s\n", name, forPrintValue(escapedValue)); if err != nil {
			return err
		}
	}
	return nil
}

func (env *Env) sortedName() []string {
	names := []string{}
	for name, _ := range env.envs {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return strings.Compare(names[i], names[j]) < 0
	})
	return names
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

func forPrintValue(str *string) string {
	if str == nil {
		return "<undefined>"
	}
	return *str
}

func cloneString(str string) *string {
	clone := str
	return &clone
}
