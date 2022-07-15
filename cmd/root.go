package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
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
		fileInfo, err := os.Stat(".git")
		if err != nil || !fileInfo.IsDir() {
			fmt.Println("不是 git 文件夹根目录！")
			fmt.Errorf("%s", err)
			return
		}

		// default data
		fmt.Print("默认生成最近一年（365天）的数据, 是否确认？[Y/n] ")
		var yes string
		fmt.Scanln(&yes)
		if yes == "y" || strings.ToLower(yes) == "y" || yes == "" {
			// 具体业务
			// 0. 确保处于一个干净的具体上传权限的 git 项目
			// 0.1 确保存在 git 命令

			branch := "gh-pages"
			// 1. checkout branch
			s, _ := os.Stat(".git/refs/heads/" + branch)
			if s != nil {
				fmt.Println(branch + " 分支已存在，切换分支")
				RunCommand("git", "checkout", branch)
			} else {
				fmt.Println(branch + " 分支不存在，新建分支")
				RunCommand("git", "checkout", "-b", branch)
			}
			// 2. for each
			filename := "git-green.md"

			// 2.1 open file
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("文件打开失败: {}", err)
			}
			defer f.Close()

			count := 365
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
				RunCommand("git", "add", filename)
				RunCommand("git", "commit", "-m", message)
				RunCommand("git", "commit", "--amend", "--no-edit", "--date=" + date.Format("2006-01-02"))
			}

			// 4. push
			// isPush config
			isPush := cmd.Flag("push").Value.String() == "true"
			if isPush {
				RunCommand("git", "push", "origin", "gh-pages")
			}
		} else {
			fmt.Println("Bye.")
			return
		}

		openFlag := cmd.Flag("open").Value.String() == "true"
		if openFlag {
			// TODO: 打开主页
			Open("https://wangrunlin.com")
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
	//rootCmd.Flags().IntVarP(&year, "year", "y", 2022, "generate year")
	rootCmd.Flags().BoolP("open", "o", false, "Open home page")
	rootCmd.Flags().BoolP("push", "p", false, "Push remote repository")
}

var commands = map[string] string{
	"windows": "cmd /c start ",
	"darwin": "open",
	"linux": "xdg-open",
}

func Open(url string) error {
	command, ok := commands[runtime.GOOS]
	if !ok {
		fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	return exec.Command(command, url).Start()
}

func RunCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	//fmt.Printf("out: %s\n", out)
}
