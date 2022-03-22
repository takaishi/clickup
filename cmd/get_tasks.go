/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type GetTasksRequestParameters struct {
	ListID int `json:"list_id"`
}

type Task struct {
	Name    string  `json:"name"`
	URL     string  `json:"url"`
	Creator Creator `json:"creator"`
}

type Creator struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type GetTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

// getTasksCmd represents the tasks command
var getTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := getTasks()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if createdBy != 0 {
			for _, task := range tasks {
				if task.Creator.ID == createdBy {
					fmt.Printf("- [%s](%s)\n", task.Name, task.URL)
				}
			}
		} else {
			for _, task := range tasks {
				fmt.Printf("- [%s](%s)\n", task.Name, task.URL)
			}
		}
	},
}

func getTasks() ([]Task, error) {
	tasks := []Task{}
	page := 0

	for {
		endpoint := fmt.Sprintf("https://api.clickup.com/api/v2/list/%s/task", listId)
		queryArr := []string{}
		queryMap := map[string]string{
			"page":           strconv.Itoa(page),
			"subtasks":       "true",
			"include_closed": "true",
		}
		if assignee != "" {
			queryMap["assignees[]"] = assignee
		}
		if assignToMe {
			resp, err := getAuthorizedUser()
			if err != nil {
				return nil, err
			}
			queryMap["assignees[]"] = strconv.Itoa(resp.User.ID)
		}
		if updatedAtGtFlag != "" {
			updatedAt, _ := time.Parse("2006-01-02", updatedAtGtFlag)
			queryMap["date_updated_gt"] = fmt.Sprintf("%d", updatedAt.Unix()*1000)
		}
		if updatedAtLtFlag != "" {
			updatedAt, _ := time.Parse("2006-01-02", updatedAtLtFlag)
			queryMap["date_updated_lt"] = fmt.Sprintf("%d", updatedAt.Unix()*1000)
		}
		if createdAtGtFlag != "" {
			createdAt, _ := time.Parse("2006-01-02", createdAtGtFlag)
			queryMap["date_created_gt"] = fmt.Sprintf("%d", createdAt.Unix()*1000)
		}
		if createdAtLtFlag != "" {
			createdAt, _ := time.Parse("2006-01-02", createdAtLtFlag)
			queryMap["date_created_lt"] = fmt.Sprintf("%d", createdAt.Unix()*1000)
		}

		for k, v := range queryMap {
			queryArr = append(queryArr, fmt.Sprintf("%s=%s", k, v))
		}

		if len(statuses) > 0 {
			for _, v := range statuses {
				queryArr = append(queryArr, fmt.Sprintf("statuses[]=%s", v))
			}
		}
		client := &http.Client{}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", endpoint, url.PathEscape(strings.Join(queryArr, "&"))), nil)
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

		var getTaskResp GetTasksResponse
		json.Unmarshal(body, &getTaskResp)

		if len(getTaskResp.Tasks) == 0 {
			break
		}
		tasks = append(tasks, getTaskResp.Tasks...)
		page = page + 1
	}
	return tasks, nil
}

var statuses []string
var updatedAtGtFlag string
var updatedAtLtFlag string
var createdAtGtFlag string
var createdAtLtFlag string
var createdBy int
var assignee string
var assignToMe bool
var listId string

func init() {
	getCmd.AddCommand(getTasksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getTasksCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getTasksCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getTasksCmd.Flags().String("list", "", "A help for foo")

	getTasksCmd.Flags().StringArrayVar(&statuses, "status", []string{}, "")
	getTasksCmd.Flags().StringVar(&updatedAtGtFlag, "updated-at-gt", "", "")
	getTasksCmd.Flags().StringVar(&updatedAtLtFlag, "updated-at-lt", "", "")
	getTasksCmd.Flags().StringVar(&createdAtGtFlag, "created-at-gt", "", "")
	getTasksCmd.Flags().StringVar(&createdAtLtFlag, "created-at-lt", "", "")
	getTasksCmd.Flags().StringVar(&assignee, "assignee-id", "", "")
	getTasksCmd.Flags().StringVar(&listId, "list-id", "", "")
	getTasksCmd.Flags().IntVar(&createdBy, "created-by", 0, "")
	getTasksCmd.Flags().BoolVar(&assignToMe, "assign-to-me", false, "")
}
