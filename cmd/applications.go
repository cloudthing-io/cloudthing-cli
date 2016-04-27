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
    _"fmt"
    _"os"
    "github.com/spf13/cobra"
    _"github.com/spf13/viper"
    _"net/http"
    _"net/url"
    _"encoding/json"
    _"gitlab.com/cloudthing/cloudthing-cli/utils"
    api "gitlab.com/cloudthing/go-api-client"
    "strings"
)

// loginCmd represents the login command
var applicationsCmd = &cobra.Command{
    Use:   "applications",
    Short: "Prints tenants directories",
    Long: `Prints tenants directories`,
    Run: func(cmd *cobra.Command, args []string) {
        tenant, _ := ctx.Tenant.Get()
        dirs, _, _ := tenant.Applications()

        for _, dir := range dirs {
            id := strings.Split(dir.Href, "/")
            log.Infof("ID: %s", id[len(id)-1])
            log.Infof("Name: %s", dir.Name)
            log.Infof("Official: %t", dir.Official)
            log.Infof("Description: %s", dir.Description)
            log.Infof("Status: %s", dir.Status)
            log.Infof("Created at: %s", dir.CreatedAt)
            log.Infof("Updated at: %s", dir.UpdatedAt)
        }
    },
}

var applicationsInfoCmd = &cobra.Command{
    Use:   "info directoryId",
    Short: "Prints tenants directories",
    Long: `Prints tenants directories`,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) < 1 {
            log.Error("You need to specify at least directoryId")
            return
        }

        dir, _ := ctx.Applications.GetById(args[0])

        id := strings.Split(dir.Href, "/")
        log.Infof("ID: %s", id[len(id)-1])
        log.Infof("Name: %s", dir.Name)
        log.Infof("Official: %t", dir.Official)
        log.Infof("Description: %s", dir.Description)
        log.Infof("Status: %s", dir.Status)
        log.Infof("Created at: %s", dir.CreatedAt)
        log.Infof("Updated at: %s", dir.UpdatedAt)
    },
}


var applicationsUpdateCmd = &cobra.Command{
    Use:   "update directoryId [name] [description] [status]",
    Short: "Updates directory with given ID.",
    Long: ``,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) < 1 {
            log.Error("You need to specify at least directoryId")
            return
        }
        
        dir, _ := ctx.Applications.GetById(args[0])
        if len(args) > 1 && args[1] != "" {
            dir.Name = args[1]
        }
        if len(args) > 2 && args[2] != "" {
            dir.Description = args[2]
        }

        if len(args) > 3 && args[3] != "" {
            if args[3] != "ENABLED" && args[3] != "DISABLED" {
                log.Warn("Status must be ENABLED or DISABLED, won't update it.")
            } else {
                dir.Status = args[3]
            }
            
        }

        err := dir.Save()
        if err != nil {
            log.WithError(err).Fatal("Couldn't save directory")
            return
        }
        
        applicationsInfoCmd.Run(applicationsInfoCmd, args[:1])
    },
}

var applicationsCreateCmd = &cobra.Command{
    Use:   "create name [description]",
    Short: "Creates new directory",
    Long: ``,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) < 1 {
            log.Error("You have to specify at least name of a new directory")
            return
        }

        req := &api.Application {
            Name: args[0],
        }
        if len(args) > 1 {
            req.Description = args[1]
        }
        dir, err := ctx.Applications.Create(req)
        if err != nil {
            log.WithError(err).Fatal("Couldn't create directory")
            return
        }
        href := strings.Split(dir.Href, "/")
        arg := make([]string,1)
        arg[0] = href[len(href)-1]
        applicationsInfoCmd.Run(applicationsInfoCmd, arg)
    },
}

func init() {
    applicationsCmd.AddCommand(applicationsInfoCmd)
    applicationsCmd.AddCommand(applicationsUpdateCmd)
    applicationsCmd.AddCommand(applicationsCreateCmd)
    RootCmd.AddCommand(applicationsCmd)
}
    

