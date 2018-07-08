package data

import (
	"fmt"
	"sync"
	"time"
)

// This file contains an in-memory implementation of the data store interfaces
// it is useful for mocking during tests and for development.
// Given that this is a not a real project, we just use the memory implementation.
// In a real project there would also be an implementation using a permanent database.
var _ DataStore = (*MemoryStore)(nil)
var _ CompanyData = (*Company)(nil)
var _ DaySchedule = (*Day)(nil)

const minutesInDay = 24 * 60

type MemoryStore struct {
	companies map[int]*Company
	mutex     sync.Mutex
}

type Client struct {
	id   int
	name string
}

type Appointment struct {
	ClientID    int
	StartMinute int
	EndMinute   int
	Date        Date
	Available   bool
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		companies: map[int]*Company{},
	}
}

func (mem *MemoryStore) CreateCompany(name string, id, appointmentDurationMinutes int) CompanyData {
	com := NewCompany(name, id, appointmentDurationMinutes)
	mem.mutex.Lock()
	mem.companies[com.ID()] = com
	mem.mutex.Unlock()
	return com
}

func (mem *MemoryStore) GetCompany(id int) CompanyData {
	mem.mutex.Lock()
	com, ok := mem.companies[id]
	mem.mutex.Unlock()
	if ok {
		return com
	}
	return nil
}

func (apt *Appointment) Hour() int {
	return apt.StartMinute / 60
}

func (apt *Appointment) Minute() int {
	return apt.StartMinute % 60
}

func (apt *Appointment) Time() time.Time {
	return time.Date(apt.Date.Year(), apt.Date.Month(), apt.Date.Day(), 0, apt.StartMinute, 0, 0, time.Now().Location())
}

func (apt *Appointment) DisplayTime() string {
	return fmt.Sprintf("%02d:%02d", apt.Hour(), apt.Minute())
}
