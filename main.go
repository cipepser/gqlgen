package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cipepser/gqlgen/graph"
	"github.com/vektah/gqlgen/handler"
)

func main() {
	resolvers := &graph.Resolver{}
	http.Handle("/", handler.Playground("Todo", "/query"))
	http.Handle("/query", handler.GraphQL(graph.MakeExecutableSchema(resolvers)))

	fmt.Println("Listening on: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
