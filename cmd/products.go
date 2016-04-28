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


func printProduct(p *api.Product) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\tName: %s\n\n\tLink: %s\n\n\tDescription: %s\n\n\tCreated at: %s\n\tUpdated at: %s\n",
    href[len(href)-1], p.Name, p.Href, p.Description, p.CreatedAt, p.UpdatedAt)
}

func printProductShort(p *api.Product) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\tName: %s\n",
    href[len(href)-1], p.Name)
}

var productsCmd = &cobra.Command{
    Use: "products",
}

var productsListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all products",
    Long: `Prints all products of current organization`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {

        lo := &api.ListOptions{
            Limit: 25,
            Page: 1,
        }

        if s := viper.GetInt("product-limit"); s != 0 {
            lo.Limit = s
        }
        if s := viper.GetInt("product-page"); s != 0 {
            lo.Page = s
        }

        objs, _, err := ctx.Products.List(lo)
        if err != nil {
            log.WithError(err).Error("Failed to retrieve tenant products")
            os.Exit(-1)
        }

        for _, obj := range objs {
            printProductShort(&obj)
        }
    },
}

var productsShowCmd = &cobra.Command{
    Use:   "show <id>",
    Short: "Show product details",
    Long: `Prints details of product with specified ID`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of product")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }

        obj, err := ctx.Products.GetById(args[0])
        if err != nil {
            log.WithError(err).Fatal("Failed to get product")
            os.Exit(-1)
        }

        printProduct(obj)
    },
}


var productsUpdateCmd = &cobra.Command{
    Use:   "update <id>",
    Short: "Updates product properties",
    Long: `Updates product properties with specified ID. Fields for update
    must be specified in flags using lowercase (eg. --name "A new name")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of product")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }
        
        obj, _ := ctx.Products.GetById(args[0])
        if s := viper.GetString("product-update-name"); s != "" {
            obj.Name = s
        }
        if s := viper.GetString("product-update-description"); s != "" {
            obj.Description = s
        }

        err := obj.Save()
        if err != nil {
            log.WithError(err).Fatal("Couldn't save directory")
            return
        }
        
        printProduct(obj)
    },
}

var productsCreateCmd = &cobra.Command{
    Use:   "create <name>",
    Short: "Creates new product",
    Long: `Creates new product. Its fields will be populated with values
    specified in flags using lowercase (eg. --name "My product")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify name of a product")
                return
            }
            log.Warn("You can specify only name, I'll discard other arguments")
        }
        req := &api.ProductRequestCreate {
            Name: args[0],
        }

        if s := viper.GetString("product-create-description"); s != "" {
            req.Description = s
        }
        obj, err := ctx.Products.Create(req)
        if err != nil {
            log.WithError(err).Fatal("Couldn't create product")
            return
        }
        printProduct(obj)
    },
}

var productsDeleteCmd = &cobra.Command{
    Use:   "delete <id>",
    Short: "Deletes product",
    Long: `Deletes product with specified ID`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of a product")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }
        err := ctx.Products.DeleteById(args[0])
        if err != nil {
            log.WithError(err).Fatal("Couldn't delete product")
            return
        }
    },
}

var productsExportCmd = &cobra.Command{
    Use:   "export <productId> <applicationId> [[type:name:read:write:grantRead:grantWrite]]",
    Short: "Creates new product",
    Long: `Creates new product. Its fields will be populated with values
    specified in flags using lowercase (eg. --name "My product")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) < 2 {
            log.Error("You need to specify ID of product and application")
        }
        prod, err := ctx.Products.GetById(args[0])
        if err != nil {
            log.WithError(err).Error("Couldn't get product")
            os.Exit(-1)
        }
        req := &api.ExportRequestCreate {  
            ModelType: "DEVICE",
            TenExpPerm: "PRIMARY",
            Product: &api.Link {
                prod.Href,
            },
        }

        for i := 2; i < len(args); i++ {
            ar := strings.Split(args[i], ":")
            if len(ar) != 6 {
                log.Warnf("Entry %d is not correct, omitting", (i-1))
                continue
            }
            entry := api.ExportEntry {
                Type: ar[0],
                Name: ar[1],
            }

            if ar[2] == "1" {
                entry.Read = true
            } else {
                entry.Read = false
            }

            if ar[3] == "1" {
                entry.Write = true
            } else {
                entry.Write = false
            }
            if ar[4] == "1" {
                entry.GrantRead = true
            } else {
                entry.GrantRead = false
            }
            if ar[5] == "1" {
                entry.GrantWrite = true
            } else {
                entry.GrantWrite = false
            }

            req.Export = append(req.Export, entry)
        }

        _, err = ctx.Exports.CreateByApplication(args[1], req)
        if err != nil {
            log.WithError(err).Fatal("Couldn't create export")
            return
        }
       // printProduct(obj)
    },
}

func init() {
    productsCreateCmd.Flags().String("description", "", "Description of product")
    viper.BindPFlag("product-create-description", productsCreateCmd.Flags().Lookup("description"))

    productsUpdateCmd.Flags().String("name", "", "Name of product")
    viper.BindPFlag("product-update-name", productsUpdateCmd.Flags().Lookup("name"))
    productsUpdateCmd.Flags().String("description", "", "Description of product")
    viper.BindPFlag("product-update-description", productsUpdateCmd.Flags().Lookup("description"))

    productsListCmd.Flags().String("limit", "", "Max number of records")
    viper.BindPFlag("product-limit", productsListCmd.Flags().Lookup("limit"))
    productsListCmd.Flags().String("page", "", "Page of records")
    viper.BindPFlag("product-page", productsListCmd.Flags().Lookup("page"))

    productsCmd.AddCommand(productsListCmd)
    productsCmd.AddCommand(productsShowCmd)
    productsCmd.AddCommand(productsUpdateCmd)
    productsCmd.AddCommand(productsCreateCmd)
    productsCmd.AddCommand(productsDeleteCmd)
    productsCmd.AddCommand(productsExportCmd)
    RootCmd.AddCommand(productsCmd)
}
    

