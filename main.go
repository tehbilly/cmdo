package main

import (
	"fmt"
	"os"

	"github.com/tehbilly/cmdo/commands"
)

func main() {
	// TODO Make this the functionality of the root command
	//if len(os.Args) == 2 && os.Args[1] == "link" {
	//	createLinks(os.Args[2])
	//	return
	//}

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

//func createLinks(targetDir string) {
//	// If target dir doesn't exist, create it
//	if _, err := os.Stat(targetDir); err != nil && os.IsNotExist(err) {
//		if err := os.MkdirAll(targetDir, 0755); err != nil {
//			fmt.Println("Error creating target directory:", err.Error())
//			os.Exit(1)
//		}
//	}
//
//	target := filepath.Clean(os.Args[0])
//
//	switch runtime.GOOS {
//	case "windows":
//		cmdName := strings.TrimRight(filepath.Base(target), filepath.Ext(target))
//		shimPath := filepath.Join(targetDir, cmdName+".cmd")
//		shimCmd(shimPath, target)
//		fmt.Println("Shim .cmd files created in:", targetDir)
//	default:
//		if err := os.Symlink(os.Args[0], filepath.Join(targetDir, filepath.Base(target))); err != nil {
//			fmt.Println("Unable to create symolic link:", err.Error())
//			os.Exit(1)
//		}
//		fmt.Println("Links should be created in:", targetDir)
//	}
//}
