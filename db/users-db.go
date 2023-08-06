package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterUser(email string, password string) (userId string, err error) {
	ctx, client, cancel, err := DbSetup()
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)

	autosDatabase := client.Database("autos")
	usersCollection := autosDatabase.Collection("users")

	insertResult, err := usersCollection.InsertOne(ctx, bson.M{"email": email, "password": password})
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func CheckNewUser(email string) (isNew bool, err error) {
	ctx, client, cancel, err := DbSetup()
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return false, err
	}
	defer client.Disconnect(ctx)

	autosDatabase := client.Database("autos")
	usersCollection := autosDatabase.Collection("users")

	result := usersCollection.FindOne(ctx, bson.M{"email": email})

	var user primitive.D
	result.Decode(&user)

	if user != nil {
		return false, nil
	}

	return true, nil
}

func CheckUserCredentials(email string, password string) (status bool, userId string, err error) {
	ctx, client, cancel, err := DbSetup()
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return false, "", err
	}
	defer client.Disconnect(ctx)

	autosDatabase := client.Database("autos")
	usersCollection := autosDatabase.Collection("users")

	result := usersCollection.FindOne(ctx, bson.M{"email": email, "password": password})

	var user primitive.D
	result.Decode(&user)

	if user == nil {
		return false, "", nil
	}

	return true, user.Map()["_id"].(primitive.ObjectID).Hex(), nil
}
