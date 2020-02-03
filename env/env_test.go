package env

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestPutEnvWhenValueIsPresent(t *testing.T) {
	envs := New()
	n := "envname"
	v := "value"
	envs.PutEnv(n, &v)

	// NOTE assert.NotSame() not support at 1.4.0
	if &v == envs.GetEnv(n) {
		t.Errorf("Value is not cloned")
	}
}

func TestPutEnvWhenValueIsNil(t *testing.T) {
	envs := New()
	n := "envname"
	envs.PutEnv(n, nil)

	assert.Nil(t, envs.GetEnv(n))
}

func TestApplyEnv(t *testing.T) {
	envs := New()
	vars := map[string]*string {
		"N1": toPtr("value"),
		"N2": toPtr(""),
		"N3": nil,
		"N4": toPtr("hello\nworld"),
		"N5": toPtr(`hello\nworld`),
	}
	for n, v := range vars {
		envs.PutEnv(n, v)
	}

	envs.ApplyEnv()
	for n, v := range vars {
		if v != nil {
			assert.Equal(t, *v, os.Getenv(n))
		}
	}
}

func toPtr(s string) *string {
	return &s
}
