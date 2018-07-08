package data

import (
	"fmt"
	"sync"
	"time"
)

type Company struct {
	mutex                      sync.Mutex
	id                         int
	appointmentDurationMinutes int
	name                       string
	availableSchedule          [7]DaySchedule
	schedule                   map[Date]DaySchedule
	proposed                   []*Appointment
}

func (com *Company) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(
		`{"ID": %d, "Name": %q, "AppointmentDurationMinutes": %d}`,
		com.id, com.name, com.appointmentDurationMinutes)), nil
}

func (com *Company) ID() int {
	return com.id
}

func (com *Company) AppointmentMinutes() int {
	return com.appointmentDurationMinutes
}

func (com *Company) SetAvailable(day time.Weekday, schedule DaySchedule) {
	com.mutex.Lock()
	com.availableSchedule[day] = schedule
	com.mutex.Unlock()
}
func (com *Company) getOrCreateDay(date Date) DaySchedule {
	day := com.schedule[date]
	if day == nil {
		template := com.availableSchedule[date.Weekday()]
		if template == nil {
			day = NewDay(com.appointmentDurationMinutes)
			day.Reserve(&Appointment{
				EndMinute: minutesInDay,
			})
		} else {
			day = template.Copy()
		}
		day.SetDate(date)
		com.schedule[date] = day
	}
	return day
}

func (com *Company) GetDay(date Date) DaySchedule {
	com.mutex.Lock()
	defer com.mutex.Unlock()
	return com.getOrCreateDay(date)
}

func (com *Company) Reserve(apt *Appointment) bool {
	// Don't allow making appointments in the past
	if apt.Time().Before(time.Now()) {
		return false
	}
	com.mutex.Lock()
	defer com.mutex.Unlock()
	day := com.getOrCreateDay(apt.Date)
	return day.Reserve(apt)
}

func NewCompany(name string, id, appointmentDurationMinutes int) *Company {
	return &Company{
		id:   id,
		name: name,
		appointmentDurationMinutes: appointmentDurationMinutes,
		schedule:                   map[Date]DaySchedule{},
	}
}
