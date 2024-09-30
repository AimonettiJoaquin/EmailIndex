# Email Indexer

Email Indexer is a Go application designed to process and index large volumes of email data efficiently. It's particularly tailored for handling the Enron email dataset but can be adapted for other email datasets as well.

## Features

- Concurrent processing of email files
- Batch indexing to ZincSearch
- Configurable settings via environment variables
- Performance profiling with CPU and memory profiles

## Prerequisites

- Go 1.22 or later
- ZincSearch instance running
- Enron email dataset (or similar structured email data)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/email-indexer.git
   cd email-indexer
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Create a `config.env` file in the root directory with the following content:
   ```
   ZINC_HOST=http://localhost:4080
   ZINC_USER=admin
   ZINC_PASSWORD=Complexpass#123
   ZINC_INDEX=enronJELM
   BATCH_SIZE=500
   EMAIL_DATA_PATH=/path/to/your/email/data/
   ```
   Adjust the values according to your setup.

## Usage

Run the application with:

go run cmd/main.go


The application will process all email files in the specified directory, index them in batches to ZincSearch, and generate CPU and memory profiles.

## Configuration

You can configure the application by modifying the `config.env` file. The following settings are available:

- `ZINC_HOST`: URL of your ZincSearch instance
- `ZINC_USER`: Username for ZincSearch
- `ZINC_PASSWORD`: Password for ZincSearch
- `ZINC_INDEX`: Name of the index in ZincSearch
- `BATCH_SIZE`: Number of emails to process in each batch
- `EMAIL_DATA_PATH`: Path to the directory containing email data

## Performance Profiling

The application generates CPU and memory profiles:

- `cpu.prof`: CPU profile
- `memory.prof`: Memory profile

You can analyze these profiles using Go's pprof tool:

```
go tool pprof -http=:8080 cpu.prof
go tool pprof -http=:8080 memory.prof
```

## Project Structure

- `cmd/main.go`: Entry point of the application
- `pkg/config/config.go`: Configuration management
- `pkg/model/email.go`: Email data structure and processing logic
- `pkg/utils/utils.go`: Utility functions for file operations

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
