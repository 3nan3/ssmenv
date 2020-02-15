package paramstore

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/3nan3/ssmenv/env"
)

type MockSSMAPI struct {
	mock.Mock
	ssmiface.SSMAPI
}

func (mocksvc *MockSSMAPI) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
    args := mocksvc.Called(input)
    return args.Get(0).(*ssm.GetParameterOutput), args.Error(1)
}

func TestGetEnv(t *testing.T) {
	emptyPattern := "empty"
	ps := New("/path", emptyPattern)

	mocks := map[string]*string{
		"VAR_A": aws.String("value"),
		"VAR_B": aws.String(emptyPattern),
		"VAR_C": nil,
	}
	mocksvc := new(MockSSMAPI)
	for n, v := range mocks {
		input := &ssm.GetParameterInput {
			Name: aws.String(ps.nameWithPath(n)),
			WithDecryption: aws.Bool(true),
		}
		if v != nil {
			mocksvc.On("GetParameter", input).Return(
				&ssm.GetParameterOutput{
					Parameter: &ssm.Parameter{
						Name: aws.String(n),
						Value: aws.String(*v),
					},
				},
				nil,
			)
		} else {
			mocksvc.On("GetParameter", input).Return(&ssm.GetParameterOutput{}, &ssm.ParameterNotFound{})
		}
	}
	ps.svc = mocksvc

	expecteds := map[string]*string{
		"VAR_A": aws.String("value"),
		"VAR_B": aws.String(""),
		"VAR_C": nil,
	}
	for name, expected := range expecteds {
		envs, err := ps.GetEnv(name)
		if assert.Nil(t, err) {
			assert.Equal(t, expected, envs.GetEnv(name))
		}
	}
}

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
