package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Rizz404/midtrans-handler/internal/enums"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type CreateOrderWithPaymentRequest struct {
	UserID              string
	PaymentMethodID     string
	OrderType           enums.OrderType
	EstimatedReadyTime  *time.Time
	SpecialInstructions *string
	OrderItems          []OrderItem
}

type UpdateOrderRequest struct {
	OrderStatus   *enums.OrderStatus
	PaymentStatus *enums.PaymentStatus
}

func CreateOrderWithPayment(
	ctx context.Context,
	firestoreClient *firestore.Client,
	midtransClient *coreapi.Client,
	req CreateOrderWithPaymentRequest,
) (*Order, error) {
	user, err := GetUserByID(ctx, firestoreClient, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	paymentMethod, err := GetPaymentMethodByID(ctx, firestoreClient, req.PaymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method: %v", err)
	}

	var totalAmount float64
	for _, item := range req.OrderItems {
		totalAmount += item.Total
	}

	orderID := firestoreClient.Collection("orders").NewDoc().ID
	chargeReq := buildMidtransChargeRequest(orderID, totalAmount, user, paymentMethod, req.OrderItems)

	chargeResp, chargeErr := midtransClient.ChargeTransaction(chargeReq)
	if chargeErr != nil {
		return nil, fmt.Errorf("midtrans charge failed: %v", chargeErr.GetMessage())
	}

	// --- PERUBAHAN UTAMA: Buat map untuk ditulis ke Firestore ---
	orderData := map[string]any{
		"id":                  orderID,
		"userId":              req.UserID,
		"paymentMethodId":     req.PaymentMethodID,
		"orderType":           req.OrderType,
		"status":              enums.OrderStatusPending,
		"paymentStatus":       enums.PaymentStatusUnpaid,
		"totalAmount":         totalAmount,
		"estimatedReadyTime":  req.EstimatedReadyTime,
		"specialInstructions": req.SpecialInstructions,
		"orderItems":          req.OrderItems, // Client Go akan menangani marshaling struct ini
		"paymentDetailsRaw":   chargeResp,
		"createdAt":           firestore.ServerTimestamp, // Sekarang menjadi map value, INI BENAR
		"updatedAt":           firestore.ServerTimestamp, // Sekarang menjadi map value, INI BENAR
	}

	// Parsing dan tambahkan data pembayaran ke map
	if len(chargeResp.VaNumbers) > 0 {
		orderData["paymentCode"] = chargeResp.VaNumbers[0].VANumber
	}
	if chargeResp.PaymentCode != "" {
		orderData["paymentCode"] = chargeResp.PaymentCode
	}
	for _, action := range chargeResp.Actions {
		if action.Name == "generate-qr-code" || action.Name == "deeplink-redirect" {
			orderData["paymentDisplayUrl"] = action.URL
			break
		}
	}
	expiryTime, err := time.Parse("2006-01-02 15:04:05", chargeResp.ExpiryTime)
	if err == nil {
		orderData["paymentExpiry"] = expiryTime
	}

	// Tulis map ke Firestore, bukan struct
	_, err = firestoreClient.Collection("orders").Doc(orderID).Set(ctx, orderData)
	if err != nil {
		log.Printf("CRITICAL: Order %s created at Midtrans but failed to save to Firestore: %v", orderID, err)
		return nil, fmt.Errorf("payment created but failed to save order: %v", err)
	}

	// Buat struct Order untuk dikembalikan ke handler, karena penulisan sudah sukses
	finalOrder := &Order{
		ID:                  orderID,
		UserID:              req.UserID,
		PaymentMethodID:     req.PaymentMethodID,
		OrderType:           req.OrderType,
		Status:              enums.OrderStatusPending,
		PaymentStatus:       enums.PaymentStatusUnpaid,
		TotalAmount:         totalAmount,
		EstimatedReadyTime:  req.EstimatedReadyTime,
		SpecialInstructions: req.SpecialInstructions,
		OrderItems:          req.OrderItems,
		PaymentDetailsRaw:   toMapPointer(orderData["paymentDetailsRaw"]),
		PaymentCode:         toStringPointer(orderData["paymentCode"]),
		PaymentDisplayURL:   toStringPointer(orderData["paymentDisplayUrl"]),
		PaymentExpiry:       toTimePointer(orderData["paymentExpiry"]),
	}

	return finalOrder, nil
}

// toMapPointer safely converts an any to *map[string]any
func toMapPointer(val any) *map[string]any {
	if m, ok := val.(map[string]any); ok {
		return &m
	}
	return nil
}

// toStringPointer safely converts an any to *string
func toStringPointer(val any) *string {
	if s, ok := val.(string); ok {
		return &s
	}
	return nil
}

// toTimePointer safely converts an any to *time.Time
func toTimePointer(val any) *time.Time {
	if t, ok := val.(time.Time); ok {
		return &t
	}
	return nil
}

func buildMidtransChargeRequest(orderID string, totalAmount float64, user *User, paymentMethod *PaymentMethod, items []OrderItem) *coreapi.ChargeReq {
	var midtransItems []midtrans.ItemDetails
	for _, item := range items {
		var itemName string
		if item.MenuItem != nil {
			itemName = item.MenuItem.Name
		}
		midtransItems = append(midtransItems, midtrans.ItemDetails{
			ID:    item.MenuItemId,
			Price: int64(item.Price),
			Qty:   int32(item.Quantity),
			Name:  itemName,
		})
	}
	chargeReq := &coreapi.ChargeReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(totalAmount),
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: user.Username,
			LName: user.Username,
			Email: user.Email,
			Phone: user.PhoneNumber,
		},
		Items: &midtransItems,
	}
	switch paymentMethod.PaymentMethodType {
	case enums.PaymentMethodTypeVirtualAccount:
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{Bank: midtrans.Bank(paymentMethod.MidtransIdentifier)}
	case enums.PaymentMethodTypeEWallet:
		switch paymentMethod.MidtransIdentifier {
		case "gopay":
			chargeReq.PaymentType = coreapi.PaymentTypeGopay
		case "shopeepay":
			chargeReq.PaymentType = coreapi.PaymentTypeShopeepay
			chargeReq.ShopeePay = &coreapi.ShopeePayDetails{CallbackUrl: "https://your-domain.com/shopeepay/callback"}
		}
	case enums.PaymentMethodTypeQrCode:
		chargeReq.PaymentType = coreapi.PaymentTypeQris
		chargeReq.Qris = &coreapi.QrisDetails{Acquirer: paymentMethod.MidtransIdentifier}
	case enums.PaymentMethodTypeOverTheCounter:
		chargeReq.PaymentType = coreapi.PaymentTypeConvenienceStore
		chargeReq.ConvStore = &coreapi.ConvStoreDetails{Store: paymentMethod.MidtransIdentifier}
	}
	return chargeReq
}

func GetOrderByID(ctx context.Context, client *firestore.Client, id string) (*Order, error) {
	docRef := client.Collection("orders").Doc(id)
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get order %s: %v", id, err)
	}
	var order Order
	if err := docSnapshot.DataTo(&order); err != nil {
		return nil, fmt.Errorf("failed to decode order %s: %v", id, err)
	}
	return &order, nil
}

func UpdateOrder(ctx context.Context, client *firestore.Client, id string, request UpdateOrderRequest) (*Order, error) {
	docRef := client.Collection("orders").Doc(id)
	updates := []firestore.Update{}
	if request.OrderStatus != nil {
		updates = append(updates, firestore.Update{Path: "status", Value: *request.OrderStatus})
	}
	if request.PaymentStatus != nil {
		updates = append(updates, firestore.Update{Path: "paymentStatus", Value: *request.PaymentStatus})
	}
	if len(updates) == 0 {
		return GetOrderByID(ctx, client, id)
	}
	updates = append(updates, firestore.Update{Path: "updatedAt", Value: firestore.ServerTimestamp})
	if _, err := docRef.Update(ctx, updates); err != nil {
		return nil, fmt.Errorf("failed to update order %s: %v", id, err)
	}
	return GetOrderByID(ctx, client, id)
}
