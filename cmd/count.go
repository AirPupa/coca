package cmd

import (
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"github.com/phodal/coca/cmd/cmd_util"
	"github.com/phodal/coca/cmd/config"
	"github.com/phodal/coca/core/context/count"
	"github.com/phodal/coca/core/domain"
	"github.com/phodal/coca/core/infrastructure/string_helper"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

type CountCmdConfig struct {
	DependencePath string
	Top            int
}

var (
	countCmdConfig CountCmdConfig
)

var cparsedDeps []domain.JClassNode

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "count most refs function",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dependence := countCmdConfig.DependencePath
		if dependence == "" {
			return
		}

		file := cmd_util.ReadFile(dependence)
		if file == nil {
			log.Fatal("lost file:" + dependence)
		}

		_ = json.Unmarshal(file, &cparsedDeps)

		callMap := count.BuildCallMap(cparsedDeps)

		callMapSort := string_helper.RankByWordCount(callMap)

		if countCmdConfig.Top > 0 {
			callMapSort = callMapSort[:countCmdConfig.Top]
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Refs Count", "method"})

		for _, count := range callMapSort {
			table.Append([]string{strconv.Itoa(count.Value), count.Key})
		}

		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(countCmd)

	countCmd.PersistentFlags().StringVarP(&countCmdConfig.DependencePath, "dependence", "d", config.CocaConfig.ReporterPath+"/deps.json", "get dependence file")
	countCmd.PersistentFlags().IntVarP(&countCmdConfig.Top, "top", "t", 0, "top x")
}
