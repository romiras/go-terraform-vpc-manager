/*
Copyright © 2020 Roman Miro

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/romiras/go-terraform-vpc-manager/internal/commands"
	"github.com/romiras/go-terraform-vpc-manager/internal/helpers"
	registry "github.com/romiras/go-terraform-vpc-manager/internal/registries"
	"github.com/spf13/viper"
)

const AppName = "go-terraform-vpc-manager"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-terraform-vpc-manager",
	Short: "Manages VPC with Terraform",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func getInstanceName(cmd *cobra.Command, args []string) string {
	instanceName, err := cmd.Flags().GetString("instance_name")
	if err != nil {
		panic(err)
	}
	helpers.DebugMsg("\tinstanceName=", instanceName)
	return instanceName
}

func getWorkingDir(cmd *cobra.Command, args []string) string {
	workingDir, err := cmd.Flags().GetString("working_dir")
	if err != nil {
		panic(err)
	}
	helpers.DebugMsg("\tworkingDir=", workingDir)
	return workingDir
}

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		helpers.DebugMsg("cmdCreate args: " + strings.Join(args, " "))

		instanceName := getInstanceName(cmd, args)
		workingDir := getWorkingDir(cmd, args)
		err := validate(instanceName, workingDir)
		helpers.AbortOnError(err)

		registry.Reg.WorkingDir = workingDir

		return commands.CreateVPC(instanceName)
	},
}

var cmdDestroy = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		helpers.DebugMsg("cmdDestroy args: " + strings.Join(args, " "))

		instanceName := getInstanceName(cmd, args)
		workingDir := getWorkingDir(cmd, args)
		err := validate(instanceName, workingDir)
		helpers.AbortOnError(err)

		registry.Reg.WorkingDir = workingDir

		return commands.DestroyVPC(instanceName)
	},
}

func validate(instanceName, workingDir string) error {
	if instanceName == "" {
		return errors.New("Got empty instance name")
	}

	if workingDir == "" {
		return errors.New("Got empty working directory")
	}

	return nil
}

func setSubCmdFlags(cmd *cobra.Command) {
	cmdFlags := cmd.Flags()
	cmdFlags.String("instance_name", "", "Name of instance")
	cmdFlags.String("working_dir", "", "Working directory of Terraform HCL files")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(cmdCreate, cmdDestroy)
	setSubCmdFlags(cmdCreate)
	setSubCmdFlags(cmdDestroy)

	err := rootCmd.Execute()
	helpers.AbortOnError(err)
}

func init() {
	initConfig()
	registry.Reg.ExecPath = viper.GetStringMapString("terraform")["exec_path"]

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/"+AppName+"/settings.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		helpers.AbortOnError(err)

		// Search config in home directory with name ".go-terraform-vpc-manager" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".config", AppName))
		viper.SetConfigName("settings")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	helpers.AbortOnError(err)
}
