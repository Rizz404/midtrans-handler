package database

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/Rizz404/midtrans-handler/internal/enums"
	"google.golang.org/api/iterator"
)

type CreatePaymentMethodRequest struct {
	Name               string
	Description        string
	Logo               *string
	PaymentMethodType  enums.PaymentMethodType
	MidtransIdentifier string
}

type UpdatePaymentMethodRequest struct {
	Name        *string
	Description *string
	Logo        *string
}

func CreatePaymentMethod(ctx context.Context, client *firestore.Client, request CreatePaymentMethodRequest) (*PaymentMethod, error) {
	docRef := client.Collection("paymentMethods").NewDoc()

	initialData := map[string]any{
		"id":                 docRef.ID,
		"name":               request.Name,
		"description":        request.Description,
		"logo":               request.Logo,
		"paymentMethodType":  request.PaymentMethodType,
		"midtransIdentifier": request.MidtransIdentifier,
		"createdAt":          firestore.ServerTimestamp,
		"updatedAt":          firestore.ServerTimestamp,
	}

	_, err := docRef.Set(ctx, initialData)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment method: %v", err)
	}

	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get new payment method %s: %v", docRef.ID, err)
	}

	var newPaymentMethod PaymentMethod
	if err := docSnapshot.DataTo(&newPaymentMethod); err != nil {
		return nil, fmt.Errorf("failed to decode payment method %s: %v", docRef.ID, err)
	}

	return &newPaymentMethod, nil
}

func BulkCreatePaymentMethods(ctx context.Context, client *firestore.Client, requests []CreatePaymentMethodRequest) ([]PaymentMethod, error) {
	if len(requests) == 0 {
		return []PaymentMethod{}, nil
	}

	batch := client.Batch()
	newDocIDs := make([]string, 0, len(requests))

	for _, request := range requests {
		docRef := client.Collection("paymentMethods").NewDoc()
		newDocIDs = append(newDocIDs, docRef.ID)

		initialData := map[string]any{
			"id":                 docRef.ID,
			"name":               request.Name,
			"description":        request.Description,
			"logo":               request.Logo,
			"paymentMethodType":  request.PaymentMethodType,
			"midtransIdentifier": request.MidtransIdentifier,
			"createdAt":          firestore.ServerTimestamp,
			"updatedAt":          firestore.ServerTimestamp,
		}
		batch.Set(docRef, initialData)
	}

	_, err := batch.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit batch create payment methods: %v", err)
	}

	var createdPaymentMethods []PaymentMethod
	iter := client.Collection("paymentMethods").Where("id", "in", newDocIDs).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate newly created payment methods: %v", err)
		}

		var pm PaymentMethod
		if err := doc.DataTo(&pm); err != nil {
			return nil, fmt.Errorf("failed to decode newly created payment method: %v", err)
		}
		createdPaymentMethods = append(createdPaymentMethods, pm)
	}

	return createdPaymentMethods, nil
}

func GetAllPaymentMethods(ctx context.Context, client *firestore.Client) ([]PaymentMethod, error) {
	var paymentMethods []PaymentMethod
	iter := client.Collection("paymentMethods").Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate payment methods: %v", err)
		}

		var pm PaymentMethod
		if err := doc.DataTo(&pm); err != nil {
			return nil, fmt.Errorf("failed to decode payment method: %v", err)
		}
		paymentMethods = append(paymentMethods, pm)
	}

	return paymentMethods, nil
}

func GetPaymentMethodByID(ctx context.Context, client *firestore.Client, id string) (*PaymentMethod, error) {
	docRef := client.Collection("paymentMethods").Doc(id)
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method %s: %v", id, err)
	}

	var paymentMethod PaymentMethod
	if err := docSnapshot.DataTo(&paymentMethod); err != nil {
		return nil, fmt.Errorf("failed to decode payment method %s: %v", id, err)
	}

	return &paymentMethod, nil
}

func UpdatePaymentMethod(ctx context.Context, client *firestore.Client, id string, request UpdatePaymentMethodRequest) (*PaymentMethod, error) {
	docRef := client.Collection("paymentMethods").Doc(id)

	updates := []firestore.Update{}
	if request.Name != nil {
		updates = append(updates, firestore.Update{Path: "name", Value: *request.Name})
	}
	if request.Description != nil {
		updates = append(updates, firestore.Update{Path: "description", Value: *request.Description})
	}
	if request.Logo != nil {
		updates = append(updates, firestore.Update{Path: "logo", Value: request.Logo})
	}

	if len(updates) == 0 {
		return GetPaymentMethodByID(ctx, client, id)
	}

	updates = append(updates, firestore.Update{Path: "updatedAt", Value: firestore.ServerTimestamp})

	if _, err := docRef.Update(ctx, updates); err != nil {
		return nil, fmt.Errorf("failed to update payment method %s: %v", id, err)
	}

	return GetPaymentMethodByID(ctx, client, id)
}

func DeletePaymentMethod(ctx context.Context, client *firestore.Client, id string) error {
	_, err := client.Collection("paymentMethods").Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete payment method %s: %v", id, err)
	}
	return nil
}
