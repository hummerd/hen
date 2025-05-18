package main

import (
	"context"
	"fmt"

	"github.com/hummerd/hen/samples/github"
)

func main() {
	c := github.NewClient()

	branches, err := c.GetBranches(context.Background(), "hummerd", "hen")
	if err != nil {
		panic(err)
	}

	fmt.Println(branches)
}
