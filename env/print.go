package env

import (
	"os"
	"fmt"
	"io"
	"strings"
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
				fmt.Printf("  old_value: %s\n  new_value: %s\n", forPrintValue(oldv), forPrintValue(newv))
			}
		}
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

func forPrintValue(str *string) string {
	if str == nil {
		return "<undefined>"
	}
	return *str
}
