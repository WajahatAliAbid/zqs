/*
Copyright © 2024 WajahatAliAbid
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"

	"github.com/WajahattAliAbid/zqs/helper"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type JsonInfo struct {
	Message  map[string]interface{} `json:"message"`
	FileName string                 `json:"fileName"`
}

func parseFile(filePath string) (*[]map[string]interface{}, error) {
	results := []map[string]interface{}{}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	byteValue, _ := io.ReadAll(file)
	file.Close()
	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)
	if err == nil {
		results = append(results, result)
		return &results, nil
	}

	file, _ = os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {

		var result map[string]interface{}
		err = json.Unmarshal(scanner.Bytes(), &result)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &results, nil
}

func Find(root, ext string) ([]string, error) {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a, nil
}

func chunkBy[T any](items []T, chunkSize int) [][]T {
	var _chunks = make([][]T, 0, (len(items)/chunkSize)+1)
	for chunkSize < len(items) {
		items, _chunks = items[chunkSize:], append(_chunks, items[0:chunkSize:chunkSize])
	}
	return append(_chunks, items)
}

func RunSendMessages(cmd *cobra.Command, args []string) {

	config := cmd.Context().Value(AwsConfig{}).(aws.Config)
	logger := pterm.DefaultBasicText

	fileName := cmd.Flag("file").Value.String()
	directoryName := cmd.Flag("directory").Value.String()

	if fileName == "" && directoryName == "" {
		logger.Printfln(pterm.Red("Either --file or --directory must be provided"))
		os.Exit(1)
	}

	if fileName != "" && directoryName != "" {
		logger.Printfln(pterm.Red("Only one of --file or --directory can be provided"))
		os.Exit(1)
	}

	jsons := []map[string]interface{}{}

	if fileName != "" {

		statInfo, err := os.Stat(fileName)
		if os.IsNotExist(err) {
			logger.Printfln(pterm.Red(err.Error()))
			os.Exit(1)
		}

		if statInfo.IsDir() {
			logger.Printfln(pterm.Red("File path is a directory"))
			os.Exit(1)
		}

		filePath, err := filepath.Abs(fileName)
		if err != nil {
			logger.Printfln(pterm.Red(err.Error()))
			os.Exit(1)
		}

		infos, err := parseFile(filePath)
		if err != nil {
			logger.Printfln(pterm.Red(err.Error()))
			os.Exit(1)
		}

		jsons = append(jsons, *infos...)

	}

	if directoryName != "" {
		statInfo, err := os.Stat(directoryName)
		if os.IsNotExist(err) {
			logger.Printfln(pterm.Red(err.Error()))
			os.Exit(1)
		}

		if !statInfo.IsDir() {
			logger.Printfln(pterm.Red("File path is not a directory"))
			os.Exit(1)
		}

		dirPath, err := filepath.Abs(directoryName)
		if err != nil {
			logger.Printfln(pterm.Red(err.Error()))
			os.Exit(1)
		}

		files, err := Find(dirPath, ".json")
		if err != nil {
			logger.Printfln(pterm.Red(err.Error()))
			os.Exit(1)
		}

		for _, file := range files {
			infos, err := parseFile(file)
			if err != nil {
				logger.Printfln(pterm.Red(err.Error()))
				os.Exit(1)
			}

			jsons = append(jsons, *infos...)
		}
	}
	client := helper.New(&config)
	queueName := args[0]
	queueUrlParsed, err := url.ParseRequestURI(queueName)
	var queueUrl = ""
	if err != nil {

		queueUrlFromResponse, err := client.GetQueueUrl(cmd.Context(), &queueName)
		if err != nil {
			logger.Printfln(pterm.Red(err.Error()))
			os.Exit(1)
		}
		queueUrl = *queueUrlFromResponse
	} else {
		queueUrl = queueUrlParsed.String()
	}

	progress_bar, err := pterm.DefaultProgressbar.
		WithTitle("Sending messages").
		WithTotal(len(jsons)).
		WithCurrent(0).
		WithShowElapsedTime(true).
		WithShowCount(true).
		WithRemoveWhenDone(false).
		WithBarCharacter(pterm.DefaultBarChart.HorizontalBarCharacter).
		WithShowPercentage(true).
		WithShowTitle(true).
		Start()
	if err != nil {
		logger.Printfln(pterm.Red(err.Error()))
		os.Exit(1)
	}
	chunks := chunkBy(jsons, 10)
	for _, chunk := range chunks {

		err := client.SendMessages(
			cmd.Context(),
			&queueUrl,
			&chunk,
		)

		if err != nil {
			logger.Printfln(pterm.Red(err.Error()))
			os.Exit(1)
		}

		for i := 0; i < len(chunk); i++ {
			progress_bar.Increment()
		}
	}

	logger.Printfln(pterm.Green("Messages sent successfully"))

}
