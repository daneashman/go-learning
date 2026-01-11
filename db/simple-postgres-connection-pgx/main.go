package main

import (
	"fmt"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	db, err := sql.Open("pgx", "postgresql://localhost:5432/daneashman")
	if err != nil {
		fmt.Printf("Err opening db: %v\n", err)
		return;
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM authors")
	if err != nil {
		fmt.Printf("Err sending query: %v\n", err)
		return;
	}

	hasNext := rows.Next()
	for hasNext {
		var name string
		var age int
		scanError := rows.Scan(&name, &age)
		if scanError != nil {
			fmt.Printf("Err scanning rows: %v\n", scanError)
			return;
		}

		fmt.Printf("name: %s, age: %d\n", name, age)

		hasNext = rows.Next()
	}

	fmt.Printf("No next row. Rows.Err = %v\n", rows.Err()) 
}
