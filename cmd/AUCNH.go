/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	OZT = 31.1034768
)

// AUCNHCmd represents the AUCNH command
var AUCNHCmd = &cobra.Command{
	Use:     "AUCNH",
	Aliases: []string{"auc"},
	Short:   "XAUUSD->AUCNH",
	Long: `
现货黄金->人民币黄金

使用方法:
  tool ac [XAUUSD] [USDCNH]
  tool ac [XAUUSD] [USDCNH]  # 使用别名

其中：
  `,
	Run: AUCNH,
}

func init() {
	rootCmd.AddCommand(AUCNHCmd)
}

func AUCNH(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: tool ac [XAUUSD] [USDCNH]")
		return
	}

	XAUUSD, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		fmt.Println("XAUUSD err: ", err)
		return
	}

	USDCNH, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Println("USDCNH err: ", err)
		return
	}

	AUCNH := XAUUSD / OZT * USDCNH

	fmt.Println("此时黄金的价格应为: ", AUCNH, "人民币/克")
}
