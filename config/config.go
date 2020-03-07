package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type configStruct struct {
	Token string `json:"Token"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// ReadConfig copied from previous project and simplified.
// panic if we don't get a Token
func ReadConfig() string {
	fmt.Println("Reading from config file...")
	file, err := ioutil.ReadFile("./config.json")
	checkError(err)

	config := configStruct{}
	err = json.Unmarshal(file, &config)
	checkError(err)

	fmt.Println("Got Token: " + config.Token)
	return config.Token
}
