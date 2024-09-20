
# Zqs
[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v2-blue.svg)](https://opensource.org/licenses/)

Zqs makes it easier to send messages to the Amazon SQS. Just provide it file or directory containing the file and watch it run the file(s) seamlessly

## Installation

Install zqs using go cli

```bash
go install github.com/github.com/ZenExtensions/zqs
```
    
## Features

- Run messages from a file
- Run messages from a directory
- Run providing queue name instead of queue url


## Usage/Examples
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
git clone https://github.com/github.com/ZenExtensions/zqs
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

