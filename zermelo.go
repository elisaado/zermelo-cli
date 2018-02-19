package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type AppointmentResponse struct {
	Status    int           `json:"status"`
	Message   string        `json:"message"`
	StartRow  int           `json:"startRow"`
	EndRow    int           `json:"endRow"`
	TotalRows int           `json:"totalRows"`
	Data      []Appointment `json:"data"`
}

type Appointment struct {
	ID                  int      `json:"id"`
	AppointmentInstance int      `json:"appointmentInstance"`
	Start               int      `json:"start"`
	End                 int      `json:"end"`
	StartTimeSlot       int      `json:"startTimeSlot"`
	EndTimeSlot         int      `json:"endTimeSlot"`
	Subjects            []string `json:"subjects"`
	Teachers            []string `json:"teachers"`
	Groups              []string `json:"groups"`
	GroupsInDepartments []int    `json:"groupsInDepartment"`
	Locations           []string `json:"locations"`
	LocationsOfBranch   []int    `json:"locationsOfBranch"`
	Type                string   `json:"type"`
	Remark              string   `json:"remark"`
	Valid               bool     `json:"valid"`
	Cancelled           bool     `json:"cancelled"`
	Modified            bool     `json:"modified"`
	Moved               bool     `json:"moved"`
	New                 bool     `json:"new"`
	ChangeDescription   string   `json:"changeDescription"`
}

// TODO:
// Write functions that do get rqequests to the zermelo api and implement some funcs, maybe a class <- struct zermelo or smth idk

var baseurl string

func fetchAuthToken(organisation string, code int) string {
	// Do a http request to get the auth token
	response, err := http.PostForm("https://"+organisation+".zportal.nl/api/v3"+"/oauth/token", url.Values{"grant_type": {"authorization_code"}, "code": {strconv.Itoa(code)}})
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return ""
	}

	// Read the body from the response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Parse json and return token
	var jsonBody map[string]string
	json.Unmarshal(body, &jsonBody) // We expect an error, because the expires_at field is [string]int, we don't need it anyway though
	return jsonBody["access_token"]
}

func fetchAppointments(organisation string, token string, start int, end int) []Appointment {
	response, err := http.Get(baseurl + "/appointments?user=~me&start=" + strconv.Itoa(start) + "&end=" + strconv.Itoa(end) + "&access_token=" + token)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		fmt.Println(response)
		return []Appointment{}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var appointments map[string]AppointmentResponse
	json.Unmarshal(body, &appointments)
	return appointments["response"].Data
}
