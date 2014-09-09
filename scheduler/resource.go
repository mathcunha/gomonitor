package scheduler

import (
	"bytes"
	"encoding/json"
	"github.com/mathcunha/gomonitor/db"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Evaluate(m db.Monitor) {
	total, hits := getResource(m)
	log.Printf("monitor = (%v) total = (%v) - hits = (%v)", m.Id, total, hits)

	if total > 0 {
		var alert db.Alert
		alert.Monitor = m
		alert.Total = total
		PostAlert(alert)
	}
}

func getResource(m db.Monitor) (int, string) {
	body := callElasticsearch(strings.Replace(m.Query, "\\", "", -1))
	//log.Printf(body)
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal([]byte(body), &objmap)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	err = json.Unmarshal([]byte(*objmap["hits"]), &objmap)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	b := []byte(*objmap["total"])
	total, err := strconv.Atoi(string(b[:]))

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	b = []byte(*objmap["hits"])

	return total, string(b[:])
}

func PostAlert(alert db.Alert) {
	var postData []byte
	w := bytes.NewBuffer(postData)
	json.NewEncoder(w).Encode(alert)
	http.Post("http://127.0.0.1:8080/alert", "application/json", w)
}

func callElasticsearch(query string) string {
	var postData []byte
	w := bytes.NewBuffer(postData)
	w.Write([]byte(query))

	res, err := http.Post("http://127.0.0.1:9200/"+getIndex()+"/_search", "application/json", w)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	robots, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(robots[:])
}

func getIndex() string {
	t := time.Now()
	y := t.AddDate(0, 0, -1)

	return "logstash-" + t.Format("2006.01.02") + ",logstash-" + y.Format("2006.01.02")
}
