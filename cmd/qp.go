package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"qp"
)

func main() {
	file := flag.String("file", "", "--file example.qp")
	flag.Parse()

	if *file == "" {
		fmt.Println("require source file to execute")
		os.Exit(1)
		return
	}
	data, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}
	qp.Parse(string(data)).Invoke()
}