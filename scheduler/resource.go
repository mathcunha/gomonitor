package scheduler

import (
	"bytes"
	"github.com/mathcunha/gomonitor/db"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Evaluate(m db.Monitor) {
	log.Printf(callElasticsearch(m.Query))
}

func callElasticsearch(query string) string {
	var postData []byte
	w := bytes.NewBuffer(postData)
	body := strings.Replace(query, "\\", "", -1)
	w.Write([]byte(body))

	res, err := http.Post("http://127.0.0.1:9200/logstash-2014.09.08/_search", "application/json", w)

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
