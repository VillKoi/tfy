package main

import (
	"flag"
	"fmt"
	"os"
	"tfy/lab2/cmd/solution"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config10.txt", "Used for set path to config file.")
	flag.Parse()

	equationSystem, err := solution.GetData(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	ответ := solution.Pешение(equationSystem)

	fmt.Println(ответ.Var + "=" + ответ.Regex)
}
