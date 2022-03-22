/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/raksul/go-clickup/clickup"
	"github.com/spf13/cobra"
	"os"
)

// getFoldersCmd represents the folders command
var getFoldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		folders, err := getFolders()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, folder := range folders {
			fmt.Printf("%s %s\n", folder.ID, folder.Name)
		}
	},
}

func getFolders() ([]clickup.Folder, error) {
	client := clickup.NewClient(nil, os.Getenv("CLICKUP_TOKEN"))
	folders, _, err := client.Folders.GetFolders(context.Background(), spaceId, false)
	if err != nil {
		return nil, err
	}
	return folders, nil
}

var spaceId string

func init() {
	getCmd.AddCommand(getFoldersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getFoldersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getFoldersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getFoldersCmd.Flags().StringVar(&spaceId, "space-id", "", "")
}
