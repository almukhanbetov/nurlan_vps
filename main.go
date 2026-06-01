package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:Zxcvbnm123@localhost:5432/nurlan?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Nurlan API"})
	})

	r.GET("/persons", func(c *gin.Context) {
		rows, err := pool.Query(context.Background(), `
			SELECT id, name, age
			FROM persons
			ORDER BY id
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		persons := []Person{}

		for rows.Next() {
			var p Person
			if err := rows.Scan(&p.ID, &p.Name, &p.Age); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			persons = append(persons, p)
		}

		c.JSON(http.StatusOK, persons)
	})

	r.Run(":8081")
}
