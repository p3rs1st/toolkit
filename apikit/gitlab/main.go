package main

import (
	"os"
	"toolkit/apikit/gitlab/cmd"
)

func main() {
	err := cmd.NewGitlabCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
