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
    "github.com/howeyc/gopass"
    _"net/http"
    _"encoding/json"
    "github.com/mitchellh/go-homedir"
    "path"
    "gitlab.com/cloudthing/cloudthing-cli/utils"
    "gitlab.com/cloudthing/go-api-client"
)

var AuthsFileName string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
    Use:   "login [email]",
    Short: "Authenticates user and obtains authorization yokens",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            log.Error("You must specify e-mail only")
            os.Exit(-1)
        }
        //fmt.Println(cmd)

        fmt.Print("Password: ")
        pass, err := gopass.GetPasswd()
        if err != nil {
            log.WithError(err).Fatal("Failed to get password")
            os.Exit(-1)
        }

        apiServer := viper.GetString("api-server")
        ctx, err = api.NewClient(nil, fmt.Sprintf("http://%s", apiServer))
        if err != nil {
            log.WithError(err).Fatal("Couldn't create API client")
            return
        }

        err = ctx.SetBasicAuth(args[0], string(pass))

        if err != nil {
            log.WithError(err).Fatal("Authentication failed")
            return 
        }

        err = utils.SaveAuth(apiServer, ctx.GetToken())
        if err != nil {
            log.WithError(err).Fatal("Can't save credentials to file")
            return
        }

        log.Infof("Successfully logged in! Auth file created in %s", AuthsFileName)


    },
}

func init() {
    RootCmd.AddCommand(loginCmd)

    dir, err := homedir.Dir()
    if err != nil {
        panic(err)
    }
    expanded, err := homedir.Expand(dir)
    if err != nil {
        panic(err)
    }
    AuthsFileName = path.Join(expanded, ".cloudthing-cli/auths.json")
}
    

