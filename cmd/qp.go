package main

import (
	"flag"
	"fmt"
	"github.com/akzj/qp"
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
	qp.NewParse2(string(data)).Parse().Invoke()
}
