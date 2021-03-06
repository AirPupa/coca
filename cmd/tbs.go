package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/phodal/coca/cmd/cmd_util"
	"github.com/phodal/coca/core/adapter/coca_file"
	"github.com/phodal/coca/core/context/analysis"
	"github.com/phodal/coca/core/context/tbs"
	"github.com/phodal/coca/core/domain"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

type TbsCmdConfig struct {
	Path   string
	IsSort bool
}

var (
	tbsCmdConfig TbsCmdConfig
)

var tbsCmd = &cobra.Command{
	Use:   "tbs",
	Short: "generate tests bad smell",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		files := coca_file.GetJavaTestFiles(tbsCmdConfig.Path)
		var identifiers []domain.JIdentifier

		identifiers = cmd_util.LoadTestIdentify(files)
		identifiersMap := domain.BuildIdentifierMap(identifiers)

		var classes []string = nil
		for _, node := range identifiers {
			classes = append(classes, node.Package+"."+node.ClassName)
		}

		analysisApp := analysis.NewJavaFullApp()
		classNodes := analysisApp.AnalysisFiles(identifiers, files, classes)

		nodeContent, _ := json.MarshalIndent(classNodes, "", "\t")
		cmd_util.WriteToCocaFile("tdeps.json", string(nodeContent))

		app := tbs.NewTbsApp()
		result := app.AnalysisPath(classNodes, identifiersMap)

		fmt.Println("Test Bad Smell nums: ", len(result))
		resultContent, _ := json.MarshalIndent(result, "", "\t")

		if tbsCmdConfig.IsSort {
			var tbsMap = make(map[string][]tbs.TestBadSmell)
			for _, tbs := range result {
				tbsMap[tbs.Type] = append(tbsMap[tbs.Type], tbs)
			}

			resultContent, _ = json.MarshalIndent(tbsMap, "", "\t")
		}

		cmd_util.WriteToCocaFile("tbs.json", string(resultContent))

		if len(result) <= 20 {
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
	tbsCmd.PersistentFlags().BoolVarP(&tbsCmdConfig.IsSort, "sort", "s", false, "-s")
}
