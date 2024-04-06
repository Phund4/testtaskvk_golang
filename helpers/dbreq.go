package helpers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	env "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func DBExec(query string, values ...interface{}) (sql.Result, int, string, error) {
	err := env.Load(".env")
	if err != nil {
		err = env.Load("../.env") // чтобы работали тесты
		if err != nil {
			return nil, 
				http.StatusInternalServerError,
				"unexpected error\n",
				fmt.Errorf("error in load environments: %v", err.Error());
		}
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, 
			http.StatusInternalServerError,
			"unexpected error\n",
			fmt.Errorf("error in connection to database: %s", err.Error());
	}
	defer db.Close()

	result, err := db.Exec(query, values...)
	if err != nil {
		return nil,
			http.StatusBadRequest,
			"incorrect values\n",
			fmt.Errorf("error in get task: %s", err.Error());
	}

	return result,
		http.StatusOK,
		"",
		nil;
}

func DBQuery(query string, values ...interface{}) (*sql.Rows, int, string, error) {
	err := env.Load(".env")
	if err != nil {
		err = env.Load("../.env") // чтобы работали тесты
		if err != nil {
			return nil, 
				http.StatusInternalServerError,
				"unexpected error\n",
				fmt.Errorf("error in load environments: %v", err.Error());
		}
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, 
			http.StatusInternalServerError,
			"unexpected error\n",
			fmt.Errorf("error in connection to database: %s", err.Error());
	}
	defer db.Close()

	result, err := db.Query(query, values...)
	if err != nil {
		return nil,
			http.StatusBadRequest,
			"incorrect values\n",
			fmt.Errorf("error in get task: %s", err.Error());
	}

	return result,
		http.StatusOK,
		"",
		nil;
}