package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"sigs.k8s.io/yaml"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var conf config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal(err)
	}
	output, err := generate(conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}

type config struct {
	RegoFiles []string `json:"regos"`
	BaseFile  string   `json:"base"`
}
