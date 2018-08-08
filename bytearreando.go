package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	var file string
	flag.StringVar(&file, "f", "file", "")
	flag.Parse()

	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("%v", content)
}
