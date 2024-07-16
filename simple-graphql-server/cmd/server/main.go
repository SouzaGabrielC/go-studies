package main

import (
	"database/sql"
	"graphql-go-expert/graph"
	"graphql-go-expert/internal/database"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = "8080"

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("failed to open DB connection: %v\n", err)
	}
	defer db.Close()

	categoryDb, err := database.NewCategory(db)
	if err != nil {
		log.Fatalf("failed to create category table: %v\n", err)
	}

	courseDb, err := database.NewCourse(db)
	if err != nil {
		log.Fatalf("failed to create course table: %v\n", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CategoryDB: categoryDb,
		CourseDB:   courseDb,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
