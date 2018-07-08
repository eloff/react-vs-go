package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/eloff/scheduler/server/api"
	"github.com/eloff/scheduler/server/data"
)

func main() {
	var host string
	var port int
	flag.IntVar(&port, "port", 8080, "port the http server listens on")
	flag.StringVar(&host, "host", "", "host the http server binds to")
	flag.Parse()

	ds := data.NewMemoryStore()
	company := ds.CreateCompany("ACME", 1, 30) // 30 minute appointments
	day := data.NewDay(company.AppointmentMinutes())
	// 9am-5pm available with 12-1pm (lunch) blocked off
	day.Reserve(&data.Appointment{
		EndMinute: 9 * 60, // block off everything before 9am
	})
	day.Reserve(&data.Appointment{
		StartMinute: 12 * 60, // block off 1 hour at noon for lunch
		EndMinute:   13 * 60,
	})
	day.Reserve(&data.Appointment{
		StartMinute: 17 * 60, // block off everything after 5pm
		EndMinute:   24 * 60,
	})
	for i := time.Monday; i < time.Saturday; i++ {
		company.SetAvailable(i, day)
	}

	server, err := api.APIServer(ds, fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("cannot create server: %s", err)
	}

	fmt.Printf("Serving on %s:%d\n", host, port)
	err = server.ListenAndServe()
	if err != nil {
		// TODO check for error returned on close
		log.Fatalf("error in ListenAndServe(): %s", err)
	}
}
