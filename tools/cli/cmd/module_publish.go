package cmd

import (
	"fmt"
	module2 "github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"github.com/terrariumcloud/terrarium/tools/cli/pkg/module"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type maturityValue struct {
	maturity module2.Maturity
}

var (
	moduleMetadata = module.Metadata{
		Name:        "",
		Version:     "",
		Description: "",
		Source:      "",
		Maturity:    module2.Maturity_STABLE,
	}
	maturityVar = maturityValue{maturity: module2.Maturity_STABLE}
)

// modulePublishCmd represents the publish command
var modulePublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publishes a zip file as a module version in terrarium.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conn, client, err := getModulePublisherClient()
		if err != nil {
			printErrorAndExit("Failed to connect to terrarium", err, 1)
		}
		defer func() { _ = conn.Close() }()
		if source, err := os.Open(args[0]); err != nil {
			printErrorAndExit(fmt.Sprintf("Failed to open %s", args[0]), err, 1)
		} else {
			moduleMetadata.Maturity = maturityVar.maturity
			err := module.Publish(client, source, moduleMetadata)
			if err != nil {
				printErrorAndExit("Publishing module failed", err, 1)
			} else {
				fmt.Println("Module published.")
			}
		}
	},
}

func init() {
	moduleCmd.AddCommand(modulePublishCmd)

	modulePublishCmd.Flags().StringVar(&moduleMetadata.Name, "name", "", "Name of the module in the form \"<organisation>/<name>/<provider>\".")
	_ = modulePublishCmd.MarkFlagRequired("name")
	modulePublishCmd.Flags().StringVar(&moduleMetadata.Version, "version", "", "Semantic version of the module to publish.")
	_ = modulePublishCmd.MarkFlagRequired("version")
	modulePublishCmd.Flags().StringVar(&moduleMetadata.Description, "description", "", "Description of the module.")
	modulePublishCmd.Flags().StringVar(&moduleMetadata.Source, "source", "", "URL containing the module source.")
	modulePublishCmd.Flags().Var(&maturityVar, "maturity", "The maturity of the module, one of IDEA, PLANNING, DEVELOPING, ALPHA, BETA, STABLE, DEPRECATED or END-OF-LIFE.")
}

var maturityToString = map[module2.Maturity]string{
	module2.Maturity_IDEA:        "IDEA",
	module2.Maturity_PLANNING:    "PLANNING",
	module2.Maturity_DEVELOPING:  "DEVELOPING",
	module2.Maturity_ALPHA:       "ALPHA",
	module2.Maturity_BETA:        "BETA",
	module2.Maturity_STABLE:      "STABLE",
	module2.Maturity_DEPRECATED:  "DEPRECATED",
	module2.Maturity_END_OF_LIFE: "END-OF-LIFE",
}

func (m maturityValue) String() string {
	if s, ok := maturityToString[m.maturity]; ok {
		return s
	}
	return "UNKNOWN"
}

func (m *maturityValue) Set(s string) error {
	for maturity, maturityStr := range maturityToString {
		if strings.ToUpper(s) == maturityStr {
			m.maturity = maturity
			return nil
		}
	}
	return fmt.Errorf("unknown maturity: %s", s)
}

func (m maturityValue) Type() string {
	return "maturity"
}
