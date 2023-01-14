package conf

import (
	"encoding/json"
	"io/ioutil"
)

func LoadFromJsonFile(path string) *Root {
	var root = new(Root)
	bts, err := ioutil.ReadFile(path)
	throwError(err, json.Unmarshal(bts, root))
	return root.prepare()
}

func throwError(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
