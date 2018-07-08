package data

import "encoding/json"

type Day struct {
	appointmentDurationMinutes int
	reservations               []Appointment
}

func NewDay(appointmentDurationMinutes int) *Day {
	day := &Day{
		appointmentDurationMinutes: appointmentDurationMinutes,
		reservations:               make([]Appointment, 24*60/appointmentDurationMinutes),
	}
	startMinute := 0
	for i := 0; i < len(day.reservations); i++ {
		day.reservations[i].StartMinute = startMinute
		startMinute += appointmentDurationMinutes
		day.reservations[i].EndMinute = startMinute
		day.reservations[i].Available = true
	}
	return day
}

func (day *Day) Copy() DaySchedule {
	// Make an independant copy of reservations
	newDay := &Day{
		appointmentDurationMinutes: day.appointmentDurationMinutes,
		reservations:               append([]Appointment(nil), day.reservations...),
	}
	for i := range newDay.reservations {
		// Clear any booked appointments
		if newDay.reservations[i].ClientID != 0 {
			newDay.reservations[i].Available = true
			newDay.reservations[i].ClientID = 0
		}
	}
	return newDay
}

func (day *Day) SetDate(date Date) {
	for i := range day.reservations {
		day.reservations[i].Date = date
	}
}

func (day *Day) Date() Date {
	return day.reservations[0].Date
}

func (day *Day) getIndex(startMinute int) int {
	i := startMinute / day.appointmentDurationMinutes
	if i > len(day.reservations) {
		return -1 // out of range
	}
	return i
}

// Reserve reserves one or more Appointment slots for apt.ClientID.
// If apt.Available is false, the time slot will be blocked off.
func (day *Day) Reserve(apt *Appointment) bool {
	start := day.getIndex(apt.StartMinute)
	end := day.getIndex(apt.EndMinute + day.appointmentDurationMinutes - 1)
	if start < 0 || end < 0 {
		return false
	}
	for i := start; i < end; i++ {
		if !day.reservations[i].Available {
			// Clear the slots we've already reserved
			for j := i - 1; j >= start; j-- {
				day.reservations[j].ClientID = 0
				day.reservations[j].Available = false
			}
		}
		day.reservations[i].ClientID = apt.ClientID
		day.reservations[i].Available = apt.Available
	}
	return true
}

func (day *Day) MarshalJSON() ([]byte, error) {
	return json.Marshal(day.reservations)
}

func (day *Day) AppointmentBlocks() []Appointment {
	return day.reservations
}

func (day *Day) HasAvailable(clientID int) bool {
	for _, apt := range day.reservations {
		if apt.Available || apt.ClientID == clientID {
			return true
		}
	}
	return false
}
