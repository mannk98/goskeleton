// Copyright © 2021 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/mannk98/goske/interfaces"
	"github.com/mannk98/goske/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os/exec"
)

var (
	goSke   interfaces.GoskeProject
	initCmd = &cobra.Command{
		Use:     "init [path]",
		Aliases: []string{"initialize", "initialise", "create"},
		Short:   "Initialize a Cobra Application",
		Long: `Initialize (cobra-cli init) will create a new application, with a license
and the appropriate structure for a Cobra-based CLI application.

Cobra init must be run inside of a go module (please run "go mod init <MODNAME>" first)
`,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var comps []string
			var directive cobra.ShellCompDirective
			if len(args) == 0 {
				comps = cobra.AppendActiveHelp(comps, "Optionally specify the path of the go module to initialize")
				directive = cobra.ShellCompDirectiveDefault
			} else if len(args) == 1 {
				comps = cobra.AppendActiveHelp(comps, "This command does not take any more arguments (but may accept flags)")
				directive = cobra.ShellCompDirectiveNoFileComp
			} else {
				comps = cobra.AppendActiveHelp(comps, "ERROR: Too many arguments specified")
				directive = cobra.ShellCompDirectiveNoFileComp
			}
			return comps, directive
		},
		Run: func(_ *cobra.Command, args []string) {
			goSke = models.NewProject()
			projectPath, err := goSke.InitializeProject(args)
			cobra.CheckErr(err)
			cobra.CheckErr(goGet("github.com/spf13/cobra"))
			if viper.GetBool("useViper") {
				cobra.CheckErr(goGet("github.com/spf13/viper"))
			}

			cobra.CheckErr(err)
			fmt.Printf("Your Cobra application is ready at\n%s\n", projectPath)
		},
	}
)

func goGet(mod string) error {
	return exec.Command("go", "get", mod).Run()
}
