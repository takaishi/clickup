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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

type Folder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetFoldersResponse struct {
	Folders []Folder `json:"folders"`
}

// foldersCmd represents the folders command
var foldersCmd = &cobra.Command{
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

func getFolders() ([]Folder, error) {
	page := 0
	url := fmt.Sprintf("https://api.clickup.com/api/v2/space/%s/folder?archived=false&page=%d", spaceId, page)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", os.Getenv("CLICKUP_TOKEN"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()

	var getFoldersResp GetFoldersResponse
	json.Unmarshal(body, &getFoldersResp)
	return getFoldersResp.Folders, nil
}

var spaceId string

func init() {
	getCmd.AddCommand(foldersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// foldersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// foldersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	foldersCmd.Flags().StringVar(&spaceId, "space-id", "", "")
}
