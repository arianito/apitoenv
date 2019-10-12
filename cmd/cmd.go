package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func Init() {
	address := flag.String("api", "", "api address")
	path := flag.String("path", "", "json path")
	key := flag.String("key", "", "key")
	flag.Parse()

	resp, err := http.Get(*address)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	value, err := getDeep(data, strings.Split(*path, "."))
	if key != nil && *key != "" {
		err = os.Setenv(*key, value.(string))
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Print(value)
}
func getDeep(rest map[string]interface{}, keys []string) (interface{}, error) {
	if len(keys) < 1 {
		return nil, errors.New("not found")
	}
	key := keys[0]

	value, ok := rest[key]
	if !ok {
		return nil, errors.New(key + " not found")
	}
	nw, isMap := value.(map[string]interface{})
	if !isMap {
		return value, nil
	}
	return getDeep(nw, keys[1:])
}
