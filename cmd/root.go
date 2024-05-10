package cmd

import (
	"fmt"
	"github.com/Go-routine-4595/ude-alert/udealarm"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	config  = "conf.yml"
	version = "0.01"
)

var ConfigFile string
var ComppileDate string

var rootCmd = &cobra.Command{
	Use:     "adealarm",
	Version: version,
	Short:   "ude-alert - a simple CLI to generate times eries",
	Long: `this tool is used to simulate UDE alerts for a grafana mockup
    
with no argument it starts as servers and data will be generated on regular interval and send to the Postgres DB`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", version)
		fmt.Printf("Compiler: %s\n", runtime.Compiler)
		fmt.Printf("Platform: %s\n", runtime.GOOS)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("Compile Date: %s\n", ComppileDate)
		udealarm.StartSim(ConfigFile)
	},
}

var reverseCmd = &cobra.Command{
	Use:     "Add",
	Aliases: []string{"Add"},
	Short:   "Add equipment",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", version)
		fmt.Printf("Compiler: %s\n", runtime.Compiler)
		fmt.Printf("Platform: %s\n", runtime.GOOS)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("Compile Date: %s\n", ComppileDate)
		udealarm.AddEquipment(ConfigFile, args[0])
	},
}

func Execute(c string) {
	ComppileDate = c
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(reverseCmd)
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", "", "Config file (default is $PWD/config.yml)")
}
