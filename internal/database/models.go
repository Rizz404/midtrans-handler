package database

import (
	"time"

	"github.com/Rizz404/midtrans-handler/internal/enums"
)

type User struct {
	ID             string     `firestore:"id"`
	Username       string     `firestore:"username"`
	Email          string     `firestore:"email"`
	Password       *string    `firestore:"password,omitempty"`
	Role           enums.Role `firestore:"role"`
	PhoneNumber    string     `firestore:"phoneNumber"`
	Address        *string    `firestore:"address,omitempty"`
	ProfilePicture *string    `firestore:"profilePicture,omitempty"`
	CreatedAt      any        `firestore:"createdAt"`
	UpdatedAt      any        `firestore:"updatedAt"`
}

type Category struct {
	ID   string `firestore:"id"`
	Name string `firestore:"name"`
	// ! selalu pake pointer buat nullable ya
	Description *string `firestore:"description,omitempty"`
	CreatedAt   any     `firestore:"createdAt"`
	UpdatedAt   any     `firestore:"updatedAt"`
}

type DenormalizedMenuItem struct {
	ID         string    `firestore:"id"`
	Name       string    `firestore:"name"`
	Price      float64   `firestore:"price"`
	ImageUrl   *string   `firestore:"imageUrl,omitempty"`
	CategoryId *string   `firestore:"categoryId,omitempty"`
	Category   *Category `firestore:"category,omitempty"`
	CreatedAt  any       `firestore:"createdAt"`
	UpdatedAt  any       `firestore:"updatedAt"`
}

type CartItem struct {
	ID         string                `firestore:"id"`
	MenuItemId string                `firestore:"menuItemId"`
	UserId     string                `firestore:"userId"`
	Quantity   int                   `firestore:"quantity"`
	MenuItem   *DenormalizedMenuItem `firestore:"menuItem,omitempty"`
	CreatedAt  any                   `firestore:"createdAt"`
	UpdatedAt  any                   `firestore:"updatedAt"`
}

type OrderItem struct {
	ID              string                `firestore:"id"`
	OrderId         string                `firestore:"orderId"`
	MenuItemId      string                `firestore:"menuItemId"`
	Quantity        int                   `firestore:"quantity"`
	Price           float64               `firestore:"price"`
	Total           float64               `firestore:"total"`
	SpecialRequests *string               `firestore:"specialRequests,omitempty"`
	MenuItem        *DenormalizedMenuItem `firestore:"menuItem,omitempty"`
	CreatedAt       any                   `firestore:"createdAt"`
	UpdatedAt       any                   `firestore:"updatedAt"`
}

type Order struct {
	ID                  string              `firestore:"id"`
	UserID              string              `firestore:"userId"`
	PaymentMethodID     string              `firestore:"paymentMethodId"`
	OrderType           enums.OrderType     `firestore:"orderType"`
	Status              enums.OrderStatus   `firestore:"status"`
	TotalAmount         float64             `firestore:"totalAmount"`
	PaymentStatus       enums.PaymentStatus `firestore:"paymentStatus"`
	OrderDate           time.Time           `firestore:"orderDate"`
	EstimatedReadyTime  *time.Time          `firestore:"estimatedReadyTime,omitempty"`
	SpecialInstructions *string             `firestore:"specialInstructions,omitempty"`
	OrderItems          []OrderItem         `firestore:"orderItems,omitempty"`
	PaymentProof        *string             `firestore:"paymentProof,omitempty"`
	PaymentCode         *string             `firestore:"paymentCode,omitempty"`       // Untuk VA, kode Indomaret, dll.
	PaymentDisplayURL   *string             `firestore:"paymentDisplayUrl,omitempty"` // Untuk URL QRIS, dll.
	PaymentExpiry       *time.Time          `firestore:"paymentExpiry,omitempty"`     // Waktu kedaluwarsa
	PaymentDetailsRaw   *map[string]any     `firestore:"paymentDetailsRaw,omitempty"` // Data mentah dari Midtrans
	CreatedAt           any                 `firestore:"createdAt"`
	UpdatedAt           any                 `firestore:"updatedAt"`
}

type RestaurantTable struct {
	ID          string         `firestore:"id"`
	TableNumber string         `firestore:"tableNumber"`
	Capacity    int            `firestore:"capacity"`
	IsAvailable bool           `firestore:"isAvailable"`
	Location    enums.Location `firestore:"location"`
	CreatedAt   any            `firestore:"createdAt"`
	UpdatedAt   any            `firestore:"updatedAt"`
}

type TableReservation struct {
	ID              string                  `firestore:"id"`
	UserID          string                  `firestore:"userId"`
	TableID         string                  `firestore:"tableId"`
	OrderID         string                  `firestore:"orderId"`
	ReservationTime any                     `firestore:"reservationTime"`
	Status          enums.ReservationStatus `firestore:"status"`
	Table           *RestaurantTable        `firestore:"table,omitempty"`
	CreatedAt       any                     `firestore:"createdAt"`
	UpdatedAt       any                     `firestore:"updatedAt"`
}

type PaymentMethod struct {
	ID                        string                  `firestore:"id"`
	Name                      string                  `firestore:"name"`
	Description               string                  `firestore:"description"`
	Logo                      *string                 `firestore:"logo,omitempty"`
	PaymentMethodType         enums.PaymentMethodType `firestore:"paymentMethodType"`
	MidtransIdentifier        *string                 `firestore:"midtransIdentifier"`
	MinimumAmount             float64                 `firestore:"minimumAmount"`
	MaximumAmount             float64                 `firestore:"maximumAmount"`
	AdminPaymentCode          *string                 `firestore:"adminPaymentCode,omitempty"`
	AdminPaymentQrCodePicture *string                 `firestore:"adminPaymentQrCodePicture,omitempty"`
	CreatedAt                 any                     `firestore:"createdAt"`
	UpdatedAt                 any                     `firestore:"updatedAt"`
}
