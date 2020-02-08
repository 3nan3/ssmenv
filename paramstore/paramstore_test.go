package paramstore

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/3nan3/ssmenv/env"
)

func TestPutParamsInEnvs(t *testing.T) {
	emptyPattern := "empty"
	ps := New("/path", emptyPattern)
	params := []*ssm.Parameter{
		&ssm.Parameter{Name: aws.String("VAR_A"), Value: aws.String("value")},
		&ssm.Parameter{Name: aws.String("VAR_B"), Value: aws.String("va\nlue")},
		&ssm.Parameter{Name: aws.String("VAR_C"), Value: aws.String(emptyPattern)},
	}

	envs := env.New()
	ps.putParamsInEnvs(envs, params)

	expecteds := map[string]*string{
		"VAR_A": aws.String("value"),
		"VAR_B": aws.String("va\nlue"),
		"VAR_C": aws.String(""),
	}
	for n, expected := range expecteds {
		assert.Equal(t, *expected, *envs.GetEnv(n))
	}
}
