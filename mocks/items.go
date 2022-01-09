package mocks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func GetRecord(key string) string {
	// FUTURE: Maybe log instead of print, but this is just for tests so it doesnt really matter.
	jsonFile, err := os.Open("mocks/data.json")
	if err != nil {
		fmt.Println(err)
	}

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var record = make(map[string]interface{})
	err = json.Unmarshal(data, &record)
	if err != nil {
		fmt.Println(err)
	}

	bytes, err := json.Marshal(record[key])
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)
}
