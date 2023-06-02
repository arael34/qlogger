package main

import (
	"fmt"
	"os"

	"github.com/arael34/qlogger/app"
	"github.com/arael34/qlogger/types"
)

func main() {
	fmt.Println("starting...\n")

	env, err := app.ValidateEnvironment()
	if err != nil {
		fmt.Printf("Error loading environment: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("loaded environment.")

	client, err := app.ConnectToDatabase(env.DatabaseUrl, env.DatabaseName)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(app.CloseDatabase(client, 1))
	}
	fmt.Println("pinged database.")

	logger := types.NewQLogger(env.AuthHeader, client.Database(env.DatabaseName))

	app, err := app.NewAppBuilder().WithClient(client).WithLogger(logger).Build()
	if err != nil {
		fmt.Printf("Error building app: %v\n", err)
		os.Exit(app.CloseDatabase(client, 1))
	}

	app.Run()
}
