package enums

// * Enums atau constant lebih tepatnya
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type Location string

const (
	LocationIndoor  Location = "indoor"
	LocationOutdoor Location = "outdoor"
	LocationVIP     Location = "vip"
)

type ReservationStatus string

const (
	StatusReserved  ReservationStatus = "reserved"
	StatusOccupied  ReservationStatus = "occupied"
	StatusCompleted ReservationStatus = "completed"
	StatusCancelled ReservationStatus = "cancelled"
)

type OrderType string

const (
	OrderTypeDineIn   OrderType = "dineIn"
	OrderTypeTakeAway OrderType = "takeAway"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type PaymentStatus string

const (
	PaymentStatusChallenge PaymentStatus = "challenge"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusDeny      PaymentStatus = "deny"
	PaymentStatusFailure   PaymentStatus = "failure"
	PaymentStatusPending   PaymentStatus = "pending"
)

type PaymentMethodType string

const (
	PaymentMethodTypeCash           PaymentMethodType = "cash"
	PaymentMethodTypeCard           PaymentMethodType = "card"
	PaymentMethodTypeDirectDebit    PaymentMethodType = "directDebit"
	PaymentMethodTypeOverTheCounter PaymentMethodType = "overTheCounter"
	PaymentMethodTypeQrCode         PaymentMethodType = "qrCode"
	PaymentMethodTypeVirtualAccount PaymentMethodType = "virtualAccount"
	PaymentMethodTypeEWallet        PaymentMethodType = "eWallet"
	PaymentMethodTypeEchannel       PaymentMethodType = "echannel"
)
