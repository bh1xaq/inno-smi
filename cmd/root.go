/*
Copyright © 2025 XiaoTan <tanxiao@16iot.cn>
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"inno-smi/libraries/mem"
	"inno-smi/libraries/status"
	"inno-smi/libraries/topo"
	"inno-smi/libraries/version"
	"os"
	"os/exec"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "inno-smi",
	Short: "An application for monitoring Inno Fantasy II-M",
	Long: `An application for monitoring Inno Fantasy II-M
It can be used to check driver versions, GPU usage, etc...
meminfo: 获取显存分配情况
memusage: 获取显存使用情况
gputopo: 获取 GPU 拓扑
version: 获取版本信息
默认获取 status 信息
`,
	Run: func(cmd *cobra.Command, args []string) {
		// 测试显卡类型
		hased, err := hasInnoFantasyII()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !hased {
			fmt.Println("未识别到 Inno Fantasy II-M 显卡，不支持运行 SMI 工具")
			os.Exit(1)
		}
		if len(args) < 1 {
			fmt.Println(status.GetMainInnoStatusToTable())
			return
		}
		switch args[0] {
		case "meminfo":
			fmt.Println(mem.GetInfoToTable())
		case "memusage":
			fmt.Println(mem.GetUsageToTable())
		case "gputopo":
			fmt.Println(topo.GetGPUTopoToTable())
		case "version":
			fmt.Println(version.GetVersionToTable())
		default:
			fmt.Println(status.GetMainInnoStatusToTable())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.inno-smi.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func hasInnoFantasyII() (bool, error) {
	cmd := exec.Command("/usr/bin/lspci")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Innosilicon Co Ltd Fantasy II-M") {
			return true, nil
		}
	}
	return false, nil
}
