package main

import (
	"flag"
	"fmt"
	"gitlab.com/akzj/qp/parser"
	"io/ioutil"
	"os"
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
	parser.New(string(data)).Parse().Invoke()
}
