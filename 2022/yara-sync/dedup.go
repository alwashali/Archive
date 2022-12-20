package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type rule struct {
	rulename string
	sha1     string
	path     string
}

var rules []rule

func sha1hash(filepath string) (string, error) {
	var returnSHA1String string

	//Open the filepath passed by the argument and check for any error
	file, err := os.Open(filepath)
	if err != nil {
		return returnSHA1String, err
	}

	defer file.Close()

	hash := sha1.New()

	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA1String, err
	}

	hashInBytes := hash.Sum(nil)[:20]

	returnSHA1String = hex.EncodeToString(hashInBytes)

	return returnSHA1String, nil

}

func makeList() {
	root := "repos"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {

			fmt.Println(err)
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".yar" {
			rulesha1, err := sha1hash(path)

			if err != nil {
				fmt.Println(err)
			}
			abs_path, err := filepath.Abs(path)
			filename := filepath.Base(path)
			r := rule{
				rulename: filename,
				path:     abs_path,
				sha1:     rulesha1,
			}
			rules = append(rules, r)

		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

}

func dedup() {
	makeList()

	var uniq_rules []rule
	count := 0
	for _, item := range rules {
		if contains(uniq_rules, item) == false {
			uniq_rules = append(uniq_rules, item)
		} else {
			count++
			fmt.Println(item.rulename)
		}
	}
	fmt.Printf("Found %d duplicate rules\n", count)

	abs_dst, err := filepath.Abs("Rules")
	if err != nil {
		fmt.Println(err)
	}

	for _, rule := range uniq_rules {

		copy(rule.path, abs_dst)
	}
	fmt.Printf("Number of rules: %d\n", len(uniq_rules))

}

func contains(rules []rule, rule rule) bool {
	for _, r := range rules {
		if r.sha1 == rule.sha1 {
			fmt.Println(r.rulename)
			return true

		}
	}

	return false
}

func copy(src string, dst string) {

	data, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(err)
	}
	filename := filepath.Base(src)

	err = ioutil.WriteFile(filepath.Join(dst, filename), data, 0644)
	if err != nil {
		fmt.Println(err)
	}

}
