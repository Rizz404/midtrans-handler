package database

import "time"

type CreateTableReservationRequest struct {
	TableId         string
	ReservationTime time.Time
	Table           *RestaurantTable
}
