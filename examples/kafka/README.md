# Kafka Examples

This directory contains examples showing how to use the microkit Kafka adapter.

## Prerequisites

You need Docker and Docker Compose installed to run Kafka locally.

## Quick Start

1. **Start Kafka**:
   ```bash
   docker compose up -d
   ```

2. **Run the consumer** (in one terminal):
   ```bash
   cd consumer
   go run main.go
   ```

3. **Run the producer** (in another terminal):
   ```bash
   cd producer
   go run main.go
   ```

4. **Stop Kafka**:
   ```bash
   docker compose down
   ```

## What You'll See

- Producer publishes 3 messages with keys and payloads
- Consumer receives and prints each message
- Clean error handling when Kafka is unavailable

## Examples

- `producer/` - Shows how to publish messages to Kafka
- `consumer/` - Shows how to consume messages from Kafka