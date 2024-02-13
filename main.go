package main

import (
	"context"
	"cookvs/application"
	"fmt"
	"os/signal"
	"os"
)


func main() {
	app := 	application.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failder to start app")
	}


}


