package raw

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

//go:embed test-fixtures/template.pkr.hcl
var testDatasourceHCL2Basic string

func TestAccExternalDatasource(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "external_datasource_basic_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testDatasourceHCL2Basic,
		Type:     "external-raw",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable find %s", logfile)
			}
			defer logs.Close()

			logsBytes, err := io.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			logsString := string(logsBytes)

			fooLog := `null.test: result: raboof`
			if matched := strings.Contains(logsString, fooLog); !matched {
				t.Fatalf("logs doesn't contain expected foo value %q", logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
