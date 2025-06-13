package database

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// type CreateUserRequest struct {
// 	Name               string
// 	Description        string
// 	Logo               *string
// 	UserType           enums.UserType
// 	MidtransIdentifier string
// }

// type UpdateUserRequest struct {
// 	Name        *string
// 	Description *string
// 	Logo        *string
// }

// func CreateUser(ctx context.Context, client *firestore.Client, request CreateUserRequest) (*User, error) {
// 	docRef := client.Collection("users").NewDoc()

// 	initialData := map[string]any{
// 		"id":                 docRef.ID,
// 		"name":               request.Name,
// 		"description":        request.Description,
// 		"logo":               request.Logo,
// 		"userType":  request.UserType,
// 		"midtransIdentifier": request.MidtransIdentifier,
// 		"createdAt":          firestore.ServerTimestamp,
// 		"updatedAt":          firestore.ServerTimestamp,
// 	}

// 	_, err := docRef.Set(ctx, initialData)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create user: %v", err)
// 	}

// 	docSnapshot, err := docRef.Get(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get new user %s: %v", docRef.ID, err)
// 	}

// 	var newUser User
// 	if err := docSnapshot.DataTo(&newUser); err != nil {
// 		return nil, fmt.Errorf("failed to decode user %s: %v", docRef.ID, err)
// 	}

// 	return &newUser, nil
// }

func GetAllUsers(ctx context.Context, client *firestore.Client) ([]User, error) {
	var users []User
	iter := client.Collection("users").Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate users: %v", err)
		}

		var pm User
		if err := doc.DataTo(&pm); err != nil {
			return nil, fmt.Errorf("failed to decode user: %v", err)
		}
		users = append(users, pm)
	}

	return users, nil
}

func GetUserByID(ctx context.Context, client *firestore.Client, id string) (*User, error) {
	docRef := client.Collection("users").Doc(id)
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %s: %v", id, err)
	}

	var user User
	if err := docSnapshot.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user %s: %v", id, err)
	}

	return &user, nil
}

// func UpdateUser(ctx context.Context, client *firestore.Client, id string, request UpdateUserRequest) (*User, error) {
// 	docRef := client.Collection("users").Doc(id)

// 	updates := []firestore.Update{}
// 	if request.Name != nil {
// 		updates = append(updates, firestore.Update{Path: "name", Value: *request.Name})
// 	}
// 	if request.Description != nil {
// 		updates = append(updates, firestore.Update{Path: "description", Value: *request.Description})
// 	}
// 	if request.Logo != nil {
// 		updates = append(updates, firestore.Update{Path: "logo", Value: request.Logo})
// 	}

// 	if len(updates) == 0 {
// 		return GetUserByID(ctx, client, id)
// 	}

// 	updates = append(updates, firestore.Update{Path: "updatedAt", Value: firestore.ServerTimestamp})

// 	if _, err := docRef.Update(ctx, updates); err != nil {
// 		return nil, fmt.Errorf("failed to update user %s: %v", id, err)
// 	}

// 	return GetUserByID(ctx, client, id)
// }

// func DeleteUser(ctx context.Context, client *firestore.Client, id string) error {
// 	_, err := client.Collection("users").Doc(id).Delete(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete user %s: %v", id, err)
// 	}
// 	return nil
// }
