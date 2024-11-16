package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	// Create a shared context
	ctx := context.Background()

	// Run the stages of the pipeline
	if err := Build(ctx); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
}

func Build(ctx context.Context) error {
	// Initialize Dagger client
	client, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	LANGTRACE_API_KEY := client.SetSecret("LANGTRACE_API_KEY", os.Getenv("LANGTRACE_API_KEY"))

	python := client.Container().From("python:3.12.2-bookworm").
		WithSecretVariable("LANGTRACE_API_KEY", LANGTRACE_API_KEY).
		WithDirectory("python", client.Host().Directory("python")).
		WithExec([]string{"python", "--version"}).
		WithExec([]string{"pip", "install", "crewai"}).
		WithExec([]string{"pip", "install", "langtrace-python-sdk"})

	output, err := python.WithExec([]string{"python", "python/run-agents.py"}).Stdout(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Python script output:", output)

	// _, err = python.
	// 	Directory("output").
	// 	Export(ctx, "output")
	// if err != nil {
	// 	return err
	// }

	return nil
}
