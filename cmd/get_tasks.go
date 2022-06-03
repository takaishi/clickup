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
	"context"
	"encoding/json"
	"fmt"
	"github.com/raksul/go-clickup/clickup"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type GetTasksRequestParameters struct {
	ListID int `json:"list_id"`
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
		if err := renderTasks(tasks); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func renderTasks(tasks []clickup.Task) error {
	if output == "text" {
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
	} else if output == "json" {
		items := []clickup.Task{}
		if createdBy != 0 {
			for _, task := range tasks {
				if task.Creator.ID == createdBy {
					items = append(items, task)
				}
			}
			s, err := json.Marshal(items)
			if err != nil {
				return err
			}
			fmt.Println(string(s))
		} else {
			s, err := json.Marshal(tasks)
			if err != nil {
				return err
			}
			fmt.Println(string(s))
		}
	}
	return nil
}

func getTasks() ([]clickup.Task, error) {
	options := clickup.GetTasksOptions{
		Subtasks:      true,
		IncludeClosed: true,
	}

	if assignee != "" {
		options.Assignees = []string{assignee}
	}
	if assignToMe {
		resp, err := getAuthorizedUser()
		if err != nil {
			return nil, err
		}
		options.Assignees = []string{strconv.Itoa(resp.ID)}
	}
	if updatedAtGtFlag != "" {
		updatedAt, _ := time.Parse("2006-01-02", updatedAtGtFlag)
		options.DateUpdatedGt = updatedAt.Unix() * 1000
	}
	if updatedAtLtFlag != "" {
		updatedAt, _ := time.Parse("2006-01-02", updatedAtLtFlag)
		options.DateUpdatedLt = updatedAt.Unix() * 1000
	}
	if createdAtGtFlag != "" {
		createdAt, _ := time.Parse("2006-01-02", createdAtGtFlag)
		options.DateCreatedGt = createdAt.Unix() * 1000
	}
	if createdAtLtFlag != "" {
		createdAt, _ := time.Parse("2006-01-02", createdAtLtFlag)
		options.DateCreatedLt = createdAt.Unix() * 1000
	}

	if len(statuses) > 0 {
		options.Statuses = statuses
	}

	client := clickup.NewClient(nil, os.Getenv("CLICKUP_TOKEN"))
	tasks, resp, err := client.Tasks.GetTasks(context.Background(), listId, &options)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		return nil, err
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
var output string

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
	getTasksCmd.Flags().StringVar(&output, "output", "text", "")
}
