package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
)

func main() {

	requestsTypes := []string{"dates-compare", "charts-count", "charts-compare", "entities-compare"}

	for _, requestType := range requestsTypes {
		fmt.Print("\n ====================================================================\n")

		charts := make(map[string][]Chart)
		charts["wfm"] = []Chart{
			{"T_READY_PLAN", "line"},
			{"PCT_READY", "line"},
		}
		charts["pds"] = []Chart{
			{"JOB_TALK", "line"}, {"JOB_CALLS", "line"}, {"JOB_UPDATE", "line"},
		}

		option := getTypeOption(requestType)

		data := &RequestOptions{
			RequestTypeOptions: option,
			Filters: Filter{
				Cluster:            "clusterHour",
				DayStart:           "2017-01-22",
				DayEnd:             "2017-01-24",
				TimePeriods:        []TimePeriod{},
				Lines:              []string{},
				Users:              getUsers(requestType),
				OrgStruct:          []string{},
				OrgClusteredGroups: map[string][]string{},
			},
			Charts: charts,
		}
		jsonData, _ := json.Marshal(data)

		sendRequest(jsonData, requestType)

	}
}

func sendRequest(data []byte, requestType string) {

	request := gorequest.New()
	url := "http://localhost:8080/chart"

	fmt.Println("Request:>", requestType)

	resp, body, errs := request.Post(url).
		Set("X-Request-type", requestType).
		Send(string(data)).
		End()
	if errs != nil {
		fmt.Println(errs)
		os.Exit(1)
	}

	fmt.Println("response Headers:", resp.Header)
	fmt.Println("response Body:", body)

}

func getUsers(requestType string) []string {
	switch requestType {
	case "entities-compare":
		return []string{"2109", "1546", "5580"}

	default:
		return []string{}

	}
}

func getTypeOption(requestType string) interface{} {

	switch requestType {
	case "dates-compare":
		options := DatesCompareRequestOptions{}
		options.Periods = [][2]string{
			{"2017-01-22", "2017-01-23"},
			{"2017-01-23", "2017-01-24"},
			{"2017-01-07", "2017-01-09"},
		}
		return options
	case "entities-compare":

		options := EntitiesCompareRequestOptions{}
		options.Entity = "Users"
		return options
	default:
		return nil

	}
}

type RequestOptions struct {
	RequestTypeOptions interface{}

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
type DatesCompareRequestOptions struct {
	Periods [][2]string
}

type EntitiesCompareRequestOptions struct {
	Entity string
}
