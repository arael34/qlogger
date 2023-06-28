package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("starting...")

	env, err := ValidateEnvironment()
	if err != nil {
		fmt.Printf("Error loading environment: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\nloaded environment.")

	client, err := ConnectToDatabase(&env.DatabaseUrl, &env.DatabaseName)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(CloseDatabase(client, 1))
	}
	fmt.Println("pinged database.")

	logger := NewQLogger(
		&env.AuthHeader,
		client.Database(env.DatabaseName).Collection("logs"),
	)

	app, err :=
		NewAppBuilder().
			WithClient(client).
			WithLogger(logger).
			Build()

	if err != nil {
		fmt.Printf("Error building app: %v\n", err)
		os.Exit(CloseDatabase(client, 1))
	}

	app.Run()
}
