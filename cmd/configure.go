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
    _"github.com/howeyc/gopass"
    _"net/http"
    _"encoding/json"
    "bufio"
    "gopkg.in/yaml.v2"
    "path"
    "io/ioutil"
    "github.com/mitchellh/go-homedir"
)

// loginCmd represents the login command
var configureCmd = &cobra.Command{
    Use:   "configure",
    Short: "Configures cloudthing CLI",
    Long: `You can configure cloudthing CLI, set your server API URL and credentials.
    It is required if you want to use cloudthing with API key.`,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 0 {
            log.Warn("You're not supposed to pass any arguments, but we'll go on.")
        }

        reader := bufio.NewReader(os.Stdin)
        conf := make(map[string]interface{})

        log.Infof("Do you have an account? [y/n] (y)")
        fmt.Printf("\t")
        answerWithNewline, _ := reader.ReadString('\n')
        answer := answerWithNewline[:(len(answerWithNewline)-1)]

        if answer == "y" || answer == "Y" || answer == "" {
            log.Info("Great! Please enter URL of your API without scheme (typically it would be short-name.cloudthing.io):")
        
            fmt.Printf("\t")
            answerWithNewline, _ = reader.ReadString('\n')
            answer = answerWithNewline[:(len(answerWithNewline)-1)]
            conf["api-server"] = answer

            if answer != "" {
                log.Info("Please enter API key:")
            }
            fmt.Printf("\t")
            answerWithNewline, _ = reader.ReadString('\n')
            answer = answerWithNewline[:(len(answerWithNewline)-1)]
            conf["api"] = make(map[string]interface{})
            conf["api"].(map[string]interface{})["key"] = answer 

            if answer != "" {
                log.Info("Please enter API secret:")
            }
            fmt.Printf("\t")
            answerWithNewline, _ = reader.ReadString('\n')
            answer = answerWithNewline[:(len(answerWithNewline)-1)]
            conf["api"].(map[string]interface{})["secret"] = answer 




        } else if answer == "n" || answer == "N" {
            log.Info("No problem! I'll run cloudthing register for your convenience")
            log.Info("Please enter your email:")
            fmt.Printf("\t")
            answerWithNewline, _ = reader.ReadString('\n')
            answer = answerWithNewline[:(len(answerWithNewline)-1)]
            arg := make([]string, 1)
            arg[0] = answer
            registerCmd.Run(registerCmd, arg)

            if viper.GetString("api-server") == "" {
                log.Error("Something went wrong, sorry :(")
                return
            }

            log.Info("Tenant created.")
            conf["api-server"] = viper.GetString("api-server")


        } else {
            log.Error("This is not the answer I was expected, sorry.")
            return
        }

        d, _ := yaml.Marshal(&conf)
        cpath := viper.GetString("config")
        if cpath == "" {
                dir, err := homedir.Dir()
            if err != nil {
                panic(err)
            }
            expanded, err := homedir.Expand(dir)
            if err != nil {
                panic(err)
            }
            cpath = path.Join(expanded, ".cloudthing-cli.yaml")
        }
        
        if err := os.MkdirAll(path.Dir(cpath), 0755); err != nil {
            return
        }
        if err := ioutil.WriteFile(cpath, d, 0755); err != nil {
            return
        }

        log.Infof("Successfully created configuration in %s", cpath)
    },
}

func init() {
    RootCmd.AddCommand(configureCmd)


    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // loginCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    

}

