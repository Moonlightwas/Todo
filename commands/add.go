package commands

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"todo_app/db"
	"todo_app/models"

	"github.com/spf13/cobra"
)

func AddCommand(DB *sql.DB) *cobra.Command {
	var (
		priority string
		projects string
		contexts string
		deadline string
	)
	cmd := &cobra.Command{
		Use:   "add [description]",
		Short: "Add a new task",
		Long: `Add a new task to your todo list with optional parameters:
		- Priority (A-Z)
		- Projects (comma separated)
		- Contexts (comma separated)
		- Deadline (YYYY-MM-DD)`,
		Example: "add Buy milk -p A --proj shopping --ctx home --dl 2025-04-10",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			description := strings.Join(args, " ")

			todo := &models.Todo{
				Description: description,
				CreatedAt:   time.Now(),
			}

			if priority != "" {
				if len(priority) != 1 || priority[0] < 'A' || priority[0] > 'Z' {
					fmt.Println("Priority must be a single letter A-Z")
					return
				}
				todo.Priority = string(priority[0])
			} else {
				fmt.Println("Priority is required")
				return
			}

			if deadline != "" { //deadline parse
				if deadlineParsed, err := time.Parse("2006-01-02", deadline); err != nil {
					fmt.Println("Invalid date format")
					return
				} else {
					todo.Deadline = deadlineParsed
				}
			}

			todo.Project = projects
			todo.Context = contexts

			if err := db.AddTodo(DB, todo); err != nil {
				fmt.Println("Failed to add", err)
				return
			}

			fmt.Println("Task added")
		},
	}

	cmd.Flags().StringVarP(&priority, "priority", "p", "", "Task priority (A-Z)")
	cmd.Flags().StringVar(&projects, "prj", "", "Projects (comma separated)")
	cmd.Flags().StringVar(&contexts, "ctx", "", "Contexts (comma separated)")
	cmd.Flags().StringVar(&deadline, "dl", "", "Deadline (YYYY-MM-DD)")

	return cmd
}
