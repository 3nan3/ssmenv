package paramstore

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"fmt"
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

func (mocksvc *MockSSMAPI) GetParametersByPath(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error) {
    args := mocksvc.Called(input)
    return args.Get(0).(*ssm.GetParametersByPathOutput), args.Error(1)
}

func (mocksvc *MockSSMAPI) PutParameter(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
    args := mocksvc.Called(input)
    if args.Get(0) == nil {
	    return nil, args.Error(1)
    }
    return args.Get(0).(*ssm.PutParameterOutput), args.Error(1)
}

func (mocksvc *MockSSMAPI) DeleteParameters(input *ssm.DeleteParametersInput) (*ssm.DeleteParametersOutput, error) {
    args := mocksvc.Called(input)
    return args.Get(0).(*ssm.DeleteParametersOutput), args.Error(1)
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

func TestGetEnvs(t *testing.T) {
	emptyPattern := "empty"
	ps := New("/path", emptyPattern)

	mocksvc := new(MockSSMAPI)
	mocks := map[*string](*ssm.GetParametersByPathOutput){
		aws.String("nil"): &ssm.GetParametersByPathOutput{
			Parameters: []*ssm.Parameter{
				&ssm.Parameter{Name: aws.String("VAR_A"), Value: aws.String("value_a")},
				&ssm.Parameter{Name: aws.String("VAR_B"), Value: aws.String("value_b")},
			},
			NextToken: aws.String("token1"),
		},
		aws.String("token1"): &ssm.GetParametersByPathOutput{
			Parameters: []*ssm.Parameter{
				&ssm.Parameter{Name: aws.String("VAR_C"), Value: aws.String("value_c")},
				&ssm.Parameter{Name: aws.String("VAR_D"), Value: aws.String("value_d")},
			},
			NextToken: aws.String("token2"),
		},
		aws.String("token2"): &ssm.GetParametersByPathOutput{
			Parameters: []*ssm.Parameter{
				&ssm.Parameter{Name: aws.String("VAR_E"), Value: aws.String(emptyPattern)},
			},
		},
	}
	for token, output := range mocks {
		input := &ssm.GetParametersByPathInput{
			Path: aws.String(ps.path),
			WithDecryption: aws.Bool(true),
		}
		if *token != "nil" {
			input.SetNextToken(*token)
		}
		mocksvc.On("GetParametersByPath", input).Return(output, nil)
	}
	ps.svc = mocksvc

	expecteds := map[string]*string{
		"VAR_A": aws.String("value_a"),
		"VAR_B": aws.String("value_b"),
		"VAR_C": aws.String("value_c"),
		"VAR_D": aws.String("value_d"),
		"VAR_E": aws.String(""),
	}
	envs, err := ps.GetEnvs()
	if assert.Nil(t, err) {
		for name, expected := range expecteds {
			assert.Equal(t, expected, envs.GetEnv(name))
		}
	}
}

func TestPutEnvs(t *testing.T) {
	emptyPattern := "empty"
	ps := New("/path", emptyPattern)

	mocksvc := new(MockSSMAPI)

	// Set the currently stored value
	mocksvc.On("GetParametersByPath", mock.Anything).Return(
		&ssm.GetParametersByPathOutput{
			Parameters: []*ssm.Parameter{
				&ssm.Parameter{Name: aws.String("VAR_B"), Value: aws.String("value_b")},
				&ssm.Parameter{Name: aws.String("VAR_C"), Value: aws.String("value_cc")},
				&ssm.Parameter{Name: aws.String("VAR_D"), Value: aws.String(emptyPattern)},
			},
		},
		nil,
	)

	// Mock values to store this time
	mocks := map[string]string{
		"VAR_A": "value_a",
		//"VAR_B": "value_b", This value has not changed
		"VAR_C": "value_cc",
		"VAR_D": "value_d",
		"VAR_E": emptyPattern,
	}
	for n, v := range mocks {
		input := &ssm.PutParameterInput {
			Name: aws.String(ps.nameWithPath(n)),
			Overwrite: aws.Bool(true),
			Type: aws.String("SecureString"),
			Value: aws.String(v),
		}
		mocksvc.On("PutParameter", input).Return(nil, nil)
	}
	mocksvc.On("PutParameter", mock.Anything).Return(nil, fmt.Errorf("PutParameter was called with invalid input"))
	ps.svc = mocksvc

	stored := map[string]string{
		"VAR_A": "value_a",
		"VAR_B": "value_b",
		"VAR_C": "value_cc",
		"VAR_D": "value_d",
		"VAR_E": "",
	}
	envs := env.New()
	for n, v := range stored {
		envs.PutEnv(n, &v)
	}
	_, err := ps.PutEnvs(envs)
	assert.Nil(t, err)
}

func TestDeleteEnvs(t *testing.T) {
	emptyPattern := "empty"
	ps := New("/path", emptyPattern)

	mocksvc := new(MockSSMAPI)
	output := &ssm.DeleteParametersOutput{
		DeletedParameters: []*string{aws.String("VAR_A"), aws.String("VAR_B")},
	}
	mocksvc.On("DeleteParameters", mock.Anything).Return(output, nil)
	ps.svc = mocksvc

	deleted := []string{"VAR_A", "VAR_B", "VAR_C"}
	actuals, err := ps.DeleteEnvs(deleted)

	expecteds := []string{"VAR_A", "VAR_B"}
	if assert.Nil(t, err) {
		assert.ElementsMatch(t, expecteds, actuals)
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
