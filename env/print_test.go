package env

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, "", actual, "Origin value: %s", nil)
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
