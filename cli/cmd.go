package cli

import (
	"github.com/spf13/cobra"
	"github.com/zibilal/logwrapper"
	"os"
	"simplims/cli/cmdhttp"
	"strings"
)

func Execute() {

	cmdRoot := &cobra.Command{Use: selfName()}
	cmdhttp.ExecuteServeHttp(cmdRoot)

	err := cmdRoot.Execute()
	if err != nil {
		logwrapper.Fatal(err)
	}

	selfName()
}

func selfName() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dt := strings.Split(ex, "/")
	return dt[len(dt) - 1]
}
