package scheduler

import (
	"bytes"
	"github.com/mathcunha/gomonitor/db"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

func getIndex() string{
	t := time.Now()
        y := t.AddDate(0,0,-1)

        return t.Format("2006.01.02")+","+y.Format("2006.01.02")
}
