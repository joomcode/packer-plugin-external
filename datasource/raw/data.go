//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package raw

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
)

type Config struct {
	Program    []string `mapstructure:"program"`
	WorkingDir string   `mapstructure:"working_dir"`
	Query      string   `mapstructure:"query"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	Result string `mapstructure:"result"`
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
	program := d.config.Program

	// first element is assumed to be an executable command, possibly found
	// using the PATH environment variable.
	_, err := exec.LookPath(program[0])
	if err != nil {
		return cty.Value{}, fmt.Errorf("external program lookup failed: %v", err)
	}

	cmd := exec.CommandContext(context.Background(), program[0], program[1:]...)
	cmd.Dir = d.config.WorkingDir
	cmd.Stdin = bytes.NewBufferString(d.config.Query)

	log.Printf("executing command %+v in %v, input='%q'", program, cmd.Dir, d.config.Query)
	resultRaw, err := cmd.Output()
	if err != nil {
		var stderr string
		if eerr, ok := err.(*exec.ExitError); ok {
			stderr = string(eerr.Stderr)
		}
		log.Printf("command failed, err=%v, stdout='%q', stderr='%q'", err, string(resultRaw), stderr)
		return cty.Value{}, fmt.Errorf("external program run failed: %v", err)
	}
	log.Println("command result:", string(resultRaw))

	output := DatasourceOutput{Result: string(resultRaw)}
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
