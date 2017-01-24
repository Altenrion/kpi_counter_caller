package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
)

func main() {
	url := "http://localhost:8080/chart"
	fmt.Println("URL:>", url)

	request := gorequest.New()
	requestsList := []string{"request1", "request2", "request3"}


	for _, request_name := range requestsList {

		charts := make(map[string][]Chart)
		charts["wfm"] = []Chart{
			{"T_READY_PLAN", "line"}, {"PCT_READY", "line"},
		}
		charts["pds"] = []Chart{
			{"JOB_TALK", "line"}, {"JOB_CALLS", "line"}, {"JOB_UPDATE", "line"},
		}

		data := &RequestOptions{
			Filters: Filter{
				Cluster:            "clusterDay",
				DayStart:           "2017-01-09",
				DayEnd:             "2017-01-11",
				TimePeriods:        []TimePeriod{},
				Lines:              []string{},
				Users:              []string{},
				OrgStruct:          []string{},
				OrgClusteredGroups: map[string][]string{},
			},
			Charts: charts,
		}

		jsonData, _ := json.Marshal(data)

		requestHeader := "test-kpi-counter-service-" + request_name
		fmt.Println("Request:>", requestHeader)

		resp, body, errs := request.Post(url).
			Set("X-Custom-Header", requestHeader).
			Send(string(jsonData)).
			End()
		if errs != nil {
			fmt.Println(errs)
			os.Exit(1)
		}
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		fmt.Println("response Body:", body)
	}
}

type RequestOptions struct {
	Filters Filter             `json:"filters"`
	Charts  map[string][]Chart `json:"charts"`
}

type Filter struct {
	Cluster            string              `json:"cluster"`
	DayStart           string              `json:"dayStart"`
	DayEnd             string              `json:"dayEnd"`
	TimePeriods        []TimePeriod        `json:"timePeriods"`
	Lines              []string            `json:"lines"`
	Users              []string            `json:"users"`
	OrgStruct          []string            `json:"orgStruct"`
	OrgClusteredGroups map[string][]string `json:"orgClusteredGroups"`
}
type TimePeriod []string
type Chart struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
