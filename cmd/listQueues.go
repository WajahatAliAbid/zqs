/*
Copyright © 2024 WajahatAliAbid
*/
package cmd

import (
	"ZenExtensions/sqs-plus/helper"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func RunListQueues(cmd *cobra.Command, _ []string) {
	config := cmd.Context().Value(AwsConfig{}).(aws.Config)

	client := helper.New(&config)
	logger := pterm.DefaultBasicText

	queueNames, err := cmd.Flags().GetStringArray("queue-name")
	if err != nil {
		logger.Printfln(pterm.Red(err.Error()))
		os.Exit(1)
	}
	logger.Printfln("Listing queues")
	queues, err := client.ListQueuesWithAttributes(cmd.Context(), &queueNames)
	if err != nil {
		logger.Printfln(pterm.Red(err.Error()))
		os.Exit(1)
	}
	if len(*queues) == 0 {
		logger.Printfln(pterm.Yellow("No queues found"))
		return
	}

	treeChildren := []pterm.TreeNode{}

	for _, info := range *queues {
		subTree := []pterm.TreeNode{}
		subTree = append(subTree, pterm.TreeNode{
			Text: info.Url,
		})

		if info.Tags != nil {
			tagsTree := []pterm.TreeNode{}

			for key, value := range *info.Tags {
				tagsTree = append(tagsTree, pterm.TreeNode{
					Text: fmt.Sprintf("%s: %s", key, value),
				})
			}
			subTree = append(subTree, pterm.TreeNode{
				Text:     "Tags",
				Children: tagsTree,
			})
		}

		if info.Attributes != nil {
			attributesTree := []pterm.TreeNode{}

			for key, value := range *info.Attributes {
				attributesTree = append(attributesTree, pterm.TreeNode{
					Text: fmt.Sprintf("%s: %s", key, value),
				})
			}
			subTree = append(subTree, pterm.TreeNode{
				Text:     "Attributes",
				Children: attributesTree,
			})
		}
		treeChildren = append(treeChildren, pterm.TreeNode{
			Text:     info.Name,
			Children: subTree,
		})
	}

	tree := pterm.TreeNode{
		Text:     "Queues",
		Children: treeChildren,
	}
	pterm.DefaultTree.WithRoot(tree).Render()
}

// listQueuesCmd represents the listQueues command
var listQueuesCmd = &cobra.Command{
	Use:    "list-queues",
	Short:  "Lists aws queues",
	Long:   `Lists aws queues based on name or all queues if no name is provided`,
	PreRun: PreExecute,
	Run:    RunListQueues,
}

func init() {
	rootCmd.AddCommand(listQueuesCmd)

	listQueuesCmd.Flags().StringArrayP(
		"queue-name",
		"q",
		[]string{},
		"List queues with this name",
	)
}
