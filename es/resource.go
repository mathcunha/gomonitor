package es

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/mathcunha/gomonitor/prop"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Search(query string) (int, string, error) {
	body, _ := callElasticsearch(query)
	//log.Printf(body)
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal([]byte(body), &objmap)

	if err != nil {
		log.Printf("error parsing ElasticSearch results - %v", err)
		return -1, "", err
	}

	if objmap["hits"] == nil {
		log.Printf("no results found, perhaps its a missing indice")
		return -1, "", errors.New("no results found")
	}

	err = json.Unmarshal([]byte(*objmap["hits"]), &objmap)
	if err != nil {
		log.Printf("error parsing ElasticSearch hits - %v", err)
		return -1, "", err
	}

	b := []byte(*objmap["total"])
	total, err := strconv.Atoi(string(b[:]))

	if err != nil {
		log.Printf("error parsing ElasticSearch total - %v", err)
		return -1, "", err
	}

	b = []byte(*objmap["hits"])

	return total, string(b[:]), nil
}

func callElasticsearch(query string) (string, error) {
	var postData []byte
	w := bytes.NewBuffer(postData)
	w.Write([]byte(query))
	endpoint := os.Getenv(prop.Property("elasticsearch"))
	endpoint = strings.Replace(endpoint, "tcp", "http", 1)
	endpoint = endpoint + "/" + getIndex() + "/_search"

	res, err := http.Post(endpoint, "application/json", w)

	if err != nil {
		log.Printf("error calling ElasticSearch at %v  [%v]", endpoint, err)
		return "", err
	}
	defer res.Body.Close()

	robots, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Printf("error reading the ElasticSearch results at %v  [%v]", endpoint, err)
		return "", err
	}
	return string(robots[:]), nil
}

func getIndex() string {
	t := time.Now()
	y := t.AddDate(0, 0, -1)

	return "logstash-" + t.Format("2006.01.02") + ",logstash-" + y.Format("2006.01.02")
}
