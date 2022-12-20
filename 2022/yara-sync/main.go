package main

import (
	"bufio"
	"fmt"
	"os"
)

func list_load() []string {
	var repos []string
	readFile, err := os.Open("yara_repo_list.txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		repos = append(repos, fileScanner.Text())
	}

	readFile.Close()
	return repos
}

func main() {

	list_repo := list_load()

	gitClone(list_repo)
	dedup()

}
