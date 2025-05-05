package commands

import (
	"database/sql"
	"fmt"
	"time"
	"todo_app/db"
	"todo_app/models"

	"github.com/spf13/cobra"
)

func ListCommand(DB *sql.DB) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [all | date end-date]",
		Short: "List tasks",
		Long: `List tasks in format 'id (-p) dcr +prj @ctx dl:' with different filters:
		- list: show today's tasks (default)
		- list <date>: show tasks for specific date (YYYY-MM-DD)
		- list <start-date> <end-date>: show tasks in date range(YYYY-MM-DD)`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var tasks []models.Todo
			var err error
			if len(args) == 0 {
				today := time.Now().Format("2006-01-02")
				if tasks, err = db.GetTodoByDate(DB, today); err != nil {
					fmt.Printf("Failed to list %s\n", err)
					return
				}
			} else if len(args) == 1 && args[0] == "all" {
				if tasks, err = db.GetAllTodo(DB); err != nil {
					fmt.Printf("Failed ot list all %s\n", err)
					return
				}
			} else if len(args) == 1 {
				if tasks, err = db.GetTodoByDate(DB, args[0]); err != nil {
					fmt.Println("Invalid date format")
					return
				}
			} else if len(args) == 2 {
				if tasks, err = db.GetTodoByDates(DB, args[0], args[1]); err != nil {
					fmt.Println("Invalid date format")
					return
				}
			} else {
				fmt.Printf("Failed to list\n")
				return
			}
			for _, task := range tasks {
				db.PrintTodo(&task)
			}
		},
	}

	return cmd
}
