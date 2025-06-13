package database

import "time"

type CreateTableReservationParams struct {
	TableId         string           `json:"table_id"`
	ReservationTime time.Time        `json:"reservation_time"`
	Table           *RestaurantTable `json:"table"`
}
