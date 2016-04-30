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
    _"net/http"
    _"net/url"
    _"encoding/json"
    _"gitlab.com/cloudthing/cloudthing-cli/utils"
    api "gitlab.com/cloudthing/go-api-client"
    "strings"
)

var _ = fmt.Printf


func printUser(p *api.User) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\tEmail: %s\n\n\tLink: %s\n\n\tFirst name: %s\n\tSurname: %s\n\n\tLast successful login:  %s\n\tLast failed login: %s\n\n\tCreated at: %s\n\tUpdated at: %s\n",
    href[len(href)-1], p.Email, p.Href, p.FirstName, p.Surname, p.LastSuccessfulLogin, p.LastFailedLogin, p.CreatedAt, p.UpdatedAt)
}

func printUserShort(p *api.User) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\tEmail: %s\n",
    href[len(href)-1], p.Email)
}

var usersCmd = &cobra.Command{
    Use: "users",
}

var usersListCmd = &cobra.Command{
    Use:   "list <directoryId>",
    Short: "List all users",
    Long: `Prints all users of specified directory`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of directory")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }
        lo := &api.ListOptions{
            Limit: 25,
            Page: 1,
        }

        if s := viper.GetInt("user-limit"); s != 0 {
            lo.Limit = s
        }
        if s := viper.GetInt("user-page"); s != 0 {
            lo.Page = s
        }

        objs, _, err := ctx.Users.ListByDirectory(args[0], lo)
        if err != nil {
            log.WithError(err).Error("Failed to retrieve users")
            os.Exit(-1)
        }

        for _, obj := range objs {
            printUserShort(&obj)
        }
    },
}

var userCmd = &cobra.Command{
    Use:   "user",
    Short: "Show user details",
    Long: `Prints details of user with specified ID`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        obj, err := ctx.Users.GetCurrent()
        if err != nil {
            log.WithError(err).Fatal("Failed to get user")
            os.Exit(-1)
        }

        printUser(obj)
    },
}

var usersShowCmd = &cobra.Command{
    Use:   "show <id>",
    Short: "Show user details",
    Long: `Prints details of user with specified ID`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of user")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }

        obj, err := ctx.Users.GetById(args[0])
        if err != nil {
            log.WithError(err).Fatal("Failed to get user")
            os.Exit(-1)
        }

        printUser(obj)
    },
}


var usersUpdateCmd = &cobra.Command{
    Use:   "update <id>",
    Short: "Updates user properties",
    Long: `Updates user properties with specified ID. Fields for update
    must be specified in flags using lowercase (eg. --name "A new name")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of user")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }
        
        obj, _ := ctx.Users.GetById(args[0])
        if s := viper.GetString("user-update-email"); s != "" {
            obj.Email = s
        }
        if s := viper.GetString("user-update-password"); s != "" {
            obj.Password = s
        }
        if s := viper.GetString("user-update-firstname"); s != "" {
            obj.FirstName = s
        }
        if s := viper.GetString("user-update-surname"); s != "" {
            obj.Surname = s
        }

        err := obj.Save()
        if err != nil {
            log.WithError(err).Fatal("Couldn't save user")
            return
        }
        
        printUser(obj)
    },
}

var usersCreateCmd = &cobra.Command{
    Use:   "create <email> <password> <directoryId>",
    Short: "Creates new user",
    Long: `Creates new user. Its fields will be populated with values
    specified in flags using lowercase (eg. --name "My user")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 3 {
            if len(args) == 0 {
                log.Error("You need to specify user email, password and directory ID")
                return
            }
            log.Warn("You can specify only email, password and directory ID, I'll discard other arguments")
        }

        req := &api.UserRequestCreate {
            Email: args[0],
            Password: args[1],            
        }

        if s := viper.GetString("user-create-firstname"); s != "" {
            req.FirstName = s
        }
        if s := viper.GetString("user-create-surname"); s != "" {
            req.Surname = s
        }
        obj, err := ctx.Users.CreateByDirectory(args[2], req)
        if err != nil {
            log.WithError(err).Fatal("Couldn't create user")
            return
        }
        printUser(obj)
    },
}

var usersDeleteCmd = &cobra.Command{
    Use:   "delete <id>",
    Short: "Deletes user",
    Long: `Deletes user with specified ID`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of a user")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }
        err := ctx.Users.DeleteById(args[0])
        if err != nil {
            log.WithError(err).Fatal("Couldn't delete user")
            return
        }
    },
}
func init() {
    usersCreateCmd.Flags().String("firstname", "", "User's first name")
    viper.BindPFlag("user-create-firstname", usersCreateCmd.Flags().Lookup("firstname"))
    usersCreateCmd.Flags().String("surname", "", "Users's surname")
    viper.BindPFlag("user-create-surname", usersCreateCmd.Flags().Lookup("surname"))

    usersUpdateCmd.Flags().String("firstname", "", "User's first name")
    viper.BindPFlag("user-update-firstname", usersUpdateCmd.Flags().Lookup("firstname"))
    usersUpdateCmd.Flags().String("surname", "", "Users's surname")
    viper.BindPFlag("user-update-surname", usersUpdateCmd.Flags().Lookup("surname"))
    usersUpdateCmd.Flags().String("email", "", "User's email address")
    viper.BindPFlag("user-update-email", usersUpdateCmd.Flags().Lookup("email"))
    usersUpdateCmd.Flags().String("password", "", "Users's password")
    viper.BindPFlag("user-update-password", usersUpdateCmd.Flags().Lookup("password"))

    usersListCmd.Flags().String("limit", "", "Max number of records")
    viper.BindPFlag("user-limit", usersListCmd.Flags().Lookup("limit"))
    usersListCmd.Flags().String("page", "", "Page of records")
    viper.BindPFlag("user-page", usersListCmd.Flags().Lookup("page"))

    usersCmd.AddCommand(usersListCmd)
    usersCmd.AddCommand(usersShowCmd)
    usersCmd.AddCommand(usersUpdateCmd)
    usersCmd.AddCommand(usersCreateCmd)
    usersCmd.AddCommand(usersDeleteCmd)
    RootCmd.AddCommand(usersCmd)
    RootCmd.AddCommand(userCmd)
}
    

