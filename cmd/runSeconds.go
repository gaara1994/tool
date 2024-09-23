/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/cobra"
)

// runSecondsCmd represents the runSeconds command
var runSecondsCmd = &cobra.Command{
	Use:     "runSeconds",
	Aliases: []string{"rs"},
	Short:   "运行指定秒数",
	Long:    ``,
	Run:     runSeconds,
}

func init() {
	rootCmd.AddCommand(runSecondsCmd)
}

func runSeconds(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: tool rs [seconds] [succeed|failed|exit|loop|cpu|memory|oom]")
		return
	}

	seconds, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("runSeconds seconds error:", err)
		return
	}

	mode := args[1]
	fmt.Println("runSeconds", seconds, "秒", "模式为：", mode)

	switch mode {
	case "succeed":
		sleep(seconds)
		fmt.Println("runSeconds succeed")
		return

	case "failed":
		sleep(seconds)
		panic("runSeconds failed")

	case "exit":
		fmt.Println("runSeconds failed")
		os.Exit(1)
		return

	case "loop":
		for {
			time.Sleep(time.Hour * 24)
		}

	case "cpu":
		var wg sync.WaitGroup                                   // 定义一个sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background()) // 创建一个context

		// 获取当前系统可用的逻辑CPU数量
		numCPU := runtime.NumCPU()
		runtime.GOMAXPROCS(numCPU) // 设置Go调度器使用的最大线程数
		// 启动与逻辑CPU数量相同的goroutine
		for i := 0; i < numCPU; i++ {
			go busyWork(ctx, &wg)
		}

		sleep(seconds)
		cancel()
		wg.Wait()

	case "memory":
		memory()
		// oom()
		sleep(seconds)
		return
	case "oom":
		oom()

	default:
		fmt.Println("未定义的参数:", mode)
		return
	}
}

// busyWork 模拟一个高负荷工作的函数
// 该函数没有输入参数和返回值
// 它的设计目的是模拟一个持续进行计算任务的工作流
func busyWork(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		fmt.Println("busyWork 退出")
		return
	default:
		for {

		}
	}
}

// sleep 函数使程序暂停指定的秒数。
// 在暂停期间，每过一秒，函数会打印当前运行的秒数。
// 参数 seconds 指定了程序暂停的总秒数。
func sleep(seconds int) {
	// 每运行一秒，打印正在运行的秒数
	for i := 0; i < seconds; i++ {
		fmt.Println("正在运行", i+1, "秒")
		time.Sleep(time.Second)
	}
}

// memory 函数
// 使内存达到90%
func memory() {
	// 获取内存信息
	v, _ := mem.VirtualMemory()
	data := make([]byte, v.Free)
	for i, _ := range data {
		data[i] = 1
	}
}

// oom函数
// 触发OOM
func oom() {
	// 获取内存信息
	v, _ := mem.VirtualMemory()
	data := make([]byte, v.Total)
	for i, _ := range data {
		data[i] = 1
	}
}
