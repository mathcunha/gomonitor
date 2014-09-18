package scheduler

import (
	"bytes"
	"encoding/json"
	"github.com/mathcunha/gomonitor/db"
	"github.com/mathcunha/gomonitor/prop"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Evaluate(m db.Monitor) {
	total, hits, err := getResource(m)
	log.Printf("monitor = (%v) total = (%v) - hits = (%v)", m.Id, total, hits)

	if err == nil && total > 0 {
		var alert db.Alert
		alert.Monitor = m
		alert.Total = total
		alert.When = time.Now()
		PostAlert(alert)
	}
}

func getResource(m db.Monitor) (int, string, error) {
	body, _ := callElasticsearch(strings.Replace(m.Query, "\\", "", -1))
	//log.Printf(body)
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal([]byte(body), &objmap)

	if err != nil {
		log.Printf("error parsing ElasticSearch results - %v", err)
		return -1, "", err
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

func PostAlert(alert db.Alert) {
	var postData []byte
	w := bytes.NewBuffer(postData)
	json.NewEncoder(w).Encode(alert)
	http.Post("http://"+prop.Property("gomonitor")+"/alert", "application/json", w)
}

func callElasticsearch(query string) (string, error) {
	var postData []byte
	w := bytes.NewBuffer(postData)
	w.Write([]byte(query))
	endpoint := "http://" + prop.Property("elasticsearch") + "/" + getIndex() + "/_search"

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
