package env

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"bytes"
)

func TestPrintAll(t *testing.T) {
	vars := map[string]*string {
		"VAR_A": toPtr(`value`),
		"VAR_B": toPtr(`"value"`),
		"VAR_C": toPtr(`va lue`),
		"VAR_D": toPtr("va\nlue"),
		"VAR_E": toPtr(`{
  "a": "v",
  "b": "ABC\nDEF\n"
}`),
		"VAR_F": nil,
	}
	envs := New()
	for n, v := range vars {
		envs.PutEnv(n, v)
	}

	io := bytes.NewBufferString("")
	envs.printAll(io)

	expected :=
`VAR_A=value
VAR_B=\"value\"
VAR_C="va lue"
VAR_D="va
lue"
VAR_E="{
  \"a\": \"v\",
  \"b\": \"ABC\nDEF\n\"
}"
VAR_F=
`
	assert.Equal(t, expected, io.String())
}

func TestPrintDiff(t *testing.T) {
	vars := map[string]([]*string) {
		"VAR_A": []*string{ toPtr(`value`), toPtr(`value`)},
		"VAR_B": []*string{ toPtr(`value`), toPtr(`new_value`)},
		"VAR_C": []*string{ toPtr(`value`), nil},
		"VAR_D": []*string{ nil,            toPtr(`value`) },
		"VAR_E": []*string{ nil,            nil },
	}
	oldenvs, newenvs := New(), New()
	for n, vs := range vars {
		old, new := vs[0], vs[1]
		oldenvs.PutEnv(n, old)
		newenvs.PutEnv(n, new)
	}

	io := bytes.NewBufferString("")
	printDiff(io, oldenvs, newenvs, "all")
	expected :=
`- key: VAR_B
  old_value: value
  new_value: new_value
- key: VAR_C
  old_value: value
  new_value: <undefined>
- key: VAR_D
  old_value: <undefined>
  new_value: value
`
	assert.Equal(t, expected, io.String())

	io = bytes.NewBufferString("")
	printDiff(io, oldenvs, newenvs, "key")
	expected =
`- key: VAR_B
- key: VAR_C
- key: VAR_D
`
	assert.Equal(t, expected, io.String())
}

func TestToEnvValue(t *testing.T) {
	convertions := map[string]string {
		"value":   `value`,
		"'value'": `'value'`,
		`"value"`: `\"value\"`,
		`va"lue`:  `va\"lue`,
		"va lue":  `"va lue"`,
		"va\nlue": "\"va\nlue\"",
		`va\nlue`: `va\nlue`,
		"va=lue":  `va=lue`,
		"": "",
`{
  "a": "v",
  "b": "ABC\nDEF\n"
}`: 
`"{
  \"a\": \"v\",
  \"b\": \"ABC\nDEF\n\"
}"`,
	}

	for value, expected := range convertions {
		actual := toEnvValue(&value)
		assert.Equal(t, expected, actual, "Origin value: %s", value)
	}

	// nil -> ""
	actual := toEnvValue(nil)
	assert.Equal(t, "", actual, "Origin value: %s", "nil")
}

func TestToDiffValue(t *testing.T) {
	value := `{
  "a": "v",
  "b": "ABC\nDEF\n"
}`
	assert.Equal(t, value, toDiffValue(&value))
	assert.Equal(t, "<undefined>", toDiffValue(nil))
}

func TestSortedName(t *testing.T) {
	envs := New()
	names := []string{"VAR_B", "VAR_A", "VAR_C"}
	for _, n := range names {
		envs.PutEnv(n, nil)
	}

	expecteds := []string{"VAR_A", "VAR_B", "VAR_C"}
	actuals := envs.sortedName()
	for i := 0; i < 3; i++ {
		assert.Equal(t, expecteds[i], actuals[i])
	}
}
