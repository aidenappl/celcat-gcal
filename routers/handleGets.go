package routers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
)

type HandledResponse []struct {
	ID              string      `json:"id"`
	Start           string      `json:"start"`
	End             string      `json:"end"`
	AllDay          bool        `json:"allDay"`
	Description     string      `json:"description"`
	BackgroundColor string      `json:"backgroundColor"`
	TextColor       string      `json:"textColor"`
	Department      string      `json:"department"`
	Faculty         interface{} `json:"faculty"`
	EventCategory   string      `json:"eventCategory"`
	Sites           []string    `json:"sites"`
	Modules         []string    `json:"modules"`
	RegisterStatus  int         `json:"registerStatus"`
	StudentMark     int         `json:"studentMark"`
	Custom1         interface{} `json:"custom1"`
	Custom2         interface{} `json:"custom2"`
	Custom3         interface{} `json:"custom3"`
}

func HandleGetCalendar(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")

	if name == "" {
		http.Error(w, "a valid name is required", http.StatusBadRequest)
		return
	}

	data := url.Values{}
	data.Set("start", "2022-09-01")
	data.Set("end", "2022-12-17")
	data.Set("resType", "104")
	data.Set("calView", "month")
	data.Set("federationIds[]", name)

	apiUrl := "https://timetable.nchlondon.ac.uk"
	resource := "/Home/GetCalendarData"

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	responseOb := HandledResponse{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(bodyBytes, &responseOb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if responseOb == nil || len(responseOb) == 0 {
		http.Error(w, "no events found", http.StatusConflict)
		return
	}

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	cal.SetCalscale("GREGORIAN")
	cal.SetName("NCH Timetable")
	cal.SetDescription("NCH Timetable")
	cal.SetXWRCalName("NCH Timetable")
	cal.SetXWRCalDesc("NCH Timetable for " + name)
	cal.SetVersion("2.0")

	for _, respEv := range responseOb {

		stTime, err := time.Parse(time.RFC3339, respEv.Start+"Z")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		endTime, err := time.Parse(time.RFC3339, respEv.End+"Z")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		formattedDescription := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(respEv.Description, "\n", ""), "<br />", "\r\n"), "\r", "")

		pattern := regexp.MustCompile(`\([^)]*\)`)
		bracketPattern := regexp.MustCompile(`\[[^)]*\]`)
		fountTitlePTRN := pattern.ReplaceAllString(respEv.Modules[0], "")
		fountTitle := bracketPattern.ReplaceAllString(fountTitlePTRN, "")
		foundBracket := bracketPattern.FindString(fountTitlePTRN)
		formattedTitle := strings.TrimSpace(fountTitle)

		formattedDescription = "ClassID: " + foundBracket + " " + formattedDescription

		event := cal.AddEvent(respEv.ID)
		event.SetStartAt(stTime.Add(-time.Hour * 1))
		event.SetDtStampTime(time.Now())
		event.SetEndAt(endTime.Add(-time.Hour * 1))

		// Check is Custom 1, 2, 3 has a value
		var urlVal string
		if respEv.Custom1 != nil {
			urlVal = respEv.Custom1.(string)
		} else if respEv.Custom2 != nil {
			urlVal = respEv.Custom2.(string)
		} else if respEv.Custom3 != nil {
			urlVal = respEv.Custom3.(string)
		} else {
			urlVal = "https://timetable.nchlondon.ac.uk/"
		}

		event.SetURL(urlVal)
		event.SetClass("PUBLIC")
		event.SetSequence(0)
		event.SetSummary(formattedTitle)

		var loc string

		if len(respEv.Sites) > 0 {
			loc = respEv.Sites[0]
		} else {
			loc = "Northeastern London"
		}

		event.SetLocation(loc)
		event.SetDescription(formattedDescription)
	}

	serializedCal := cal.Serialize()

	dtx := []byte(serializedCal)

	w.Header().Set("Content-Type", "text/calendar")
	w.WriteHeader(http.StatusOK)
	w.Write(dtx)

}
