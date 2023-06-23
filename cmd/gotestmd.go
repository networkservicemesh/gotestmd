package cmd

import (
	"fmt"

	"github.com/networkservicemesh/gotestmd/cmd/gen"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	gotestCmd := &cobra.Command{
		Use: "gotestmd",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Add flag --help to get commands")
			fmt.Println("See more information about NSM https://networkservicemesh.io/docs/concepts/enterprise_users/")
		},
	}

	genCmd := gen.New()
	genCmd.PersistentFlags().Bool("bash", false, "Generate a test as a bash script")

	gotestCmd.AddCommand(genCmd)

	return gotestCmd
}
