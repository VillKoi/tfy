package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config6.txt", "Used for set path to config file.")
	flag.Parse()

	rules, err := GetData(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	example, err := checkConfluence(rules)
	if err != nil {
		log.Fatal(err.Error())
	}

	if example != nil {
		fmt.Printf("Cистема, возможно, не конфлюэнтна, есть перекрытие:\n\t%s\n\t%s\n", rules[example[0]], rules[example[1]])
	} else {
		fmt.Println("Все прошло успешно")
	}
}

type Rule struct {
	First  string
	Second string
}

func GetData(configPath string) ([]Rule, error) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	data := strings.Split(string(file), "\n")

	rules := make([]Rule, 0, len(data))
	var rule Rule
	for i := range data {
		if !strings.Contains(data[i], "->") {
			continue
		}

		str := strings.ReplaceAll(data[i], " ", "")

		terms := strings.Split(str, "->")

		rule = Rule{
			First:  terms[0],
			Second: terms[1],
		}

		rules = append(rules, rule)

		prefix := prefixFunc(rule.First)
		if prefix {
			fmt.Printf("Cистема, возможно, не конфлюэнтна, есть перекрытие внутри терма %s\n", rule.First)
			os.Exit(0)
		}
	}

	return rules, nil
}

func prefixFunc(str string) bool {
	runes := []rune(str)

	for i := 1; i < len(str); i++ {
		suffix := string(runes[0:i])
		if strings.HasSuffix(str, suffix) {
			return true
		}
	}

	return false
}

func checkConfluence(rules []Rule) ([]int, error) {
	for i, ruleI := range rules {
		for j, ruleJ := range rules {
			if i != j && ruleI.First != "" && ruleJ.First != "" {
				concat := ruleI.First + "~" + ruleJ.First

				if prefixFunc(concat) {
					return []int{i, j}, nil
				}
			}
		}
	}

	return nil, nil
}
