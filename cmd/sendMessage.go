/*
Copyright © 2024 WajahatAliAbid
*/
package cmd

import (
	"ZenExtensions/sqs-plus/helper"
	"bufio"
	"encoding/json"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type JsonInfo struct {
	Message  map[string]interface{} `json:"message"`
	FileName string                 `json:"fileName"`
}

func parseFile(filePath string) (*[]JsonInfo, error) {
	results := []JsonInfo{}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)
	if err == nil {
		results = append(results, JsonInfo{
			Message:  result,
			FileName: filePath,
		})
		return &results, nil
	}

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		var result map[string]interface{}
		err = json.Unmarshal(scanner.Bytes(), &result)
		if err != nil {
			return nil, err
		}

		results = append(results, JsonInfo{
			Message:  result,
			FileName: filePath,
		})
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

func RunSendMessages(cmd *cobra.Command, _ []string) {
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

	jsons := []JsonInfo{}

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

		infos, err := parseFile(fileName)
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

		files, err := Find(directoryName, ".json")
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
	queueName := cmd.Flag("queue-name").Value.String()
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

	logger.Printfln(queueUrl)

}

// sendMessagesCmd represents the sendMessages command
var sendMessagesCmd = &cobra.Command{
	Use:    "send-message",
	Short:  "Send messages to the queue",
	Long:   `Send messages to the queue, from either a single file or from a directory containing json files`,
	PreRun: PreExecute,
	Run:    RunSendMessages,
}

func init() {
	sendMessagesCmd.Flags().StringP(
		"queue-name",
		"q",
		"",
		"Name of the queue",
	)

	sendMessagesCmd.Flags().StringP(
		"file",
		"f",
		"",
		"File to send",
	)

	sendMessagesCmd.Flags().StringP(
		"directory",
		"d",
		"",
		"Directory to send",
	)

	rootCmd.AddCommand(sendMessagesCmd)
}
