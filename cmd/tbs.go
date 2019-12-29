package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/phodal/coca/core/adapter"
	"github.com/phodal/coca/core/adapter/call"
	"github.com/phodal/coca/core/domain/tbs"
	"github.com/phodal/coca/core/models"
	"github.com/phodal/coca/core/support"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

type TbsCmdConfig struct {
	Path string
}

var (
	tbsCmdConfig TbsCmdConfig
)

var tbsCmd = &cobra.Command{
	Use:   "tbs",
	Short: "generate tests bad smell",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		files := support.GetJavaTestFiles(tbsCmdConfig.Path)
		var identifiers []models.JIdentifier

		identifiers = adapter.LoadTestIdentify(files)
		identifiersMap := adapter.BuildIdentifierMap(identifiers)

		var classes []string = nil
		for _, node := range identifiers {
			classes = append(classes, node.Package+"."+node.ClassName)
		}

		analysisApp := call.NewJavaCallApp()
		classNodes := analysisApp.AnalysisFiles(identifiers, files, classes)

		nodeContent, _ := json.MarshalIndent(classNodes, "", "\t")
		support.WriteToCocaFile("tdeps.json", string(nodeContent))

		app := tbs.NewTbsApp()
		result := app.AnalysisPath(classNodes, identifiersMap)

		fmt.Println("Test Bad Smell nums: ", len(result))
		resultContent, _ := json.MarshalIndent(result, "", "\t")
		support.WriteToCocaFile("tbs.json", string(resultContent))

		if len(result) <= 20  {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Type", "FileName", "Line"})

			for _, result := range result {
				table.Append([]string{result.Type, result.FileName, strconv.Itoa(result.Line)})
			}

			table.Render()
		}
	},
}

func init() {
	rootCmd.AddCommand(tbsCmd)

	tbsCmd.PersistentFlags().StringVarP(&tbsCmdConfig.Path, "path", "p", ".", "example -p core/main")
}
