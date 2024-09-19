/*
Copyright © 2024 WajahatAliAbid
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sendMessagesCmd represents the sendMessages command
var sendMessagesCmd = &cobra.Command{
	Use:    "send-message",
	Short:  "Send messages to the queue",
	Long:   `Send messages to the queue, from either a single file or from a directory containing json files`,
	PreRun: PreExecute,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sendMessages called")
	},
}

func init() {
	rootCmd.AddCommand(sendMessagesCmd)
}
