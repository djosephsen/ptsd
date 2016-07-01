package main

import (
	"net/http"
	"fmt"
	"os"
	"time"
	"strings"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"github.com/jmoiron/jsonq"
	statsd "github.com/etsy/statsd/examples/go"
)

var metrics = statsd.New(`localhost`, 8125)

func increment(key string, value int){
	fmt.Printf("%s:: +%d\n",key,value)
	metrics.IncrementByValue(key, value)
}

func pdQuery(url string) []byte{
	client := &http.Client{}
	authToken := fmt.Sprintf("Token token=%s",os.Getenv("PDTOKEN"))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Add("Authorization", authToken)
	resp,_ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func hrAgo() string{
	offset,_ := strconv.Atoi(os.Getenv("PDHOURS"))
	return time.Now().UTC().Add(time.Duration(offset)*time.Hour).Format("2006-01-02T15:04:05Z")
}

func countLogs(id string){
	l := map[string]interface{}{}
	lURL := fmt.Sprintf("https://api.pagerduty.com/incidents/%s/log_entries?time_zone=UTC&is_overview=false",id)
	dec := json.NewDecoder(strings.NewReader(string(pdQuery(lURL))))
	dec.Decode(&l)
	jq := jsonq.NewQuery(l)
	log_array, _ := jq.ArrayOfObjects("log_entries")
	for _,log := range log_array{
		q := jsonq.NewQuery(log)
		entry_type, _ := q.String("type")
		if	entry_type == `notify_log_entry`{
			user,_ := q.String("user","summary")
			key := fmt.Sprintf("pagerduty.%s.notified",user)
			increment(key,1)
		}else if	entry_type == `assign_log_entry`{
			assignee_array,_ := q.ArrayOfObjects("assignees")
			for _,assignee := range assignee_array {
				j := jsonq.NewQuery(assignee)
				user,_ := j.String("summary")
				key := fmt.Sprintf("pagerduty.%s.assigned",user)
				increment(key,1)
			}
		}
	}
}

func main(){
	inc := map[string]interface{}{}
	iURL := fmt.Sprintf("https://api.pagerduty.com/incidents?time_zone=UTC&since=%s",hrAgo())
	dec := json.NewDecoder(strings.NewReader(string(pdQuery(iURL))))
	dec.Decode(&inc)
	jq := jsonq.NewQuery(inc)

	incident_array,_ := jq.ArrayOfObjects("incidents")
	for _,incident := range incident_array{
		increment("pagerduty.incidents",1)
		q := jsonq.NewQuery(incident)
		id,_ := q.String("id")
		countLogs(id)
	}
}
