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
    "gitlab.com/cloudthing/cloudthing-cli/utils"
    api "gitlab.com/cloudthing/go-api-client"
    "strings"
)

var _ = fmt.Printf


func printApikey(p *api.Apikey) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\tName: %s\n\n\tLink: %s\n\n\tDescription: %s\n\n\tKey: %s\n\n\tCreated at: %s\n\tUpdated at: %s\n",
    href[len(href)-1], p.Name, p.Href, p.Description, p.Key, p.CreatedAt, p.UpdatedAt)
}

func printApikeyShort(p *api.Apikey) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\tName: %s\n",
    href[len(href)-1], p.Name)
}

var apikeysCmd = &cobra.Command{
    Use: "apikeys",
}

var apikeysListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all apikeys",
    Long: `Prints all apikeys of current organization`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {

        lo := &api.ListOptions{
            Limit: 25,
            Page: 1,
        }

        if s := viper.GetInt("apikey-limit"); s != 0 {
            lo.Limit = s
        }
        if s := viper.GetInt("apikey-page"); s != 0 {
            lo.Page = s
        }

        objs, _, err := ctx.Apikeys.List(lo)
        if err != nil {
            log.WithError(err).Error("Failed to retrieve tenant apikeys")
            os.Exit(-1)
        }

        for _, obj := range objs {
            printApikeyShort(&obj)
        }
    },
}

var apikeysShowCmd = &cobra.Command{
    Use:   "show <id>",
    Short: "Show apikey details",
    Long: `Prints details of apikey with specified ID`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of apikey")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }

        obj, err := ctx.Apikeys.GetById(args[0])
        if err != nil {
            log.WithError(err).Fatal("Failed to get apikey")
            os.Exit(-1)
        }

        printApikey(obj)
    },
}


var apikeysUpdateCmd = &cobra.Command{
    Use:   "update <id>",
    Short: "Updates apikey properties",
    Long: `Updates apikey properties with specified ID. Fields for update
    must be specified in flags using lowercase (eg. --name "A new name")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of apikey")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }
        
        obj, _ := ctx.Apikeys.GetById(args[0])
        if s := viper.GetString("apikey-update-name"); s != "" {
            obj.Name = s
        }
        if s := viper.GetString("apikey-update-description"); s != "" {
            obj.Description = s
        }

        err := obj.Save()
        if err != nil {
            log.WithError(err).Fatal("Couldn't save apikey")
            return
        }
        
        printApikey(obj)
    },
}

var apikeysCreateCmd = &cobra.Command{
    Use:   "create <name>",
    Short: "Creates new apikey",
    Long: `Creates new apikey. Its fields will be populated with values
    specified in flags using lowercase (eg. --name "My apikey")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify name of a apikey")
                return
            }
            log.Warn("You can specify only name, I'll discard other arguments")
        }
        req := &api.ApikeyRequestCreate {
            Name: args[0],
        }

        if s := viper.GetString("apikey-create-description"); s != "" {
            req.Description = s
        }

        obj, err := ctx.Apikeys.Create(req)
        if err != nil {
            log.WithError(err).Fatal("Couldn't create apikey")
            return
        }

        path, err := utils.SaveApikey(obj.Key, obj.Secret)
        if err != nil {
            log.WithError(err).Error("Couldn't save API key/secret pair to file")
            os.Exit(-1)
        }
        log.Infof("Successfully created API key/secret pair and stored it in %s", path)

        viper.Set("apikey", path)
        log.Info("Your new API key is:")
        printApikey(obj)
    },
}

func init() {
    apikeysCreateCmd.Flags().String("description", "", "Description of apikey")
    viper.BindPFlag("apikey-create-description", apikeysCreateCmd.Flags().Lookup("description"))

    apikeysUpdateCmd.Flags().String("name", "", "Name of apikey")
    viper.BindPFlag("apikey-update-name", apikeysUpdateCmd.Flags().Lookup("name"))
    apikeysUpdateCmd.Flags().String("description", "", "Description of apikey")
    viper.BindPFlag("apikey-update-description", apikeysUpdateCmd.Flags().Lookup("description"))

    apikeysListCmd.Flags().String("limit", "", "Max number of records")
    viper.BindPFlag("apikey-limit", apikeysListCmd.Flags().Lookup("limit"))
    apikeysListCmd.Flags().String("page", "", "Page of records")
    viper.BindPFlag("apikey-page", apikeysListCmd.Flags().Lookup("page"))

    apikeysCmd.AddCommand(apikeysListCmd)
    apikeysCmd.AddCommand(apikeysShowCmd)
    apikeysCmd.AddCommand(apikeysUpdateCmd)
    apikeysCmd.AddCommand(apikeysCreateCmd)
    RootCmd.AddCommand(apikeysCmd)
}
    

