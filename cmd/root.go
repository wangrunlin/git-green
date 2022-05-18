/*
Copyright © 2022 Leo Wang leo@wangrunlin.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var cfgFile string
var year int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "green",
	Short: "A brief description of your application",
	Long:  `green long description`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// verify default data

		// default data
		fmt.Print("默认生成最近一年（365天）的数据, 是否确认？[Y/n] ")
		var yes string
		fmt.Scanln(&yes)
		if yes == "y" || strings.ToLower(yes) == "y" || yes == "" {
			// 具体业务
			// 0. 确保处于一个干净的具体上传权限的 git 项目
			// 0.1 确保存在 git 命令

			branch := "git-green"
			// 1. checkout
			s, _ := os.Stat(".git/refs/heads/" + branch)
			if s != nil {
				fmt.Println(branch + " 分支已存在")
				exec.Command("git", "checkout", branch)
			} else {

				fmt.Println(branch + " 分支不存在, 新建分支")
				exec.Command("git", "checkout", "-b", branch)
			}
			// 2. for each
			filename := "git-green.md"

			// 2.1 open file
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("文件打开失败: {}", err)
			}
			defer f.Close()

			count := 10
			date := time.Now()
			var content string
			for i := 0; i < count; i++ {
				date = time.Now().AddDate(0, 0, -i).Local()

				// write file
				// 2.2 write file
				content = date.String() + ": update by [Git Green](https://github.com/wangrunlin/git-green)\n"
				_, err := f.WriteString(content)
				if err != nil {
					fmt.Println("不能写入文件: {}", err)
					return
				}

				// 3. modify date
				message := "Update git-green.md"
				exec.Command("git", "add", filename)
				exec.Command("git", "commit", "--amend", "--date=" + date.Format("2006-01-02"), "-m", message)
			}

			// 4. push
			fmt.Print("是否上传到远程客户端？[Y/n] ")
			var isPush string
			fmt.Scanln(&isPush)
			if isPush == "" || strings.ToLower(isPush) == "y" {
				exec.Command("git", "push", "origin", "git-green")
			}
		} else {
			fmt.Println("Bye.")
			return
		}


		openFlag := cmd.Flag("open").Value.String() == "true"
		if openFlag {
			// TODO: 打开主页
			open("https://wangrunlin.com")
			fmt.Println("Open home page")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	//rootCmd.PersistentFlags().StringVarP(&year, "year", "y", "2022", "generate year")
	rootCmd.Flags().IntVarP(&year, "year", "y", 2022, "generate year")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("open", "o", false, "Open home page")
}

var commands = map[string] string{
	"windows": "cmd /c start ",
	"darwin": "open",
	"linux": "xdg-open",
}

func open(url string) error {
	command, ok := commands[runtime.GOOS]
	if !ok {
		fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	return exec.Command(command, url).Start()
}
