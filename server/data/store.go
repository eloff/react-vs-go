package data

import (
	"encoding/json"
	"time"
)

type DataStore interface {
	CreateCompany(name string, id, appointmentDurationMinutes int) CompanyData
	GetCompany(id int) CompanyData
}

type CompanyData interface {
	json.Marshaler
	ID() int
	AppointmentMinutes() int
	SetAvailable(day time.Weekday, schedule DaySchedule)
	GetDay(date Date) DaySchedule
	Reserve(apt *Appointment) bool
}

type DaySchedule interface {
	json.Marshaler
	Copy() DaySchedule
	SetDate(date Date)
	Date() Date
	Reserve(apt *Appointment) bool
	HasAvailable(clientID int) bool
}
