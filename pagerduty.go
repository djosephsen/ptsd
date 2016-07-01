package main

import (
	"net/http"
	"fmt"
	"os"
	"time"
	"strconv"
	jq "github.com/antonholmquist/jason"
)

type PDIncident struct{
	Incident 	*jq.Object
	Logs			[]*jq.Object
}

func hrAgo() string{
	offset,_ := strconv.Atoi(os.Getenv("PDHOURS"))
	return time.Now().UTC().Add(time.Duration(offset)*time.Hour).Format("2006-01-02T15:04:05Z")
}

func pdQuery(url string, token string) *jq.Object {
	client := &http.Client{}
	authToken := fmt.Sprintf("Token token=%s",token)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Add("Authorization", authToken)
	resp,err := client.Do(req)
	if err != nil{
		debug(fmt.Sprintf("error was %s",err))
	}
	defer resp.Body.Close()
	body,_ := jq.NewObjectFromReader(resp.Body)
	return body
}

func pdGetLogs(id string, token string) []*jq.Object {
	ret := []*jq.Object{}
	lURL := fmt.Sprintf("https://api.pagerduty.com/incidents/%s/log_entries?time_zone=UTC&is_overview=false",id)
	raw := pdQuery(lURL,token)
	logs, _ := raw.GetObjectArray("log_entries")
	for _,log := range logs{
		ret = append(ret, log)
	}
	return ret
}

func pdProcessLog(log *jq.Object) error { 
	entry_type, _ := log.GetString("type")
	if	entry_type == `notify_log_entry`{
		user,_ := log.GetString("user","summary")
		key := fmt.Sprintf("pagerduty.%s.notified",user)
		increment(key,1)
	}else if	entry_type == `assign_log_entry`{
		assignees,_ := log.GetObjectArray("assignees")
		for _,assignee := range assignees {
			user,_ := assignee.GetString("summary")
			key := fmt.Sprintf("pagerduty.%s.assigned",user)
			increment(key,1)
		}
	}
	//todo proper error checking
	return nil
}

func pdGetIncidents(token string) []*PDIncident{
	ret := []*PDIncident{}
	iURL := fmt.Sprintf("https://api.pagerduty.com/incidents?time_zone=UTC&since=%s",hrAgo())
	debug(iURL)
	raw := pdQuery(iURL,token)
	debug(raw.String())
	incidents,_ := raw.GetObjectArray("incidents")
	for _,incident := range incidents{
		id,_ := incident.GetString("id")
		me := &PDIncident{
			Incident : incident,
			Logs : pdGetLogs(id,token),
		}
		ret = append(ret, me)
	}
	return ret
}

func runPagerDuty(token string) {
	for _,incident := range pdGetIncidents(token){
		increment("pagerduty.incidents",1)
		for _,log := range incident.Logs{
			pdProcessLog(log)
		}
	}
}
