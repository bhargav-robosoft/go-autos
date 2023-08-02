package db

import (
	"autos/entity"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbSetup() (context.Context, *mongo.Client, context.CancelFunc, error) {
	var dbUrl string
	dbUrl = "mongodb+srv://" + os.Getenv("MONGO_USERNAME") + ":" + os.Getenv("MONGO_PASSWORD") + "@cluster1.5jqwhvz.mongodb.net/?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(dbUrl))
	if err != nil {
		return nil, nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	return ctx, client, cancel, nil
}

func GetAutos() error {
	ctx, client, cancel, err := DbSetup()
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	autosDatabase := client.Database("autos")
	carsCollection := autosDatabase.Collection("cars")

	cursor, err := carsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var cars []bson.M
	if err = cursor.All(ctx, &cars); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cars)

	return nil
}

func InsertAuto(car entity.Car) (string, error) {
	ctx, client, cancel, err := DbSetup()
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)

	autosDatabase := client.Database("autos")
	carsCollection := autosDatabase.Collection("cars")

	result, err := carsCollection.InsertOne(ctx, bson.D{
		{Key: "id", Value: car.Id},
		{Key: "name", Value: car.Name},
		{Key: "company", Value: car.Company},
		{Key: "model", Value: car.Model},
	})

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
