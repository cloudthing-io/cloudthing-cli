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

var productsCmd = &cobra.Command{
    Use: "products",
}

var productsListCmd = &cobra.Command{
    Use:   "list",
    Short: "Prints tenants directories",
    Long: `Prints tenants directories`,
    PreRun: func(cmd *cobra.Command, args []string) {
        if !ctx.IsAuthenticated() {
            log.Error("Yu need to login first")
            os.Exit(-1)
        }
    },
    Run: func(cmd *cobra.Command, args []string) {
        tenant, _ := ctx.Tenant.Get()
        dirs, _, _ := tenant.Products()

        for _, dir := range dirs {
            id := strings.Split(dir.Href, "/")
            log.Infof("ID: %s", id[len(id)-1])
            log.Infof("Name: %s", dir.Name)
            log.Infof("Description: %s", dir.Description)
            //log.Infof("Resources: %s", dir.Description)
            log.Infof("Created at: %s", dir.CreatedAt)
            log.Infof("Updated at: %s", dir.UpdatedAt)
        }
    },
}

var productsInfoCmd = &cobra.Command{
    Use:   "info directoryId",
    Short: "Prints tenants directories",
    Long: `Prints tenants directories`,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) < 1 {
            log.Error("You need to specify at least directoryId")
            return
        }

        dir, _ := ctx.Products.GetById(args[0])

        id := strings.Split(dir.Href, "/")
        log.Infof("ID: %s", id[len(id)-1])
        log.Infof("Name: %s", dir.Name)
       // log.Infof("Official: %t", dir.Official)
        log.Infof("Description: %s", dir.Description)
       /// log.Infof("Status: %s", dir.Status)
        log.Infof("Created at: %s", dir.CreatedAt)
        log.Infof("Updated at: %s", dir.UpdatedAt)
    },
}


var productsUpdateCmd = &cobra.Command{
    Use:   "update <productId>",
    Short: "Updates directory with given ID.",
    Long: ``,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) < 1 {
            log.Error("You need to specify at least directoryId")
            return
        }
        
        dir, _ := ctx.Products.GetById(args[0])
        if s := viper.GetString("products-name"); s != "" {
            dir.Name = s
        }
        if s := viper.GetString("products-description"); s != "" {
            dir.Description = s
        }
fmt.Println(dir)
        err := dir.Save()
        if err != nil {
            log.WithError(err).Fatal("Couldn't save directory")
            return
        }
        
        productsInfoCmd.Run(productsInfoCmd, args[:1])
    },
}

var productsCreateCmd = &cobra.Command{
    Use:   "create name [description]",
    Short: "Creates new directory",
    Long: ``,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) < 1 {
            log.Error("You have to specify at least name of a new directory")
            return
        }

        req := &api.Product {
            Name: args[0],
        }
        if len(args) > 1 {
            req.Description = args[1]
        }
        dir, err := ctx.Products.Create(req)
        if err != nil {
            log.WithError(err).Fatal("Couldn't create directory")
            return
        }
        href := strings.Split(dir.Href, "/")
        arg := make([]string,1)
        arg[0] = href[len(href)-1]
        productsInfoCmd.Run(productsInfoCmd, arg)
    },
}

func init() {
    productsCmd.PersistentFlags().String("name", "", "Name of product")
    viper.BindPFlag("products-name", productsCmd.PersistentFlags().Lookup("name"))

    productsCmd.PersistentFlags().String("description", "", "Description of product")
    viper.BindPFlag("products-description", productsCmd.PersistentFlags().Lookup("description"))

    productsCmd.AddCommand(productsListCmd)
    productsCmd.AddCommand(productsInfoCmd)
    productsCmd.AddCommand(productsUpdateCmd)
    productsCmd.AddCommand(productsCreateCmd)
    RootCmd.AddCommand(productsCmd)
}
    

