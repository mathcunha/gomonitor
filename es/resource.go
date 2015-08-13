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

type Hit struct {
	Total int
	Prop  string
}

type AggsTermItem struct {
	Key   string `bson:"key"        json:"key"`
	Count int    `bson:"doc_count"        json:"doc_count"`
}

func Search(query string) ([]Hit, error) {
	body, _ := callElasticsearch(query)
	//log.Printf(query)
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal([]byte(body), &objmap)
	noHit := []Hit{Hit{-1, ""}}

	if err != nil {
		log.Printf("error parsing ElasticSearch results - %v", err)
		return noHit, err
	}

	if objmap["hits"] == nil {
		log.Printf("no results found, perhaps its a missing indice")
		return noHit, errors.New("no results found")
	}

	if hits := LoadAggsTermsHits(objmap); len(hits) > 0 {
		return hits, nil
	}

	err = json.Unmarshal([]byte(*objmap["hits"]), &objmap)
	if err != nil {
		log.Printf("error parsing ElasticSearch hits - %v", err)
		return noHit, err
	}

	b := []byte(*objmap["total"])
	total, err := strconv.Atoi(string(b[:]))

	if err != nil {
		log.Printf("error parsing ElasticSearch total - %v", err)
		return noHit, err
	}

	b = []byte(*objmap["hits"])

	return []Hit{Hit{total, string(b[:])}}, nil
}

func LoadAggsTermsHits(objmap map[string]*json.RawMessage) []Hit {
	var aggsMap map[string]*json.RawMessage
	if objmap["aggregations"] == nil {
		return nil
	}
	err := json.Unmarshal([]byte(*objmap["aggregations"]), &aggsMap)
	if err != nil {
		return nil
	}
	if aggsMap["terms"] == nil {
		return nil
	}
	err = json.Unmarshal([]byte(*aggsMap["terms"]), &aggsMap)
	if err != nil {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader([]byte(*aggsMap["buckets"])))
	items := []AggsTermItem{}
	decoder.Decode(&items)
	if err != nil {
		log.Printf("error decoding Bucket - %v", err)
		return nil
	}
	hits := make([]Hit, len(items), len(items))
	for i, v := range items {
		hits[i] = Hit{v.Count, v.Key}
	}

	return hits
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
