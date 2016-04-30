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


func printApplication(p *api.Application) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\tName: %s\n\n\tLink: %s\n\n\tDescription: %s\n\n\tOfficial: %t\n\n\tStatus: %s\n\n\tCreated at: %s\n\tUpdated at: %s\n",
    href[len(href)-1], p.Name, p.Href, p.Description, *p.Official, p.Status, p.CreatedAt, p.UpdatedAt)
}

func printApplicationShort(p *api.Application) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\tName: %s\n",
    href[len(href)-1], p.Name)
}

var applicationsCmd = &cobra.Command{
    Use: "applications",
}

var applicationsListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all applications",
    Long: `Prints all applications of current organization`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {

        lo := &api.ListOptions{
            Limit: 25,
            Page: 1,
        }

        if s := viper.GetInt("application-limit"); s != 0 {
            lo.Limit = s
        }
        if s := viper.GetInt("application-page"); s != 0 {
            lo.Page = s
        }

        objs, _, err := ctx.Applications.List(lo)
        if err != nil {
            log.WithError(err).Error("Failed to retrieve tenant applications")
            os.Exit(-1)
        }

        for _, obj := range objs {
            printApplicationShort(&obj)
        }
    },
}

var applicationsShowCmd = &cobra.Command{
    Use:   "show <id>",
    Short: "Show application details",
    Long: `Prints details of application with specified ID`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of application")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }

        exp := &api.ExpandParams {
            "tenant": nil,
        }
        fmt.Println(exp)
        obj, err := ctx.Applications.GetById(args[0], exp)
        if err != nil {
            log.WithError(err).Fatal("Failed to get application")
            os.Exit(-1)
        }

        printApplication(obj)
        fmt.Println(obj.Tenant)

        if viper.GetBool("application-show-exports") {
            exp, _, err := ctx.Exports.ListByApplication(args[0], nil)
            if err != nil {
                log.WithError(err).Error("Couldn;t get exports")
                os.Exit(-1)
            }

            for _, e := range exp {
                fmt.Printf("\tExport from tenant %s\n", e.TenantExp.Href)
               /* for _, item := range e.Export {
                    //fmt.Printf("\t\t%s %s")
                }*/
            }
        }


    },
}


var applicationsUpdateCmd = &cobra.Command{
    Use:   "update <id>",
    Short: "Updates application properties",
    Long: `Updates application properties with specified ID. Fields for update
    must be specified in flags using lowercase (eg. --name "A new name")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of application")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }
        
        obj, _ := ctx.Applications.GetById(args[0])
        if s := viper.GetString("application-update-name"); s != "" {
            obj.Name = s
        }
        if s := viper.GetString("application-update-description"); s != "" {
            obj.Description = s
        }
        if s := viper.GetString("application-update-status"); s != "" {
            if s != "ENABLED" && s != "DISABLED" {
                log.Warn("Status must be either ENABLED or DISABLED, omitting this field.")
            } else {
                obj.Status = s
            }        
        }

        err := obj.Save()
        if err != nil {
            log.WithError(err).Fatal("Couldn't save application")
            return
        }
        
        printApplication(obj)
    },
}

var applicationsCreateCmd = &cobra.Command{
    Use:   "create <name>",
    Short: "Creates new application",
    Long: `Creates new application. Its fields will be populated with values
    specified in flags using lowercase (eg. --name "My application")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify name of a application")
                return
            }
            log.Warn("You can specify only name, I'll discard other arguments")
        }
        req := &api.ApplicationRequestCreate {
            Name: args[0],
            Status: "ENABLED",
        }

        if s := viper.GetString("application-create-description"); s != "" {
            req.Description = s
        }
        if s := viper.GetString("application-update-status"); s != "" {
            if s != "ENABLED" && s != "DISABLED" {
                log.Warn("Status must be either ENABLED or DISABLED, omitting this field.")
            } else {
                req.Status = s
            }        
        }
        if s := viper.GetString("application-create-directory"); s != "" {
            dir, err := ctx.Directories.GetById(s)
            if err != nil {
                log.WithError(err).Error("Couldn't get directory")
                os.Exit(-1)
            }
            req.Directory = &api.Link {
                dir.Href,
            }
        }
        obj, err := ctx.Applications.Create(req)
        if err != nil {
            log.WithError(err).Fatal("Couldn't create application")
            return
        }
        printApplication(obj)
    },
}

func init() {
    applicationsShowCmd.Flags().BoolP("exports", "e", false, "Print exported resources")
    viper.BindPFlag("application-show-exports", applicationsShowCmd.Flags().Lookup("exports"))

    applicationsCreateCmd.Flags().String("description", "", "Description of application")
    viper.BindPFlag("application-create-description", applicationsCreateCmd.Flags().Lookup("description"))

    applicationsCreateCmd.Flags().String("status", "", "Status of application (ENABLED/DISABLED)")
    viper.BindPFlag("application-create-status", applicationsCreateCmd.Flags().Lookup("status"))
    applicationsCreateCmd.Flags().String("directory", "", "Status of application (ENABLED/DISABLED)")
    viper.BindPFlag("application-create-directory", applicationsCreateCmd.Flags().Lookup("directory"))

    applicationsUpdateCmd.Flags().String("name", "", "Name of application")
    viper.BindPFlag("application-update-name", applicationsUpdateCmd.Flags().Lookup("name"))
    applicationsUpdateCmd.Flags().String("description", "", "Description of application")
    viper.BindPFlag("application-update-description", applicationsUpdateCmd.Flags().Lookup("description"))

    applicationsUpdateCmd.Flags().String("status", "", "Status of application (ENABLED/DISABLED)")
    viper.BindPFlag("application-update-status", applicationsUpdateCmd.Flags().Lookup("status"))

    applicationsListCmd.Flags().String("limit", "", "Max number of records")
    viper.BindPFlag("application-limit", applicationsListCmd.Flags().Lookup("limit"))
    applicationsListCmd.Flags().String("page", "", "Page of records")
    viper.BindPFlag("application-page", applicationsListCmd.Flags().Lookup("page"))

    applicationsCmd.AddCommand(applicationsListCmd)
    applicationsCmd.AddCommand(applicationsShowCmd)
    applicationsCmd.AddCommand(applicationsUpdateCmd)
    applicationsCmd.AddCommand(applicationsCreateCmd)
    RootCmd.AddCommand(applicationsCmd)
}
    

