package main

import (
	"fmt"
	"os"

	"github.com/jonasiwnl/qlogger/portal"
)

func main() {
	fmt.Println("starting...")

	env, err := portal.ValidateEnvironment()
	if err != nil {
		fmt.Printf("Error loading environment: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\nloaded environment.")

	client, err := portal.ConnectToDatabase(&env.DatabaseUrl, &env.DatabaseName)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(portal.CloseDatabase(client, 1))
	}
	fmt.Println("pinged database.")

	logger := portal.NewQLogger(
		client.Database(env.DatabaseName).Collection("logs"),
		env.AllowedOrigins,
		&env.AuthHeader,
	)

	app, err :=
		portal.NewAppBuilder().
			WithClient(client).
			WithLogger(logger).
			Build()

	if err != nil {
		fmt.Printf("Error building app: %v\n", err)
		os.Exit(portal.CloseDatabase(client, 1))
	}

	app.Run()
}
