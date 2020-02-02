package env

import (
	"os"
	"fmt"
	"io"
	"strings"
	"bytes"
	"sort"
)

func (env *Env) Stdout() error {
	return env.print(os.Stdout)
}

func PrintDiff(oldenv *Env, newenv *Env, diff string) {
	for _, name := range newenv.sortedName() {
		newv := newenv.GetEnv(name)
		oldv := oldenv.GetEnv(name)
		if newv != oldv {
			fmt.Printf("- key: %s\n", name)
			if diff == "all" {
				fmt.Printf("  old_value: %s\n  new_value: %s\n", toDiffValue(oldv), toDiffValue(newv))
			}
		}
	}
}

func (env *Env) print(io io.Writer) error {
	for _, name := range env.sortedName() {
		escaped, err := toEnvValue(env.GetEnv(name)); if err != nil {
			return err
		}

		_, err = fmt.Fprintf(io, "%s=%s\n", name, escaped); if err != nil {
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

func toEnvValue(value *string) (string, error) {
	if value == nil {
		return "", nil
	}
	if !valueNeedsQuotes(*value) {
		return *value, nil
	}

	escaped := bytes.NewBufferString("\"")
	for len(*value) > 0 {
		i := strings.IndexAny(*value, "\"\r\n")
		if i < 0 {
			i = len(*value)
		}

		if _, err := fmt.Fprint(escaped, (*value)[:i]); err != nil {
			return "", err
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
				return "", err
			}
		}
	}
	return escaped.String() + "\"", nil
}

func valueNeedsQuotes(value string) bool {
	return strings.ContainsAny(value, " =\\\"\r\n")
}

func toDiffValue(str *string) string {
	if str == nil {
		return "<undefined>"
	}
	return *str
}

func cloneString(str string) *string {
	clone := str
	return &clone
}
