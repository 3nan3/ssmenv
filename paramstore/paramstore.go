package paramstore

import (
	"path/filepath"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/3nan3/ssmenv/env"
)

type ParameterStore struct {
	path string
	emptyPattern string
	sess *session.Session
	svc  ssmiface.SSMAPI
}

func New(path string, emptyPattern string) *ParameterStore {
	sess := session.New()
	svc := ssm.New(sess)
	return &ParameterStore{path, emptyPattern, sess, svc}
}

func (ps *ParameterStore) GetEnv(envName string) (*env.Env, error) {
	input := &ssm.GetParameterInput {
		Name: aws.String(ps.nameWithPath(envName)),
		WithDecryption: aws.Bool(true),
	}
	res, err := ps.svc.GetParameter(input)
	if err != nil {
		switch err.(type) {
		case *ssm.ParameterNotFound:
			return env.New(), nil
		default:
			return nil, err
		}
	}

	envs := env.New()
	ps.putParamsInEnvs(envs, []*ssm.Parameter{res.Parameter})
	return envs, nil
}

func (ps *ParameterStore) GetEnvs() (*env.Env, error) {
	envs := env.New()

	nextToken := ""
	input := &ssm.GetParametersByPathInput {
		Path: aws.String(ps.path),
		WithDecryption: aws.Bool(true),
	}

	for {
		if nextToken != "" {
			input.SetNextToken(nextToken)
		}
		res, err := ps.svc.GetParametersByPath(input)
		if err != nil {
			return nil, err
		}

		ps.putParamsInEnvs(envs, res.Parameters)
		if res.NextToken == nil {
			break
		}
		nextToken = *res.NextToken
	}
	return envs, nil
}

func (ps *ParameterStore) PutEnvs(envs *env.Env) (*env.Env, error) {
	oldenvs, err := ps.GetEnvs()
	if err != nil {
		return nil, err
	}
	for name, value := range envs.GetEnvs() {
		if value == nil || (oldenvs.GetEnv(name) != nil && *value == *oldenvs.GetEnv(name)) {
			continue
		}
		if *value == "" {
			empty := ps.emptyPattern
			value = &empty
		}

		input := &ssm.PutParameterInput {
			Name: aws.String(ps.nameWithPath(name)),
			Overwrite: aws.Bool(true),
			Type: aws.String("SecureString"),
			Value: aws.String(*value),
		}
		_, err := ps.svc.PutParameter(input)
		if err != nil {
			return nil, err
		}
	}
	return oldenvs, nil
}

func (ps *ParameterStore) DeleteEnvs(names []string) ([]string, error) {
	params := []string{}
	for _, name := range names {
		params = append(params, ps.nameWithPath(name))
	}
	input := &ssm.DeleteParametersInput {
		Names: aws.StringSlice(params),
	}

	res, err := ps.svc.DeleteParameters(input)
	if err != nil {
		return nil, err
	}
	deleted := []string{}
	for _, d := range res.DeletedParameters {
		deleted = append(deleted, nameWithoutPath(*d))
	}
	return deleted, nil
}

func (ps *ParameterStore) putParamsInEnvs(envs *env.Env, params []*ssm.Parameter) {
	for _, param := range params {
		if *param.Value == ps.emptyPattern {
			empty := ""
			envs.PutEnv(nameWithoutPath(*param.Name), &empty)
		} else {
			envs.PutEnv(nameWithoutPath(*param.Name), param.Value)
		}
	}
}

func (ps *ParameterStore) nameWithPath(envName string) string {
	return filepath.Join(ps.path, envName)
}

func nameWithoutPath(paramName string) string {
	return filepath.Base(paramName)
}
