package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/gookit/color"
)

type TestLine struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Output  string    `json:"Output,omitempty"`
	Elapsed float32   `json:"Elapsed,omitempty"`
}

func main() {

	// read from stdin line by line
	fileScanner := bufio.NewScanner(os.Stdin)

	red := color.FgRed.Render
	green := color.FgGreen.Render
	gray := color.FgGray.Render

	result := map[string][]TestLine{}
	for fileScanner.Scan() {
		line := TestLine{}
		jsonline := fileScanner.Text()
		//fmt.Println(string(jsonline))
		err := json.Unmarshal([]byte(jsonline), &line)
		if err != nil {
			panic(err)
		}
		if line.Test != "" {
			result[line.Test] = append(result[line.Test], line)
		}
	}
	keys := make([]string, 0, len(result))
	for key := range result {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return findState(result[keys[i]]) > findState(result[keys[j]])
	})

	for _, v := range keys {
		test := result[v]
		showOutput := false
		for _, line := range test {
			if line.Action == "pass" {
				fmt.Println(green("✔"), green(line.Test))
			}
			if line.Action == "fail" {
				fmt.Println(red("✖"), red(line.Test))
				showOutput = true
			}
			if line.Action == "skip" {
				fmt.Println(line.Action, gray("-"), gray(line.Test))
			}
			if showOutput {
				if line.Action == "output" {
					fmt.Printf("\t%v", line.Output)
				}
			}
		}
	}
}

func findState(lines []TestLine) string {
	for _, line := range lines {
		if line.Action == "fail" {
			return "fail"
		}
		if line.Action == "pass" {
			return "pass"
		}
	}
	return "skip"
}
