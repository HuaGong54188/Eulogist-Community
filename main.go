package main

import (
	Eulogist "Eulogist/eulogist"
	"fmt"

	"github.com/pterm/pterm"
)

func main() {
	err := Eulogist.Eulogist()
	if err != nil {
		pterm.Error.Println(err)
	}

	fmt.Println()
	pterm.Info.Println("程序正在运行，回车即可退出。")
	fmt.Scanln()
}
