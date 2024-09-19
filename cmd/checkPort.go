/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// checkPortCmd represents the checkPort command
var checkPortCmd = &cobra.Command{
	Use:     "checkPort",
	Aliases: []string{"cp"},
	Short:   "检查端口是否被占用",
	Long:    ``,
	Run:     checkPort,
}

func init() {
	rootCmd.AddCommand(checkPortCmd)
}

func checkPort(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println(`Usage: tool cp port`)
		return
	}

	port := args[0] // 从命令行参数中获取端口号
	check(port)
}

func check(port string) {
	// 使用传入的端口号构造命令
	cmd := exec.Command("sudo", "lsof", "-i", ":"+port)
	// 注意：在某些系统上，lsof 命令可能需要其他参数或可能不可用。
	// 您可能需要使用 netstat、ss 或其他命令来检查端口。

	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut

	err := cmd.Run()
	if err != nil {
		return
	}

	// 输出命令的执行结果
	fmt.Println(stdOut.String())

}
