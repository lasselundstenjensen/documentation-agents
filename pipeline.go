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

	// Declare secrets
	LANGTRACE_API_KEY := client.SetSecret("LANGTRACE_API_KEY", os.Getenv("LANGTRACE_API_KEY"))
	OPENAI_API_KEY := client.SetSecret("OPENAI_API_KEY", os.Getenv("OPENAI_API_KEY"))

	// Create a python container and install the necessary packages
	python := client.Container().From("python:3.12.2-bookworm").
		WithSecretVariable("LANGTRACE_API_KEY", LANGTRACE_API_KEY).
		WithSecretVariable("OPENAI_API_KEY", OPENAI_API_KEY).
		WithExec([]string{"python", "--version"}).
		WithExec([]string{"pip", "install", "crewai"}).
		WithExec([]string{"pip", "install", "langtrace-python-sdk"})

	// Add the python and doccing directories to the container
	python = python.
		WithDirectory("python", client.Host().Directory("python")).
		WithDirectory("doccing", client.Host().Directory("doccing"))

	// Install crewai
	output, err := python.
		WithWorkdir("doccing").
		WithExec([]string{"crewai", "install"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	fmt.Println("'crewai install' output:", output)

	// Run the crew agents
	output, err = python.
		WithWorkdir("doccing").
		WithExec([]string{"crewai", "run"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	fmt.Println("'crewai run' output:", output)

	output, err = python.
		WithExec([]string{"python", "python/run-agents.py"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Python script output:", output)

	// Print the contents of the report.md file
	output, err = python.
		WithExec([]string{"cat", "report.md"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	fmt.Println("'report.md' output:", output)

	// _, err = python.
	// 	WithWorkdir("doccing").
	// 	Directory("output").
	// 	Export(ctx, "output/report.md")
	// if err != nil {
	// 	return err
	// }

	return nil
}
