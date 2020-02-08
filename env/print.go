package env

import (
	"os"
	"fmt"
	"io"
	"strings"
	"bytes"
	"sort"
)

func (env *Env) PrintAll() error {
	return env.printAll(os.Stdout)
}

func PrintDiff(oldenv *Env, newenv *Env, diff string) {
	printDiff(os.Stdout, oldenv, newenv, diff)
}

func (env *Env) printAll(io io.Writer) error {
	for _, name := range env.sortedName() {
		_, err := fmt.Fprintf(io, "%s=%s\n", name, toEnvValue(env.GetEnv(name)))
		if err != nil {
			return err
		}
	}
	return nil
}

func printDiff(io io.Writer, oldenv *Env, newenv *Env, diff string) {
	for _, name := range newenv.sortedName() {
		newv := toDiffValue(newenv.GetEnv(name))
		oldv := toDiffValue(oldenv.GetEnv(name))
		if newv != oldv {
			fmt.Fprintf(io, "- key: %s\n", name)
			if diff == "all" {
				fmt.Fprintf(io, "  old_value: %s\n  new_value: %s\n", oldv, newv)
			}
		}
	}
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

func toEnvValue(value *string) string {
	if value == nil {
		return ""
	}
	if !valueNeedsEscape(*value) {
		return *value
	}

	var escaped bytes.Buffer
	needsQuotes := valueNeedsQuotes(*value)
	for len(*value) > 0 {
		i := strings.IndexAny(*value, "\"\r\n$")
		if i < 0 {
			i = len(*value)
		}

		escaped.WriteString((*value)[:i])
		value = cloneString((*value)[i:])

		if len(*value) > 0 {
			switch (*value)[0] {
			case '"':
				escaped.WriteString(`\"`)
			case '$':
				escaped.WriteString(`\$`)
			case '\n', '\r':
				escaped.WriteString("\n")
			}
			value = cloneString((*value)[1:])
		}
	}
	if needsQuotes {
		return `"` + escaped.String() + `"`
	} else {
		return escaped.String()
	}
}

func valueNeedsEscape(value string) bool {
	return strings.ContainsAny(value, " \r\n\"$")
}

func valueNeedsQuotes(value string) bool {
	return strings.ContainsAny(value, " \r\n")
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
