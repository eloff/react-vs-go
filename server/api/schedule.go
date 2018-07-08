package api

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ansel1/merry"

	"github.com/eloff/scheduler/server/data"
)

type ScheduleResource struct {
	ds data.DataStore
}

func (ref *ScheduleResource) Register() error {
	http.Handle("/schedule/week.json", handlerWrapper(ref.getWeek))
	http.Handle("/schedule/reserve", handlerWrapper(ref.makeReservation))
	http.Handle("/schedule/week/render", handlerWrapper(ref.renderWeek))
	return nil
}

func (ref *ScheduleResource) getWeek(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	params := r.URL.Query()

	var err error
	names := []string{"companyID", "year", "month", "day", "clientID"}
	vals := make([]int, len(names))
	for i, name := range names {
		vals[i], err = getIntParam(params, name, false)
		if err != nil {
			return nil, err
		}
	}
	companyID := vals[0]
	year := vals[1]
	month := vals[2]
	day := vals[3]
	clientID := vals[4]
	direction := params.Get("direction")
	anchorDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	switch direction {
	case "prev":
		anchorDate = anchorDate.AddDate(0, 0, -7)
	case "next":
		anchorDate = anchorDate.AddDate(0, 0, 7)
	}

	company := ref.ds.GetCompany(companyID)
	if company == nil {
		return nil, merry.WithHTTPCode(errors.New("company not found"), http.StatusNotFound)
	}

	now := time.Now()
	fmt.Println(anchorDate)
	t := anchorDate.AddDate(0, 0, -int(anchorDate.Weekday())) // start of week (Sunday)
	hasAvailable := false
	fmt.Println(anchorDate.Weekday(), t)
	days := make([]data.DaySchedule, 7)
	// Try up to 100 times to fetch the following week if there's nothing available in this week
	for j := 0; !hasAvailable && j < 100; j++ {
		for i := 0; i < 7; i++ {
			fmt.Printf("GetDay %d %s\n", i, t)
			day := company.GetDay(data.DateFromTime(t))
			days[i] = day
			if now.Sub(t) < 24*time.Hour {
				hasAvailable = hasAvailable || day.HasAvailable(clientID)
			}
			t = t.AddDate(0, 0, 1)
		}
	}

	result := &WeekResponse{
		Now:      now,
		ClientID: clientID,
		Company:  company,
		Days:     days,
	}
	return result, nil
}

func (ref *ScheduleResource) renderWeek(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// We'd ordinarily set this up just once in Register() and save it on the ScheduleResource instance
	// but for ease of development we re-create the template object here every time so changes
	// can be made to the template without restarting the server process.
	weekTmpl, err := template.New("week.gohtml").Funcs(template.FuncMap{
		"isEven":     func(i int) bool { return i%2 == 0 },
		"inc":        func(i int) int { return i + 1 },
		"newAptSlot": NewAppointmentWithContext,
	}).ParseFiles("templates/week.gohtml")
	if err != nil {
		return nil, merry.Wrap(err)
	}

	data, err := ref.getWeek(w, r)
	if err != nil {
		return nil, err
	}
	return nil, merry.Wrap(weekTmpl.Execute(w, data))
}

func (ref *ScheduleResource) makeReservation(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var params url.Values
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			return nil, merry.WithHTTPCode(err, http.StatusBadRequest)
		}
		params = r.Form
	} else {
		params = r.URL.Query()
	}

	companyID, err := getIntParam(params, "companyID", false)
	if err != nil {
		return nil, err
	}
	clientID, err := getIntParam(params, "clientID", false)
	if err != nil {
		return nil, err
	}

	company := ref.ds.GetCompany(companyID)
	if company == nil {
		return nil, merry.WithHTTPCode(errors.New("company not found"), http.StatusNotFound)
	}

	dates := make([]int, 5)
	for i, name := range []string{"year", "month", "day", "hour", "minute"} {
		dates[i], err = getIntParam(params, name, i > 2)
		if err != nil {
			return nil, err
		}
	}

	start := dates[3]*60 + dates[4]
	end := start + company.AppointmentMinutes()
	apt := &data.Appointment{
		ClientID:    clientID,
		Date:        data.NewDate(dates[0], time.Month(dates[1]), dates[2]),
		StartMinute: start,
		EndMinute:   end,
	}
	if !company.Reserve(apt) {
		return nil, merry.WithHTTPCode(errors.New("schedule conflict"), http.StatusConflict)
	}

	if params.Get("redirect") != "" {
		return url.Parse(fmt.Sprintf("/schedule/week/render?companyID=%d&clientID=%d&year=%d&month=%d&day=%d",
			companyID, clientID, dates[0], dates[1], dates[2]))
	}

	return nil, nil
}

// getIntParam returns the integer value of the named parameter, or an error
func getIntParam(params url.Values, name string, allowZero bool) (int, error) {
	var i int
	var err error
	val := params.Get(name)
	if val == "" || (!allowZero && val == "0") {
		err = fmt.Errorf("query parameter %s is required", name)
	} else {
		i, err = strconv.Atoi(val)
	}
	return i, merry.WithHTTPCode(err, http.StatusBadRequest)
}

type WeekResponse struct {
	Now      time.Time
	ClientID int
	Company  data.CompanyData
	Days     []data.DaySchedule
}

func (res *WeekResponse) paging(direction string) string {
	start := res.Days[0].Date()
	return fmt.Sprintf("/schedule/week/render?direction=%s&companyID=%d&clientID=%d&year=%d&month=%d&day=%d", direction, res.Company.ID(), res.ClientID, start.Year(), start.Month(), start.Day())
}

func (res *WeekResponse) PrevURL() string {
	return res.paging("prev")
}

func (res *WeekResponse) NextURL() string {
	return res.paging("next")
}

type AppointmentWithContext struct {
	Data *WeekResponse
	Slot *data.Appointment
}

func NewAppointmentWithContext(res *WeekResponse, slot *data.Appointment) *AppointmentWithContext {
	return &AppointmentWithContext{Data: res, Slot: slot}
}

func (ref *AppointmentWithContext) CreateAppointmentURL() string {
	apt := ref.Slot
	return fmt.Sprintf("/schedule/reserve?companyID=%d&clientID=%d&year=%d&month=%d&day=%d&hour=%d&minute=%d",
		ref.Data.Company.ID(), ref.Data.ClientID, apt.Date.Year(), apt.Date.Month(), apt.Date.Day(), apt.Hour(), apt.Minute())
}
