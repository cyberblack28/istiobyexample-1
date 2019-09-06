package main

import (
	"fmt"
	pg "github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"log"
	"math/rand"
	"time"
)

type User struct {
	Id int
}

var db *pg.DB

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*User)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func wrt() {
	fmt.Println("writing a random number to Postgres...")
	number := rand.Intn(100000)

	p := &User{Id: number}

	err := db.Insert(p)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("✅ success")
	}
}

func startPolling() {
	for {
		<-time.After(1 * time.Second)
		go wrt()
	}
}

func main() {
	fmt.Println("⏰ starting writer...")
	// initialize postgres client

	db = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "megan",
		Addr:     "34.67.231.194:5432",
	})
	defer db.Close()

	// what is a person
	err := createSchema(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("client started...")
	go startPolling()
	time.Sleep(time.Second * 100000000)
}
