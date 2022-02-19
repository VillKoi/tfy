package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	re          = regexp.MustCompile(`([2468]|[13579][24680]*[13579])([24680]|[13579][24680]*[13579])*`)
	re_negation = regexp.MustCompile(`([^135790]|[^24680][^13579]*[^24680])([^13579]|[^24680][^13579]*[^24680])*`)
	re_lazy     = regexp.MustCompile(`^([2468]|[13579][24680]*?[13579])([24680]|[13579][24680]*?[13579])*?$`)
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.txt", "Used for set path to config file.")
	flag.Parse()

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	data := strings.Split(string(file), "\n")
	for i := range data {
		if data[i] == "" {
			continue
		}

		time_1 := time.Now()
		ss := re.FindString(data[i])
		duration_1 := time.Now().Sub(time_1)

		if !strings.EqualFold(ss, data[i]) {
			fmt.Println(i, "0", "строка не соответствует регулярному выражению")
			continue
		}

		time_2 := time.Now()
		_ = re_negation.FindString(data[i])
		duration_2 := time.Now().Sub(time_2)

		time_3 := time.Now()
		_ = re_lazy.FindString(data[i])
		duration_3 := time.Now().Sub(time_3)

		fmt.Printf("%d 1 %-8s 2 %-8s 3 %-8s\n", i, duration_1, duration_2, duration_3)
	}
}
