// This file is for doing http requests to Zermelo
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Response are the fields the server sends back when requestion (wat)
type Response struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	StartRow  int         `json:"startRow"`
	EndRow    int         `json:"endRow"`
	TotalRows int         `json:"totalRows"`
	Data      interface{} `json:"data"`
}

// Appointment is an appointment (xd)
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

// Person is a person
type Person struct {
	Archived                      bool        `json:"archived"`
	City                          string      `json:"city"`
	Code                          string      `json:"code"`
	DateOfBirth                   string      `json:"dateOfBirth"`
	Email                         string      `json:"email"`
	FirstName                     string      `json:"firstName"`
	Gender                        string      `json:"gender"`
	HasPassword                   bool        `json:"hasPassword"`
	HouseNumber                   string      `json:"houseNumber"`
	IsApplicationManager          bool        `json:"isApplicationManager"`
	IsDean                        bool        `json:"isDean"`
	IsEmployee                    bool        `json:"isEmployee"`
	IsFamilyMember                bool        `json:"isFamilyMember"`
	IsMentor                      bool        `json:"isMentor"`
	IsParentTeacherNightScheduler bool        `json:"isParentTeacherNightScheduler"`
	IsSchoolLeader                bool        `json:"isSchoolLeader"`
	IsSchoolScheduler             bool        `json:"isSchoolScheduler"`
	IsSectionLeader               bool        `json:"isSectionLeader"`
	IsStudent                     bool        `json:"isStudent"`
	IsStudentAdministrator        bool        `json:"isStudentAdministrator"`
	IsTeamLeader                  bool        `json:"isTeamLeader"`
	LastName                      string      `json:"lastName"`
	Ldap                          bool        `json:"ldap"`
	MagisterUUID                  interface{} `json:"magisterUUID"`
	PostalCode                    interface{} `json:"postalCode"`
	Prefix                        interface{} `json:"prefix"`
	Roles                         []string    `json:"roles"`
	SchoolInSchoolYears           []int       `json:"schoolInSchoolYears"`
	SomUUID                       interface{} `json:"somUUID"`
	Street                        string      `json:"street"`
	UserPrincipalName             interface{} `json:"userPrincipalName"`
	Username                      string      `json:"username"`
}

// https://bc-enschede.zportal.nl/api/v3/users/~me?fields=code,firstName,prefix,lastName,email,ldap,roles,isStudent,isEmployee,isFamilyMember,isApplicationManager,isSchoolScheduler,isSchoolLeader,isStudentAdministrator,isBranchLeader,isTeamLeader,isSectionLeader,isMentor,isParentTeacherNightScheduler,isDean

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

func fetchAppointments(token string, start int, end int) []Appointment {
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

	// <begin> What's about to come is really really bad, I'm sorry
	var appointments map[string]Response
	json.Unmarshal(body, &appointments)

	// these are the REAL APPOINTMENTS OKAY PELASE HELP
	var realAppointments []Appointment
	// MARSHALL IT INTO JSON AGAIN
	bytes, err := json.Marshal(appointments["response"].Data)
	// WOOP WOOP
	// what even
	if err != nil {
		panic(err)
	}
	// UNMARSHALL IT BACK BUT INTO A STRUCT
	json.Unmarshal(bytes, &realAppointments)
	// I'M TIRED PLEASE HELP

	// <end> Told ya
	return realAppointments
}

// returns a "pretty" json string of person
func fetchMe(token string) string {
	response, err := http.Get(baseurl + "/users/~me?access_token=" + token)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		fmt.Println(response)
		return "{}"
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	// please dont look here
	var person map[string]Response
	json.Unmarshal(body, &person)
	bytes, err := json.MarshalIndent(person["response"].Data.([]interface{})[0], "", "   ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
