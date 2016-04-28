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


func printDevice(p *api.Device) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\n\tLink: %s\n\n\tToken: %s\n\n\tActivated: %t\n\n\tCreated at: %s\n\tUpdated at: %s\n",
    href[len(href)-1], p.Href, p.Token, *p.Activated, p.CreatedAt, p.UpdatedAt)
}

func printDeviceShort(p *api.Device) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n",
    href[len(href)-1])
}

var devicesCmd = &cobra.Command{
    Use: "devices",
}

var devicesListCmd = &cobra.Command{
    Use:   "list <productId>",
    Short: "List all devices",
    Long: `Prints all devices of current organization`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of product")
                return
            }
            log.Warn("You can specify only name, I'll discard other arguments")
        }
        lo := &api.ListOptions{
            Limit: 25,
            Page: 1,
        }

        if s := viper.GetInt("device-limit"); s != 0 {
            lo.Limit = s
        }
        if s := viper.GetInt("device-page"); s != 0 {
            lo.Page = s
        }

        objs, _, err := ctx.Devices.ListByProduct(args[0], lo)
        if err != nil {
            log.WithError(err).Error("Failed to retrieve tenant devices")
            os.Exit(-1)
        }

        for _, obj := range objs {
            printDeviceShort(&obj)
        }
    },
}

var devicesShowCmd = &cobra.Command{
    Use:   "show <id>",
    Short: "Show device details",
    Long: `Prints details of device with specified ID`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of device")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }

        obj, err := ctx.Devices.GetById(args[0])
        if err != nil {
            log.WithError(err).Fatal("Failed to get device")
            os.Exit(-1)
        }

        printDevice(obj)
    },
}

/*
var devicesUpdateCmd = &cobra.Command{
    Use:   "update <id>",
    Short: "Updates device properties",
    Long: `Updates device properties with specified ID. Fields for update
    must be specified in flags using lowercase (eg. --name "A new name")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of device")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }
        
        obj, _ := ctx.Devices.GetById(args[0])
        if s := viper.GetString("device-update-name"); s != "" {
            obj.Name = s
        }
        if s := viper.GetString("device-update-description"); s != "" {
            obj.Description = s
        }

        err := obj.Save()
        if err != nil {
            log.WithError(err).Fatal("Couldn't save directory")
            return
        }
        
        printDevice(obj)
    },
}*/

var devicesCreateCmd = &cobra.Command{
    Use:   "create <productId>",
    Short: "Creates new device",
    Long: `Creates new device of specified product.`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of product")
                return
            }
            log.Warn("You can specify only name, I'll discard other arguments")
        }

        req := &api.DeviceRequestCreate{}

        obj, err := ctx.Devices.CreateByProduct(args[0], req)
        if err != nil {
            log.WithError(err).Fatal("Couldn't create device")
            return
        }
        printDevice(obj)
    },
}

func init() {
    devicesCreateCmd.Flags().String("description", "", "Description of device")
    viper.BindPFlag("device-create-description", devicesCreateCmd.Flags().Lookup("description"))
/*
    devicesUpdateCmd.Flags().String("name", "", "Name of device")
    viper.BindPFlag("device-update-name", devicesUpdateCmd.Flags().Lookup("name"))
    devicesUpdateCmd.Flags().String("description", "", "Description of device")
    viper.BindPFlag("device-update-description", devicesUpdateCmd.Flags().Lookup("description"))
*/
    devicesListCmd.Flags().String("limit", "", "Max number of records")
    viper.BindPFlag("device-limit", devicesListCmd.Flags().Lookup("limit"))
    devicesListCmd.Flags().String("page", "", "Page of records")
    viper.BindPFlag("device-page", devicesListCmd.Flags().Lookup("page"))

    devicesCmd.AddCommand(devicesListCmd)
    devicesCmd.AddCommand(devicesShowCmd)
    //devicesCmd.AddCommand(devicesUpdateCmd)
    devicesCmd.AddCommand(devicesCreateCmd)
    RootCmd.AddCommand(devicesCmd)
}
    

