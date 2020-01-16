package cmdhttp

import (
	"github.com/spf13/cobra"
	"github.com/zibilal/logwrapper"
	"simplims/bootstrap/httpserverbootstrap"
)

func ExecuteServeHttp(rootCmd *cobra.Command) {
	var (
		address  string
		cmdServe = &cobra.Command{
			Use:   "serve",
			Short: "listen and serve connection from client",
			Long:  "listen and serve connection from client",
			Args:  cobra.MinimumNArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				server := httpserverbootstrap.NewHttpServerBootstrap()
				err := server.Init()
				if err != nil {
					logwrapper.Fatal(err)
				}
				err = server.Run()
				if err != nil {
					logwrapper.Fatal(err)
				}
			},
		}
	)

	cmdServe.Flags().StringVarP(&address, "address", "a", ":8080", "the server address")

	rootCmd.AddCommand(cmdServe)
}
