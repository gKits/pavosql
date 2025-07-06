package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v0.0.0")
		},
	}
}
