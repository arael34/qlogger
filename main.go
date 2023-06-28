package main

import (
	"fmt"
	"os"

	pkg "github.com/jonasiwnl/qlogger/app"
	"github.com/jonasiwnl/qlogger/types"
)

func main() {
	fmt.Println("starting...")

	env, err := pkg.ValidateEnvironment()
	if err != nil {
		fmt.Printf("Error loading environment: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\nloaded environment.")

	client, err := pkg.ConnectToDatabase(&env.DatabaseUrl, &env.DatabaseName)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(pkg.CloseDatabase(client, 1))
	}
	fmt.Println("pinged database.")

	logger := types.NewQLogger(
		&env.AuthHeader,
		client.Database(env.DatabaseName).Collection("logs"),
	)

	app, err := pkg.
		NewAppBuilder().
		WithClient(client).
		WithLogger(logger).
		Build()

	if err != nil {
		fmt.Printf("Error building app: %v\n", err)
		os.Exit(pkg.CloseDatabase(client, 1))
	}

	app.Run()
}
