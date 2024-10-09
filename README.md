# Go Temporal E-commerce

This is a demo project to demonstrate Temporal.io's capability to receive signal and query its current "state". This is based on [Build an eCommerce App With Temporal and Go, Part 1: Getting Started](https://learn.temporal.io/tutorials/go/build-an-ecommerce-app/build-an-ecommerce-app-with-temporal-part-1/) and incorporates concepts from [Temporal Workflow message passing - Signals, Queries, & Updates](https://docs.temporal.io/encyclopedia/workflow-message-passing).

## Prerequisites

- `make` command.
- Go 1.22.
- Temporal CLI.

For the complete guide on the environment setup, check out [Set up a local development environment for Temporal and Go](https://learn.temporal.io/getting_started/go/dev_environment/).

## Steps

### Run Unit Test(s)

1) Run `make vendor`.
2) Run `make test`.

### Run Workflow

1) Run `temporal server start-dev --db-filename your_temporal.db --ui-port 8080`.
2) Run `go run worker/main.go`.
3) Run `go run start/main.go`.
