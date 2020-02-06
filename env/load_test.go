package env

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"bytes"
)

func TestLoad(t *testing.T) {
	os.Setenv("EXPORTED_VAR", "exported_value")
	vars := 
`VAR_A=value
VAR_B='value'
VAR_C=''value''
VAR_D="value"
VAR_E=""value""
VAR_F=va"lue
VAR_G=va lue
VAR_H=va\nlue
VAR_I=va=lue
VAR_J=
VAR_K={\n  "a": "v",\n  "b": "ABC\nDEF\n"\n}
VAR_L="{\n  \"a\": \"v\",\n  \"b\": \"ABC\nDEF\n\"\n}"
VAR_M=$VAR_A
VAR_N='$VAR_A'
VAR_O=$EXPORTED_VAR
VAR_P=value
VAR_P=override
#VAR_Q=value
`	// NOTE not support multiline variable: https://github.com/joho/godotenv/issues/64

	envs := New()
	reader := bytes.NewBufferString(vars)
	err := envs.load(reader)

	expecteds := map[string]*string {
		"VAR_A": toPtr(`value`),
		"VAR_B": toPtr(`value`),
		"VAR_C": toPtr(`'value'`),
		"VAR_D": toPtr(`value`),
		"VAR_E": toPtr(`"value"`),
		"VAR_F": toPtr(`va"lue`),
		"VAR_G": toPtr(`va lue`),
		"VAR_H": toPtr(`va\nlue`),
		"VAR_I": toPtr(`va=lue`),
		"VAR_J": toPtr(``),
		"VAR_K": toPtr(`{\n  "a": "v",\n  "b": "ABC\nDEF\n"\n}`),
		"VAR_L": toPtr(`{\n  "a": "v",\n  "b": "ABC\nDEF\n"\n}`),
		"VAR_M": toPtr(`value`),
		"VAR_N": toPtr(`$VAR_A`),
		"VAR_O": toPtr(``),
		"VAR_P": toPtr(`override`),
		"VAR_Q": nil,
	}
	if assert.Nil(t, err) {
		for name, actual := range envs.GetEnvs() {
			if actual == nil {
				assert.Equal(t, expecteds[name], actual)
			} else {
				assert.Equal(t, *expecteds[name], *actual)
			}
		}
	}
}
