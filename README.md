gomonitor
=========

What if your logstash metrics stored at Elasticsearch were monitored? Yeah...

This projects aims to monitor the metrics colleted by logstash and take actions when something happens. Those actions could be send a email, or run a script, which are native to the project. But you can develop your own actions.

api
=========
To achieve that goal, Gomonitor exposes an REST API, this API has two major resources, Monitor, Actions (e.g Sendmail, etc.). The first represents what is been monitored (e.g. apache response time, JVM memory) and the second represents the action that should be taken when a monitor alerts.

setup
=========
After clone the project, you should edit de <i>prop/config.json</i> file, providing some basic informations, like elasticsearch and gomonitor url.
```json
{
"elasticsearch" : "127.0.0.1:9200",
"smtp" : "127.0.0.1:25",
"gomonitor" : "127.0.0.1:8080",
"mongodb" : "127.0.0.1"
}
```

example
=========
To start monitor a metric, you should insert a monitor, as the command below.

```javascript
curl -XPOST http://localhost:8080/monitor -d '{"query":"{\"query\" : {\"bool\" : {\"must\": [{\"match\" : {\"@fields.request\" : {\"query\" : \"/MYAPP\", \"type\":\"phrase_prefix\"}}},{\"range\" : {\"@timestamp\" : {\"gte\" : \"now-15m\"}}},{\"range\" : {\"@fields.reqmsecs\" : {\"gte\" : 175328151}}}]}}}","interval":"15m","actions":["sendmail"]}'
```

The property <b>query</b> represts a Elasticseach query, the property <b>interval</b> the time between each monitor call, and the property <b>actions</b> a array of actions.

Once a monitor is inserted, the Gomonitor will schedule a Goroutine tha calls the query and alert if the itens returned are bigger than zero.



