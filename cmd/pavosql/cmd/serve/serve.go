package serve

import (
	"github.com/spf13/cobra"
)

var (
	port     uint16
	filePath string
)

func Command() *cobra.Command {
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: start server
		},
	}

	serveCmd.Flags().StringVarP(&filePath, "file", "f", "/var/lib/pavosql/pavosql.db", "")
	serveCmd.Flags().Uint16VarP(&port, "port", "p", 6677, "")

	return serveCmd
}
