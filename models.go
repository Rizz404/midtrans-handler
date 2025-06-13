package main

import (
	"time"

	"github.com/Rizz404/midtrans-handler/internal/database"
	"github.com/Rizz404/midtrans-handler/internal/enums"
)

type User struct {
	ID             string     `json:"id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	Password       *string    `json:"password,omitempty"`
	Role           enums.Role `json:"role"`
	PhoneNumber    string     `json:"phoneNumber"`
	Address        *string    `json:"address,omitempty"`
	ProfilePicture *string    `json:"profilePicture,omitempty"`
	CreatedAt      any        `json:"createdAt"`
	UpdatedAt      any        `json:"updatedAt"`
}

type Category struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedAt   any     `json:"createdAt"`
	UpdatedAt   any     `json:"updatedAt"`
}

type DenormalizedMenuItem struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	ImageUrl   *string   `json:"imageUrl,omitempty"`
	CategoryId *string   `json:"categoryId,omitempty"`
	Category   *Category `json:"category,omitempty"`
	CreatedAt  any       `json:"createdAt"`
	UpdatedAt  any       `json:"updatedAt"`
}

type OrderItem struct {
	ID              string                `json:"id"`
	OrderId         string                `json:"orderId"`
	MenuItemId      string                `json:"menuItemId"`
	Quantity        int                   `json:"quantity"`
	Price           float64               `json:"price"`
	Total           float64               `json:"total"`
	SpecialRequests *string               `json:"specialRequests,omitempty"`
	MenuItem        *DenormalizedMenuItem `json:"menuItem,omitempty"`
	CreatedAt       any                   `json:"createdAt"`
	UpdatedAt       any                   `json:"updatedAt"`
}

type Order struct {
	ID                  string              `json:"id"`
	UserID              string              `json:"userId"`
	PaymentMethodID     string              `json:"paymentMethodId"`
	OrderType           enums.OrderType     `json:"orderType"`
	Status              enums.OrderStatus   `json:"status"`
	TotalAmount         float64             `json:"totalAmount"`
	PaymentStatus       enums.PaymentStatus `json:"paymentStatus"`
	OrderDate           time.Time           `json:"orderDate"`
	EstimatedReadyTime  *time.Time          `json:"estimatedReadyTime,omitempty"`
	SpecialInstructions *string             `json:"specialInstructions,omitempty"`
	OrderItems          []OrderItem         `json:"orderItems,omitempty"`
	PaymentProof        *string             `json:"paymentProof,omitempty"`
	PaymentCode         *string             `json:"paymentCode,omitempty"`       // Untuk VA, kode Indomaret, dll.
	PaymentDisplayURL   *string             `json:"paymentDisplayUrl,omitempty"` // Untuk URL QRIS, dll.
	PaymentExpiry       *time.Time          `json:"paymentExpiry,omitempty"`     // Waktu kedaluwarsa
	PaymentDetailsRaw   *map[string]any     `json:"paymentDetailsRaw,omitempty"` // Data mentah dari Midtrans
	CreatedAt           any                 `json:"createdAt"`
	UpdatedAt           any                 `json:"updatedAt"`
}

type RestaurantTable struct {
	ID          string         `json:"id"`
	TableNumber string         `json:"tableNumber"`
	Capacity    int            `json:"capacity"`
	IsAvailable bool           `json:"isAvailable"`
	Location    enums.Location `json:"location"`
	CreatedAt   any            `json:"createdAt"`
	UpdatedAt   any            `json:"updatedAt"`
}

type TableReservation struct {
	ID              string                  `json:"id"`
	UserID          string                  `json:"userId"`
	TableID         string                  `json:"tableId"`
	OrderID         string                  `json:"orderId"`
	ReservationTime any                     `json:"reservationTime"`
	Status          enums.ReservationStatus `json:"status"`
	Table           *RestaurantTable        `json:"table,omitempty"`
	CreatedAt       any                     `json:"createdAt"`
	UpdatedAt       any                     `json:"updatedAt"`
}

type PaymentMethod struct {
	ID                 string                  `json:"id"`
	Name               string                  `json:"name"`
	Description        string                  `json:"description"`
	Logo               *string                 `json:"logo,omitempty"`
	PaymentMethodType  enums.PaymentMethodType `json:"paymentMethodType"`
	MidtransIdentifier string                  `json:"midtransIdentifier"`
	CreatedAt          any                     `json:"createdAt"`
	UpdatedAt          any                     `json:"updatedAt"`
}

// * Mapper Functions
func dbUserToUser(dbUser database.User) User {
	return User{
		ID:             dbUser.ID,
		Username:       dbUser.Username,
		Email:          dbUser.Email,
		Password:       dbUser.Password,
		Role:           dbUser.Role,
		PhoneNumber:    dbUser.PhoneNumber,
		Address:        dbUser.Address,
		ProfilePicture: dbUser.ProfilePicture,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
	}
}

func dbUsersToUsers(dbUsers []database.User) []User {
	users := make([]User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = dbUserToUser(dbUser)
	}
	return users
}

func dbCategoryToCategory(dbCategory database.Category) Category {
	return Category{
		ID:          dbCategory.ID,
		Name:        dbCategory.Name,
		Description: dbCategory.Description,
		CreatedAt:   dbCategory.CreatedAt,
		UpdatedAt:   dbCategory.UpdatedAt,
	}
}

func dbCategoriesToCategories(dbCategories []database.Category) []Category {
	categories := make([]Category, len(dbCategories))
	for i, dbCategory := range dbCategories {
		categories[i] = dbCategoryToCategory(dbCategory)
	}
	return categories
}

func dbDenormalizedMenuItemToDenormalizedMenuItem(dbMenuItem database.DenormalizedMenuItem) DenormalizedMenuItem {

	var category *Category
	if dbMenuItem.Category != nil {
		mappedCategory := dbCategoryToCategory(*dbMenuItem.Category)
		category = &mappedCategory
	}

	return DenormalizedMenuItem{
		ID:         dbMenuItem.ID,
		Name:       dbMenuItem.Name,
		Price:      dbMenuItem.Price,
		ImageUrl:   dbMenuItem.ImageUrl,
		CategoryId: dbMenuItem.CategoryId,
		Category:   category,
		CreatedAt:  dbMenuItem.CreatedAt,
		UpdatedAt:  dbMenuItem.UpdatedAt,
	}
}

func dbDenormalizedMenuItemsToDenormalizedMenuItems(dbMenuItems []database.DenormalizedMenuItem) []DenormalizedMenuItem {
	menuItems := make([]DenormalizedMenuItem, len(dbMenuItems))
	for i, dbMenuItem := range dbMenuItems {
		menuItems[i] = dbDenormalizedMenuItemToDenormalizedMenuItem(dbMenuItem)
	}
	return menuItems
}

func dbOrderItemToOrderItem(dbItem database.OrderItem) OrderItem {

	var menuItem *DenormalizedMenuItem
	if dbItem.MenuItem != nil {
		mappedMenuItem := dbDenormalizedMenuItemToDenormalizedMenuItem(*dbItem.MenuItem)
		menuItem = &mappedMenuItem
	}

	return OrderItem{
		ID:              dbItem.ID,
		OrderId:         dbItem.OrderId,
		MenuItemId:      dbItem.MenuItemId,
		Quantity:        dbItem.Quantity,
		Price:           dbItem.Price,
		Total:           dbItem.Total,
		SpecialRequests: dbItem.SpecialRequests,
		MenuItem:        menuItem,
		CreatedAt:       dbItem.CreatedAt,
		UpdatedAt:       dbItem.UpdatedAt,
	}
}

func dbOrderItemsToOrderItems(dbItems []database.OrderItem) []OrderItem {
	items := make([]OrderItem, len(dbItems))
	for i, dbItem := range dbItems {
		items[i] = dbOrderItemToOrderItem(dbItem)
	}
	return items
}

func dbOrderToOrder(dbOrder database.Order) Order {
	return Order{
		ID:                  dbOrder.ID,
		UserID:              dbOrder.UserID,
		PaymentMethodID:     dbOrder.PaymentMethodID,
		OrderType:           dbOrder.OrderType,
		Status:              dbOrder.Status,
		TotalAmount:         dbOrder.TotalAmount,
		PaymentStatus:       dbOrder.PaymentStatus,
		OrderDate:           dbOrder.OrderDate,
		EstimatedReadyTime:  dbOrder.EstimatedReadyTime,
		SpecialInstructions: dbOrder.SpecialInstructions,
		OrderItems:          dbOrderItemsToOrderItems(dbOrder.OrderItems),
		PaymentProof:        dbOrder.PaymentProof,
		PaymentCode:         dbOrder.PaymentCode,
		PaymentDisplayURL:   dbOrder.PaymentDisplayURL,
		PaymentExpiry:       dbOrder.PaymentExpiry,
		PaymentDetailsRaw:   dbOrder.PaymentDetailsRaw,
		CreatedAt:           dbOrder.CreatedAt,
		UpdatedAt:           dbOrder.UpdatedAt,
	}
}

func dbOrdersToOrders(dbOrders []database.Order) []Order {
	orders := make([]Order, len(dbOrders))
	for i, dbOrder := range dbOrders {
		orders[i] = dbOrderToOrder(dbOrder)
	}
	return orders
}

func dbRestaurantTableToRestaurantTable(dbTable database.RestaurantTable) RestaurantTable {
	return RestaurantTable{
		ID:          dbTable.ID,
		TableNumber: dbTable.TableNumber,
		Capacity:    dbTable.Capacity,
		IsAvailable: dbTable.IsAvailable,
		Location:    dbTable.Location,
		CreatedAt:   dbTable.CreatedAt,
		UpdatedAt:   dbTable.UpdatedAt,
	}
}

func dbRestaurantTablesToRestaurantTables(dbTables []database.RestaurantTable) []RestaurantTable {
	tables := make([]RestaurantTable, len(dbTables))
	for i, dbTable := range dbTables {
		tables[i] = dbRestaurantTableToRestaurantTable(dbTable)
	}
	return tables
}

func dbTableReservationToTableReservation(dbReservation database.TableReservation) TableReservation {
	var table *RestaurantTable
	if dbReservation.Table != nil {
		mappedTable := dbRestaurantTableToRestaurantTable(*dbReservation.Table)
		table = &mappedTable
	}

	return TableReservation{
		ID:              dbReservation.ID,
		UserID:          dbReservation.UserID,
		TableID:         dbReservation.TableID,
		OrderID:         dbReservation.OrderID,
		ReservationTime: dbReservation.ReservationTime,
		Status:          dbReservation.Status,
		Table:           table,
		CreatedAt:       dbReservation.CreatedAt,
		UpdatedAt:       dbReservation.UpdatedAt,
	}
}

func dbTableReservationsToTableReservations(dbReservations []database.TableReservation) []TableReservation {
	reservations := make([]TableReservation, len(dbReservations))
	for i, dbReservation := range dbReservations {
		reservations[i] = dbTableReservationToTableReservation(dbReservation)
	}
	return reservations
}

func dbPaymentMethodToPaymentMethod(dbPaymentMethod database.PaymentMethod) PaymentMethod {
	return PaymentMethod{
		ID:                 dbPaymentMethod.ID,
		Name:               dbPaymentMethod.Name,
		Description:        dbPaymentMethod.Description,
		Logo:               dbPaymentMethod.Logo,
		PaymentMethodType:  dbPaymentMethod.PaymentMethodType,
		MidtransIdentifier: dbPaymentMethod.MidtransIdentifier,
		CreatedAt:          dbPaymentMethod.CreatedAt,
		UpdatedAt:          dbPaymentMethod.UpdatedAt,
	}
}

func dbPaymentMethodsToPaymentMethods(dbPaymentMethods []database.PaymentMethod) []PaymentMethod {
	paymentMethods := make([]PaymentMethod, len(dbPaymentMethods))
	for i, dbPaymentMethod := range dbPaymentMethods {
		paymentMethods[i] = dbPaymentMethodToPaymentMethod(dbPaymentMethod)
	}
	return paymentMethods
}
