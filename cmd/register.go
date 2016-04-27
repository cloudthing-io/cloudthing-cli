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
    "net/http"
    "encoding/json"
    "bytes"
)

type RegisterRequest struct {
    Email       string      `json:"email"`
    FirstName   string      `json:"firstName,omitempty"`
    Surname     string      `json:"surname,omitempty"`
    OrgName     string      `json:"orgName, omitempty"`
    Password    string      `json:"password"`
}

// loginCmd represents the login command
var registerCmd = &cobra.Command{
    Use:   "register email [firstName] [surname] [organizationName]",
    Short: "Register new tenant (organization) and creates new user within it",
    Long: `This command will register new tenant (organization) on CloudThing IoT platform,
perform all necessary setup and create owning user for specified email and (optionally) names.
Please note, that this command will not create API key/secret pair and store it in config file.
To do so, please use our wizard by typing:
    cloudthing configure`,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            log.Error("You must specify at least email")
            return
        }
        rreq := RegisterRequest{
            Email: args[0],
        }

        if len(args) > 1 {
            rreq.FirstName = args[1]
        }
        if len(args) > 2 {
            rreq.Surname = args[2]
        }
        if len(args) > 3 {
            rreq.OrgName = args[3]
        }

        fmt.Print("Password: ")
        pass1, err := gopass.GetPasswd()
        if err != nil {
            log.WithError(err).Fatal("Failed to get password")
            return
        }

        fmt.Print("Confirm password: ")
        pass2, err := gopass.GetPasswd()
        if err != nil {
            log.WithError(err).Fatal("Failed to get password")
            return
        }

        if string(pass1) != string(pass2) {
            log.WithError(err).Fatal("Passwords do not match!")
            return
        }
        rreq.Password = string(pass1)

        enc, _ := json.Marshal(&rreq)

        buf := bytes.NewBuffer(enc)
        
        url := fmt.Sprintf("http://localhost/api/v1/tenants")
        req, err := http.NewRequest("POST", url, buf)
        if err != nil {
            log.WithError(err).Fatal("Failed to create request")
            return
        }

       // req.SetBasicAuth(args[0], string(pass))
        req.Header.Add("Accept", "application/json")
        //req.Host = "pink-cat.cloudthing.io"
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            log.WithError(err).Fatal("Failed to complete request")
            os.Exit(-1)
        }

        defer resp.Body.Close()

        if resp.StatusCode != http.StatusCreated {
            log.Fatal("Failed to create tenant")
            os.Exit(-1)
        }

        dec := json.NewDecoder(resp.Body)
        m := make(map[string]interface{})
        dec.Decode(&m)

        log.Infof("Your tenant host is: %s", m["href"].(string))
        log.Infof("You can now log in with flag --api-server=%s", m["href"].(string)[8:])

        viper.Set("api-server", m["href"].(string)[8:])

    },
}

func init() {
    RootCmd.AddCommand(registerCmd)


    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // loginCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    

}

