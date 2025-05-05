package commands

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"todo_app/db"

	"github.com/spf13/cobra"
)

func CompleteCommand(DB *sql.DB) *cobra.Command {
	var (
		id          string
		description string
	)
	cmd := &cobra.Command{
		Use:   "complete [id | description]",
		Short: "Complete task",
		Long:  "Complete an existing task by id(--id) or description(--dcr)",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if id == "" && description == "" {
				fmt.Println("An id or description is reuqired")
				return
			}
			if id != "" && description != "" {
				fmt.Println("An id or description is reuqired, not both")
				return
			}
			if cmd.Flag("dcr").Changed {
				allArgs := os.Args
				if len(allArgs) > 2 {
					for i, arg := range allArgs {
						if arg == "--dcr" {
							description = strings.Join(allArgs[i+1:], " ")
							break
						}
					}
				}
			}

			if err := db.CompleteTodo(DB, id, description); err != nil {
				fmt.Printf("Failed to complete %s\n", err)
				return
			}

			fmt.Println("Task completed")
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "ID")
	cmd.Flags().StringVar(&description, "dcr", "", "Description")

	return cmd
}
