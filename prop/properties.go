package prop

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var properties map[string]*json.RawMessage = loadConfig()

func loadConfig() map[string]*json.RawMessage {
	body, err := ioutil.ReadFile("prop/config.json")

	if err != nil {
		log.Printf("error %v", err)
	}

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(body, &objmap)
	if err != nil {
		objmap = make(map[string]*json.RawMessage)

		objmap["elasticsearch"] = new(json.RawMessage)
		objmap["smtp"] = new(json.RawMessage)
		objmap["gomonitor"] = new(json.RawMessage)
		objmap["mongodb"] = new(json.RawMessage)

		*objmap["elasticsearch"] = json.RawMessage([]byte("127.0.0.1:9200"))
		*objmap["smtp"] = json.RawMessage([]byte("127.0.0.1:25"))
		*objmap["gomonitor"] = json.RawMessage([]byte("127.0.0.1:8080"))
		*objmap["mongodb"] = json.RawMessage([]byte("127.0.0.1"))
	}

	return objmap
}

func Property(key string) string {
	return string([]byte(*properties[key])[:])
}
