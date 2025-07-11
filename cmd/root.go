package cmd

import (
	"os"

	"moeCounter/cmd/flags"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "moeCounter",
	Short: "萌萌计数器！",
	Long:  "萌萌计数器是一个用于图片拼接的计数器服务，支持多种主题和配置选项。",
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
	// 在这里添加子命令或标志
	rootCmd.PersistentFlags().IntVarP(&flags.Port, "port", "p", 8088, "服务监听端口")
	rootCmd.PersistentFlags().StringVarP(&flags.DbFile, "db", "d", "data.db", "数据库文件路径")
	rootCmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "是否开启调试模式")
}
