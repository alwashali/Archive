package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
)

func gitClone(repos []string) {

	for _, repo := range repos {

		// extract repo name
		repo_splited := strings.Split(repo, "/")
		repo_name := repo_splited[len(repo_splited)-1]

		fmt.Printf("\nClonning %s\n", repo_name)

		_, err := git.PlainClone("repos\\"+repo_name, false, &git.CloneOptions{
			URL:      repo,
			Progress: os.Stdout,
		})

		if err != nil {
			fmt.Println(err)
		}

	}

}
