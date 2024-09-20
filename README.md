
# Zqs
[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v2-blue.svg)](https://opensource.org/licenses/)

Zen extensions for AWS SQS. This commandline extends aws sqs cli to make it easier to send messages to a queue. 

## Installation

Install zqs using go cli

```bash
go install github.com/ZenExtensions/zqs
```
    
## Features

- Search for queues based on partial match
- Run messages from a file
- Run messages from a directory
- Run providing queue name instead of queue url


## Usage/Examples
List all queues
```bash
zqs list-queues
```

List queues containing name dlq

```bash
zqs list-queues -q "dlq"
```

Send json file data.json to the queue my-queue

```bash
zqs send-message my-queue -f "./file.txt"
```

Send json files from directory containing json files

```bash
zqs send-message my-queue -d "./dir-of-files`"
```
## Run Locally

Clone the project

```bash
git clone https://github.com/ZenExtensions/zqs
```

Go to the project directory

```bash
cd zqs
```

Install dependencies

```bash
go mod download
```

Run the cli

```bash
go run main.go
```


## Authors

- [@WajahatAliAbid](https://www.github.com/WajahatAliAbid)

