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


func printDirectory(p *api.Directory) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\tName: %s\n\n\tLink: %s\n\n\tDescription: %s\n\n\tOfficial: %t\n\n\tCreated at: %s\n\tUpdated at: %s\n",
    href[len(href)-1], p.Name, p.Href, p.Description, *p.Official, p.CreatedAt, p.UpdatedAt)
}

func printDirectoryShort(p *api.Directory) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\tName: %s\n",
    href[len(href)-1], p.Name)
}

var directoriesCmd = &cobra.Command{
    Use: "directories",
}

var directoriesListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all directories",
    Long: `Prints all directories of current organization`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {

        lo := &api.ListOptions{
            Limit: 25,
            Page: 1,
        }

        if s := viper.GetInt("directory-limit"); s != 0 {
            lo.Limit = s
        }
        if s := viper.GetInt("directory-page"); s != 0 {
            lo.Page = s
        }

        objs, _, err := ctx.Directories.List(lo)
        if err != nil {
            log.WithError(err).Error("Failed to retrieve tenant directories")
            os.Exit(-1)
        }

        for _, obj := range objs {
            printDirectoryShort(&obj)
        }
    },
}

var directoriesShowCmd = &cobra.Command{
    Use:   "show <id>",
    Short: "Show directory details",
    Long: `Prints details of directory with specified ID`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of directory")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }

        obj, err := ctx.Directories.GetById(args[0])
        if err != nil {
            log.WithError(err).Fatal("Failed to get directory")
            os.Exit(-1)
        }

        printDirectory(obj)
    },
}


var directoriesUpdateCmd = &cobra.Command{
    Use:   "update <id>",
    Short: "Updates directory properties",
    Long: `Updates directory properties with specified ID. Fields for update
    must be specified in flags using lowercase (eg. --name "A new name")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of directory")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }
        
        obj, _ := ctx.Directories.GetById(args[0])
        if s := viper.GetString("directory-update-name"); s != "" {
            obj.Name = s
        }
        if s := viper.GetString("directory-update-description"); s != "" {
            obj.Description = s
        }

        err := obj.Save()
        if err != nil {
            log.WithError(err).Fatal("Couldn't save directory")
            return
        }
        
        printDirectory(obj)
    },
}

var directoriesCreateCmd = &cobra.Command{
    Use:   "create <name>",
    Short: "Creates new directory",
    Long: `Creates new directory. Its fields will be populated with values
    specified in flags using lowercase (eg. --name "My directory")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify name of a directory")
                return
            }
            log.Warn("You can specify only name, I'll discard other arguments")
        }
        req := &api.DirectoryRequestCreate {
            Name: args[0],
        }

        if s := viper.GetString("directory-create-description"); s != "" {
            req.Description = s
        }

        obj, err := ctx.Directories.Create(req)
        if err != nil {
            log.WithError(err).Fatal("Couldn't create directory")
            return
        }
        printDirectory(obj)
    },
}

func init() {
    directoriesCreateCmd.Flags().String("description", "", "Description of directory")
    viper.BindPFlag("directory-create-description", directoriesCreateCmd.Flags().Lookup("description"))

    directoriesUpdateCmd.Flags().String("name", "", "Name of directory")
    viper.BindPFlag("directory-update-name", directoriesUpdateCmd.Flags().Lookup("name"))
    directoriesUpdateCmd.Flags().String("description", "", "Description of directory")
    viper.BindPFlag("directory-update-description", directoriesUpdateCmd.Flags().Lookup("description"))

    directoriesListCmd.Flags().String("limit", "", "Max number of records")
    viper.BindPFlag("directory-limit", directoriesListCmd.Flags().Lookup("limit"))
    directoriesListCmd.Flags().String("page", "", "Page of records")
    viper.BindPFlag("directory-page", directoriesListCmd.Flags().Lookup("page"))

    directoriesCmd.AddCommand(directoriesListCmd)
    directoriesCmd.AddCommand(directoriesShowCmd)
    directoriesCmd.AddCommand(directoriesUpdateCmd)
    directoriesCmd.AddCommand(directoriesCreateCmd)
    RootCmd.AddCommand(directoriesCmd)
}
    

