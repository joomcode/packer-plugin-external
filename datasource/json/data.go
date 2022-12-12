//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package json

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
)

type Config struct {
	Program    []string          `mapstructure:"program"`
	WorkingDir string            `mapstructure:"working_dir"`
	Query      map[string]string `mapstructure:"query"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	Result map[string]string `mapstructure:"result"`
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
		return err
	}
	if len(d.config.Program) == 0 {
		return fmt.Errorf("program must be specified")
	}
	if d.config.WorkingDir == "" {
		d.config.WorkingDir = "."
	}
	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	queryJson, err := json.Marshal(d.config.Query)
	if err != nil {
		return cty.Value{}, err
	}

	program := d.config.Program

	// first element is assumed to be an executable command, possibly found
	// using the PATH environment variable.
	_, err = exec.LookPath(program[0])
	if err != nil {
		return cty.Value{}, fmt.Errorf("external program lookup failed: %v", err)
	}

	cmd := exec.CommandContext(context.Background(), program[0], program[1:]...)
	cmd.Dir = d.config.WorkingDir
	cmd.Stdin = bytes.NewReader(queryJson)

	log.Printf("executing %+v in %v, input=%+v", program, cmd.Dir, d.config.Query)
	resultJson, err := cmd.Output()
	if err != nil {
		var stderr string
		if eerr, ok := err.(*exec.ExitError); ok {
			stderr = string(eerr.Stderr)
		}
		log.Printf("command failed, err=%v, stdout='%q', stderr='%q'", err, string(resultJson), stderr)
		return cty.Value{}, fmt.Errorf("external program run failed: %v", err)
	}
	log.Println("command result:", string(resultJson))

	result := map[string]string{}
	err = json.Unmarshal(resultJson, &result)
	if err != nil {
		return cty.Value{}, fmt.Errorf("failed parsing output of external program: %v", err)
	}

	output := DatasourceOutput{Result: result}
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
