# Golang, LiteFS, KeyDB Multi-Region REST API Example on Fly.io

## Overview

Example project of a multi-region Golang application using LiteFS and KeyDB, deployed on Fly.io. It features a REST API for managing records, utilizing SQLite for data storage with LiteFS for replication, and KeyDB for caching and workers.

## Features

- **REST API in Golang**: Built with Gin framework.
- **SQLite Database**: For storing records.
- **LiteFS**: To replicate SQLite across regions.
- **KeyDB**: Redis-compatible caching.
- **Fly.io Deployment**: Leveraging multi-region capabilities.

## TODO

- **ADD ASYNC WORKERS/TASKS**: currently looking into using Asynq. 
- **ADD AUTOMATED DB BACKUPS**: integrate litefs backups to s3 or other cloud long term storage soloution regularly. 

## Project Structure

- `main.go`: The main application source file.
- `Dockerfile`: Contains instructions to build the application container.
- `fly.toml`: Configuration for Fly.io deployment.
- `litefs.yml`: Configuration for LiteFS.
- `keydb.conf`: KeyDB configuration.

## Prerequisites

- Golang, Docker, and Fly.io CLI installed.
- An account on Fly.io.

## Local Development

1. Clone the repository.
2. Run `go build` in the project directory.
3. Execute `./go-rest-api`.

## Deploying to Fly.io

Follow these steps to deploy the application on Fly.io:

1. **Initialize Fly.io App**:

   ```bash
   fly launch --region dfw --no-deploy
   ```

2. **Create a Volume**:

   - This command creates a volume in the Dallas region (`dfw`).

   ```bash
   fly volumes create -r dfw --size 1 storage
   ```

3. **Attach Consul for Coordination**:

   ```bash
   fly consul attach
   ```

4. **Deploy the Application**:

   ```bash
   fly deploy
   ```

5. **Clone to Additional Regions**:
   - Clone to Miami (`mia`) and Denver (`den`) regions.
   ```bash
   fly m clone --select --region mia
   fly m clone --select --region den
   ```

## API Endpoints

- `POST /record`: Create a record.
- `GET /record/:id`: Get a record by ID.
- `GET /records`: Get all records.

## Configuration Files

- `litefs.yml`: Configures LiteFS.
- `keydb.conf`: Sets up KeyDB.

## Notes

- Adjust configurations as needed for your use case.
- Monitor performance and adjust resources on Fly.io.
