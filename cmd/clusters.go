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


func printCluster(p *api.Cluster) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\tName: %s\n\n\tLink: %s\n\n\tDescription: %s\n\n\tCreated at: %s\n\tUpdated at: %s\n",
    href[len(href)-1], p.Name, p.Href, p.Description, p.CreatedAt, p.UpdatedAt)
}

func printClusterShort(p *api.Cluster) {
    href := strings.Split(p.Href, "/")
    fmt.Printf("\n\tID: %s\n\tName: %s\n",
    href[len(href)-1], p.Name)
}

var clustersCmd = &cobra.Command{
    Use: "clusters",
}

var clustersListCmd = &cobra.Command{
    Use:   "list <applicationId>",
    Short: "List all clusters",
    Long: `Prints all clusters of current organization`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of application")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }
        lo := &api.ListOptions{
            Limit: 25,
            Page: 1,
        }

        if s := viper.GetInt("cluster-limit"); s != 0 {
            lo.Limit = s
        }
        if s := viper.GetInt("cluster-page"); s != 0 {
            lo.Page = s
        }

        objs, _, err := ctx.Clusters.ListByApplication(args[0], lo)
        if err != nil {
            log.WithError(err).Error("Failed to retrieve tenant clusters")
            os.Exit(-1)
        }

        for _, obj := range objs {
            printClusterShort(&obj)
        }
    },
}

var clustersShowCmd = &cobra.Command{
    Use:   "show <id>",
    Short: "Show cluster details",
    Long: `Prints details of cluster with specified ID`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of cluster")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }

        exp := &api.ExpandParams {
            "tenant": nil,
        }
        fmt.Println(exp)
        obj, err := ctx.Clusters.GetById(args[0], exp)
        if err != nil {
            log.WithError(err).Fatal("Failed to get cluster")
            os.Exit(-1)
        }

        printCluster(obj)
        fmt.Println(obj.Tenant)

        if viper.GetBool("cluster-show-devices") {
            exp, _, err := ctx.Devices.ListByCluster(args[0], nil)
            if err != nil {
                log.WithError(err).Error("Couldn't get devices")
                os.Exit(-1)
            }
            fmt.Printf("\tDevices:")
            for _, e := range exp {
                fmt.Printf("\t\tDevice %s\n", e.Href)

            }
        }


    },
}


var clustersUpdateCmd = &cobra.Command{
    Use:   "update <id>",
    Short: "Updates cluster properties",
    Long: `Updates cluster properties with specified ID. Fields for update
    must be specified in flags using lowercase (eg. --name "A new name")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 1 {
            if len(args) == 0 {
                log.Error("You need to specify ID of cluster")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }
        
        obj, err := ctx.Clusters.GetById(args[0])
        if err != nil {
            log.WithError(err).Error("Couldn't retrieve cluster")
            os.Exit(-1)
        }
        if s := viper.GetString("cluster-update-name"); s != "" {
            obj.Name = s
        }
        if s := viper.GetString("cluster-update-description"); s != "" {
            obj.Description = s
        }
        err = obj.Save()
        if err != nil {
            log.WithError(err).Fatal("Couldn't save cluster")
            return
        }
        
        printCluster(obj)
    },
}

var clustersAssignCmd = &cobra.Command{
    Use:   "assign <clusterId> <deviceId>",
    Short: "Updates cluster properties",
    Long: `Updates cluster properties with specified ID. Fields for update
    must be specified in flags using lowercase (eg. --name "A new name")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 2 {
            if len(args) < 2 {
                log.Error("You need to specify IDs of cluster and device")
                return
            }
            log.Warn("You can specify only ID, I'll discard other arguments")
        }
        
        dev, err := ctx.Devices.GetById(args[1])
        if err != nil {
            log.WithError(err).Error("Couldn't retrieve device")
            os.Exit(-1)
        }
        req := &api.ClusterMembershipRequestCreate {
            Device : &api.Link {
                Href: dev.Href,
            },
        }
        _, err = ctx.ClusterMemberships.CreateByCluster(args[0], req)
        if err != nil {
            log.WithError(err).Fatal("Couldn't create membership")
            return
        }
        
      //  printCluster(obj)
    },
}

var clustersCreateCmd = &cobra.Command{
    Use:   "create <applicationId> <name>",
    Short: "Creates new cluster",
    Long: `Creates new cluster. Its fields will be populated with values
    specified in flags using lowercase (eg. --name "My cluster")`,
    PreRun: isAuth,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) != 2 {
            if len(args) < 2 {
                log.Error("You need to specify ID of application and name of a cluster")
                return
            }
            log.Warn("You can specify only application and name, I'll discard other arguments")
        }
        req := &api.ClusterRequestCreate {
            Name: args[1],
        }

        if s := viper.GetString("cluster-create-description"); s != "" {
            req.Description = s
        }

        obj, err := ctx.Clusters.CreateByApplication(args[0], req)
        if err != nil {
            log.WithError(err).Fatal("Couldn't create cluster")
            return
        }
        printCluster(obj)
    },
}

func init() {
    clustersShowCmd.Flags().BoolP("exports", "e", false, "Print exported resources")
    viper.BindPFlag("cluster-show-exports", clustersShowCmd.Flags().Lookup("exports"))
    clustersShowCmd.Flags().BoolP("devices", "d", false, "Print exported resources")
    viper.BindPFlag("cluster-show-devices", clustersShowCmd.Flags().Lookup("devices"))

    clustersCreateCmd.Flags().String("description", "", "Description of cluster")
    viper.BindPFlag("cluster-create-description", clustersCreateCmd.Flags().Lookup("description"))

    clustersCreateCmd.Flags().String("status", "", "Status of cluster (ENABLED/DISABLED)")
    viper.BindPFlag("cluster-create-status", clustersCreateCmd.Flags().Lookup("status"))


    clustersUpdateCmd.Flags().String("name", "", "Name of cluster")
    viper.BindPFlag("cluster-update-name", clustersUpdateCmd.Flags().Lookup("name"))
    clustersUpdateCmd.Flags().String("description", "", "Description of cluster")
    viper.BindPFlag("cluster-update-description", clustersUpdateCmd.Flags().Lookup("description"))

    clustersUpdateCmd.Flags().String("status", "", "Status of cluster (ENABLED/DISABLED)")
    viper.BindPFlag("cluster-update-status", clustersUpdateCmd.Flags().Lookup("status"))

    clustersListCmd.Flags().String("limit", "", "Max number of records")
    viper.BindPFlag("cluster-limit", clustersListCmd.Flags().Lookup("limit"))
    clustersListCmd.Flags().String("page", "", "Page of records")
    viper.BindPFlag("cluster-page", clustersListCmd.Flags().Lookup("page"))

    clustersCmd.AddCommand(clustersListCmd)
    clustersCmd.AddCommand(clustersShowCmd)
    clustersCmd.AddCommand(clustersUpdateCmd)
    clustersCmd.AddCommand(clustersAssignCmd)
    clustersCmd.AddCommand(clustersCreateCmd)
    RootCmd.AddCommand(clustersCmd)
}
    

