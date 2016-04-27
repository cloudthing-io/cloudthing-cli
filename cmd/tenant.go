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
    _ "gitlab.com/cloudthing/go-api-client"
    "time"
)

type Tenant struct {
    Href            string            `json:"href"`
    Name            string          `json:"name"`
    ShortName       string          `json:"shortName"`
    CreatedAt       *time.Time       `json:"createdAt"`
    UpdatedAt       *time.Time       `json:"updatedAt"`
    Directories     interface{}     `json:"directories"`
    Applications    interface{}     `json:"applications"`
    Products        interface{}     `json:"products"`
    Custom          interface{}     `json:"custom"`
}

// loginCmd represents the login command
var tenantCmd = &cobra.Command{
    Use:   "tenant",
    Short: "Authenticates user and obtains authorization yokens",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        tenant, err := ctx.Tenant.Get()
        if err != nil {
        	log.WithError(err).Error("Couldn't retrieve tenant")
        	return
        }

        log.Infof("Short name: %s", tenant.ShortName)
        log.Infof("Name: %s", tenant.Name)
        log.Infof("Created at: %s", tenant.CreatedAt)
        log.Infof("Updated at: %s", tenant.UpdatedAt)
    },
}

var tenantUpdateCmd = &cobra.Command{
    Use:   "update name",
    Short: "Authenticates user and obtains authorization yokens",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        tenant, err := ctx.Tenant.Get()
        if err != nil {
        	log.WithError(err).Error("Couldn't retrieve tenant")
        	return
        }

        tenant.Name = args[0]
        err = tenant.Save()
        if err != nil {
        	log.WithError(err).Error("Couldn'tupdate tenant")
        	return
        }
        tenantCmd.Run(tenantCmd,make([]string,0))
    },
}

func init() {
	tenantCmd.AddCommand(tenantUpdateCmd)
    RootCmd.AddCommand(tenantCmd)
}
    

