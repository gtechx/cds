package cdsclient

import (
	"context"
	"net/url"

	"github.com/ovh/cds/sdk"
)

func (c *client) EnvironmentCreate(key string, env *sdk.Environment) error {
	if _, err := c.PostJSON(context.Background(), "/project/"+key+"/environment", env, nil); err != nil {
		return err
	}
	return nil
}

func (c *client) EnvironmentDelete(key string, envName string) error {
	if _, err := c.DeleteJSON(context.Background(), "/project/"+key+"/environment/"+url.QueryEscape(envName), nil, nil); err != nil {
		return err
	}
	return nil
}

func (c *client) EnvironmentGet(key string, envName string, mods ...RequestModifier) (*sdk.Environment, error) {
	env := &sdk.Environment{}
	if _, err := c.GetJSON(context.Background(), "/project/"+key+"/environment/"+url.QueryEscape(envName), env); err != nil {
		return nil, err
	}
	return env, nil
}

func (c *client) EnvironmentList(key string) ([]sdk.Environment, error) {
	envs := []sdk.Environment{}
	if _, err := c.GetJSON(context.Background(), "/project/"+key+"/environment", &envs); err != nil {
		return nil, err
	}
	return envs, nil
}
