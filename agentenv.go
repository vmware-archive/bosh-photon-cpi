package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/esxcloud/bosh-esxcloud-cpi/cmd"
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	"io/ioutil"
)

func createAgentEnv(ctx *cpi.Context, agentID, vmID, vmName string, networks, env map[string]interface{}) (res *cpi.AgentEnv) {
	res = &cpi.AgentEnv{
		AgentID:  agentID,
		VM:       cpi.VMSpec{Name: vmName, ID: vmID},
		Networks: networks,
		Env:      env,
		Mbus:     ctx.Config.Agent.Mbus,
		NTP:      ctx.Config.Agent.NTP,
	}
	return
}

func createEnvISO(env *cpi.AgentEnv, runner cmd.Runner) (path string, err error) {
	json, err := json.Marshal(env)
	if err != nil {
		return
	}
	envFile, err := ioutil.TempFile("", "agent-env-json")
	if err != nil {
		return
	}
	_, err = envFile.Write(json)
	if err != nil {
		return
	}
	err = envFile.Close()
	if err != nil {
		return
	}

	envISO, err := ioutil.TempFile("", "agent-env-iso")
	if err != nil {
		return
	}
	envISO.Close()
	output, err := runner.Run("genisoimage", "-o", envISO.Name(), envFile.Name())
	if err != nil {
		out := string(output[:])
		return "", errors.New(fmt.Sprintf("Failed to generate ISO for agent settings: %v\n%s", err, out))
	}
	return envISO.Name(), nil
}
