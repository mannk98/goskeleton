/*
Copyright © 2024 mannk khacman98@gmail.com
*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// initmodCmd represents the initmod command
var initmodCmd = &cobra.Command{
	Use:   "mod_init [module_name]",
	Short: "init golang module",

	Args: cobra.ExactArgs(1), // Ensure only one argument (module name) is provided.

	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]

		command := exec.Command("go", "mod", "init", moduleName)

		err := command.Run()
		if err != nil {
			cobra.CheckErr(fmt.Errorf("error initializing module: %v", err))
		}

		// Print the command's output if needed.
		// fmt.Println(string(output))
		fmt.Printf("Successfully initialized module %s\n", moduleName)
	},
}

var modtidyCmd = &cobra.Command{
	Use:   "mod_tidy",
	Short: "add missing and remove unused modules",

	Args: cobra.ExactArgs(0), // Ensure only one argument (module name) is provided.

	Run: func(cmd *cobra.Command, args []string) {
		command := exec.Command("go", "mod", "tidy")

		err := command.Run()
		if err != nil {
			cobra.CheckErr(fmt.Errorf("error tidy module: %v", err))
		}

		// Print the command's output if needed.
		// fmt.Println(string(output))
		fmt.Println("Successfully tidy module")
	},
}

func init() {
}
