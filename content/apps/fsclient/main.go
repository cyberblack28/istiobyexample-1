package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var plants = []string{"sunflower", "succulent", "christmas tree", "oak", "herb"}

func wrt() {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "mokeefe")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("initialized")
	}

	plantName := plants[rand.Intn(len(plants))]

	log.Println("writing a new plant to Firestore...")
	number := rand.Intn(100000)
	_, _, err = client.Collection("plants").Add(ctx, map[string]interface{}{
		"name": plantName,
		"id":   number,
	})
	if err != nil {
		log.Printf("🚫 Failed adding plant: %v", err)
	} else {
		log.Println("✅success")
	}
}

func startPolling() {
	for {
		<-time.After(3 * time.Second)
		go wrt()
	}
}

func main() {
	fmt.Println("🔥 starting writer...")
	go startPolling()
	time.Sleep(time.Second * 100000000)
}
