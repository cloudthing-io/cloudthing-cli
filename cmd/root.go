// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
    logger "github.com/apex/log"
    "github.com/apex/log/handlers/cli"
    "gitlab.com/cloudthing/go-api-client"
    "gitlab.com/cloudthing/cloudthing-cli/utils"
    "github.com/gosuri/cmdns"
)

var cfgFile string
var log logger.Interface

var ctx *api.Client

func isAuth(cmd *cobra.Command, args []string) {
    if !ctx.IsAuthenticated() {
        log.Error("You need to login first")
        os.Exit(-1)
    }
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cloudthing",
	Short: "A CLI client for CloudThing.io",
	Long: 	`A comand-line client for CloudThing.io Internet of Things cloud platform.
You may find convenient auto-configuration by typing:
    cloudthing configure
Check out our getting started guides or simply login by typing:
    cloudthing login
If you still don't have an account, just type:
    cloudthing register`,
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        var logLevel = logger.InfoLevel
        if viper.GetBool("debug") {
            logLevel = logger.DebugLevel
        }
        log = &logger.Logger{
            Level:   logLevel,
            Handler: cli.New(os.Stdout),
        }

        apiServer := viper.GetString("api-server")
        var err error

        ctx, err = api.NewClient(nil, fmt.Sprintf("http://%s", apiServer))
        if err != nil {
            log.WithError(err).Fatal("Couldn't create API client")
            return
        }

        auth := utils.LoadAuth(apiServer)

        if auth != nil {
            err = ctx.SetTokenAuth(auth)
            if err == nil {
                return
            }
        } 
        if path := viper.GetString("apikey"); path != "" {
            key, secret, err := utils.LoadApikey(path)
            if err == nil {
                err = ctx.SetBasicAuth(key, secret)
                if err == nil {
                    utils.SaveAuth(apiServer, ctx.GetToken())
                }
            }
        }
    },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

    cmdns.Namespace(RootCmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)


	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cloudthing-cli.yaml)")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

    RootCmd.PersistentFlags().String("api-server", "", "Host of API server")
    viper.BindPFlag("api-server", RootCmd.PersistentFlags().Lookup("api-server"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".cloudthing-cli") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
