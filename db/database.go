package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"todo_app/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
	SSL      string
}

func LoadEnv() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBname:   os.Getenv("DB_NAME"),
		SSL:      os.Getenv("SSL_MODE"),
	}, nil
}

func InitDB() (*sql.DB, error) {
	cnfg, err := LoadEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env: %s", err)
	}

	cnfgParams := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cnfg.Host, cnfg.Port, cnfg.User, cnfg.Password, cnfg.DBname, cnfg.SSL)
	db, err := sql.Open("postgres", cnfgParams)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping DB: %s", err)
	}

	var t models.Todo
	var params []string
	params = models.GenerateParams(t)
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS todo_table (\n\t%s);",
		strings.Join(params, ",\n\t"))

	if _, err = db.Exec(query); err != nil { //creating table if not exists
		return nil, fmt.Errorf("failed to create table %w", err)
	}

	return db, nil
}

func AddTodo(db *sql.DB, todo *models.Todo) error {
	query := `INSERT INTO todo_table 
		(priority, createdat, description`

	values := fmt.Sprintf(") \nVALUES ('%s', '%s', '%s'",
		strings.ToUpper(string(todo.Priority)),
		todo.CreatedAt.Format("2006-01-02"),
		todo.Description)

	if todo.Project != "" {
		query += ", projects"
		values += ", '" + todo.Project + "'"
	}
	if todo.Context != "" {
		query += ", contexts"
		values += ", '" + todo.Context + "'"
	}
	if todo.Deadline.Format("2006-01-02") != "0001-01-01" {
		query += ", deadline"
		values += ", '" + todo.Deadline.Format("2006-01-02") + "'"
	}
	query += values + ");"

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func GetAllTodo(db *sql.DB) ([]models.Todo, error) {
	rows, err := db.Query(`SELECT
		id,
		priority,
		createdat,
		description,
		COALESCE(projects, ''),
		COALESCE(contexts, ''),
		deadline
		FROM todo_table 
		ORDER BY createdat;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Todo
	for rows.Next() {
		var t models.Todo
		var deadline sql.NullTime

		err := rows.Scan(
			&t.ID,
			&t.Priority,
			&t.CreatedAt,
			&t.Description,
			&t.Project,
			&t.Context,
			&deadline,
		)
		if err != nil {
			return nil, err
		}

		t.Deadline = deadline.Time

		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetTodoByDate(db *sql.DB, date string) ([]models.Todo, error) {
	query := fmt.Sprintf(`SELECT
		id,
		priority,
		createdat,
		description,
		COALESCE(projects, ''),
		COALESCE(contexts, ''),
		deadline
		FROM todo_table
		WHERE DATE(createdat) = '%s'::date
		ORDER BY createdat;`, date)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Todo
	for rows.Next() {
		var t models.Todo
		var deadline sql.NullTime

		err := rows.Scan(
			&t.ID,
			&t.Priority,
			&t.CreatedAt,
			&t.Description,
			&t.Project,
			&t.Context,
			&deadline,
		)
		if err != nil {
			return nil, err
		}
		t.Deadline = deadline.Time

		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetTodoByDates(db *sql.DB, fromdate, todate string) ([]models.Todo, error) {
	query := fmt.Sprintf(`SELECT
		id,
		priority,
		createdat,
		description,
		COALESCE(projects, ''),
		COALESCE(contexts, ''),
		deadline
		FROM todo_table
		WHERE createdat between '%s'::date AND '%s'::date
		ORDER BY createdat;`, fromdate, todate)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Todo
	for rows.Next() {
		var t models.Todo
		var deadline sql.NullTime

		err := rows.Scan(
			&t.ID,
			&t.Priority,
			&t.CreatedAt,
			&t.Description,
			&t.Project,
			&t.Context,
			&deadline,
		)
		if err != nil {
			return nil, err
		}
		t.Deadline = deadline.Time

		tasks = append(tasks, t)
	}
	return tasks, nil
}

func CompleteTodo(db *sql.DB, id string, description string) error {
	var query string

	switch {
	case id != "":
		query = fmt.Sprintf(`UPDATE todo_table 
			SET priority='X' 
			WHERE id=%s;`, id)
	case description != "":
		query = fmt.Sprintf(`UPDATE todo_table 
			SET priority='X' 
			WHERE description LIKE '%%%s';`, description)

	}

	result, err := db.Exec(query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func DeleteTodo(db *sql.DB, id string, description string) error {
	var query string

	switch {
	case id != "":
		query = fmt.Sprintf(`DELETE FROM todo_table  
			WHERE id=%s;`, id)
	case description != "":
		query = fmt.Sprintf(`DELETE FROM todo_table 
			WHERE description='%s';`, description)

	}

	result, err := db.Exec(query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func PrintTodo(todo *models.Todo) {
	str := fmt.Sprintf("%d (%s) %s",
		todo.ID,
		todo.Priority,
		todo.Description)

	if todo.Project != "" {
		str += " +" + todo.Project
	}
	if todo.Context != "" {
		str += " @" + todo.Context
	}
	if !todo.Deadline.IsZero() {
		str += " due:" + todo.Deadline.Format("2006-01-02")
	}
	fmt.Println(str)
}
