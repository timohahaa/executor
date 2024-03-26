package main

import (
	"context"
	"executor/internal/repository"
	"fmt"
	"log"

	"github.com/timohahaa/postgres"
)

func main() {
	pg, err := postgres.New("postgres://timohahaa:timohahaa1337@localhost:5432/test")
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewCommandRepository(pg)
	entity, err := repo.CreateCommand(context.Background(), "chmod +x")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", entity)
}
