package es

import (
	"github.com/mathcunha/gomonitor/prop"
	"testing"
)

func TestSearch(t *testing.T) {
	prop.LoadConfig("../config/config.json")
	//hits, _ := Search("{	\"query\": 		{\"bool\": 			{\"must\": 				[{\"match_phrase\" : {\"@fields.State\" : {\"query\" : \"ERR\"}}},{\"range\" : {\"@timestamp\" : {\"gte\" : \"now-360m\"}}}]			}		}	,	\"aggregations\" : {		\"terms\" : {			\"terms\" : {				\"field\" : \"@fields.instance.raw\"			}		}	}}")
	hits, _ := Search("{\"query\":{\"bool\":{\"must\":[{\"match_phrase\":{\"@fields.State\":{\"query\":\"ERR\"}}},{\"range\":{\"@timestamp\":{\"gte\":\"now-360m\"}}}]}},\"aggregations\":{\"terms\":{\"terms\":{\"field\":\"@fields.instance.raw\"}}}}")
	t.Log(hits)
	hits, _ = Search("{\"query\": {\"bool\": {\"must\": [{\"match_phrase\" : {\"@fields.State\" : {\"query\" : \"ERR\"}}},{\"range\" : {\"@timestamp\" : {\"gte\" : \"now-360m\"}}}]}}}")
	t.Log(hits)
}
