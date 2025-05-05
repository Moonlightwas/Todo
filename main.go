package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"todo_app/commands"
	"todo_app/db"

	"github.com/spf13/cobra"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("todo> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" || input == "q" {
			fmt.Println("Exit the program...")
			break
		}

		args := strings.Fields(input)
		os.Args = append([]string{"todo"}, args...)

		rootCmd := &cobra.Command{
			Short:         "Todo is a CLI todo.txt format task manager",
			SilenceUsage:  true,
			SilenceErrors: true,
		}

		rootCmd.AddCommand(
			commands.AddCommand(db),
			commands.CompleteCommand(db),
			commands.DeleteCommand(db),
			commands.ListCommand(db))

		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
		}
	}
}
