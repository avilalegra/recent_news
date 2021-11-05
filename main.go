package main

import (
	"avilego.me/recent_news/config"
	"avilego.me/recent_news/factory"
	"avilego.me/recent_news/handler"
	"avilego.me/recent_news/persistence"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	defer func() {
		if err := persistence.Client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	monitorConfigDependantServices()

	fmt.Printf("App running at: %s\n", os.Getenv("ServerAddr"))
	fmt.Println("Mongo express running at: localhost:8081")

	log.Fatal(
		http.ListenAndServe(os.Getenv("ServerAddr"), handler.NewServerHttpHandler()),
	)
}

func monitorConfigDependantServices() {
	runServices := func(ctx context.Context) {
		go factory.Collector().Run(ctx)
		go factory.Cleaner().Run(ctx)
	}
	ctx, cancel := context.WithCancel(context.Background())
	runServices(ctx)

	go func() {
		for range config.Subject {
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			runServices(ctx)
		}
	}()
}
