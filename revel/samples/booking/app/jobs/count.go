package jobs

import (
	"fmt"

	"github.com/acsellers/multitemplate/revel/samples/booking/app/controllers"
	"github.com/acsellers/multitemplate/revel/samples/booking/app/models"
	"github.com/revel/revel"
	"github.com/revel/revel/modules/jobs/app/jobs"
)

// Periodically count the bookings in the database.
type BookingCounter struct{}

func (c BookingCounter) Run() {
	bookings, err := controllers.Dbm.Select(models.Booking{},
		`select * from Booking`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("There are %d bookings.\n", len(bookings))
}

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("@every 1m", BookingCounter{})
	})
}
