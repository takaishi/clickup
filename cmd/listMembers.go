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
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type ListMember struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type GetListMembersResponse struct {
	Members []ListMember `json:"members"`
}

// listMembersCmd represents the listMembers command
var listMembersCmd = &cobra.Command{
	Use:   "listMembers",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		members, err := getListMembers()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, member := range members {
			fmt.Printf("%d %s\n", member.ID, member.Username)
		}
	},
}

func getListMembers() ([]ListMember, error) {
	url := fmt.Sprintf("https://api.clickup.com/api/v2/list/%s/member", listId)

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

	var getListMembersResp GetListMembersResponse
	json.Unmarshal(body, &getListMembersResp)
	return getListMembersResp.Members, nil
}

func init() {
	getCmd.AddCommand(listMembersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listMembersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listMembersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listMembersCmd.Flags().StringVar(&listId, "list-id", "", "")
}
