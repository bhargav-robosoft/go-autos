package db

import (
	"autos/autoserror"
	"autos/entity"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReadCars() (cars []entity.Car, err error) {
	ctx, client, cancel, err := DbSetup()
	if err != nil {
		return []entity.Car{}, err
	}
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return []entity.Car{}, err
	}
	defer client.Disconnect(ctx)

	autosDatabase := client.Database("autos")
	carsCollection := autosDatabase.Collection("cars")

	cursor, err := carsCollection.Find(ctx, bson.M{})
	if err != nil {
		return []entity.Car{}, err
	}

	cars = []entity.Car{}
	for cursor.Next(ctx) {
		var car entity.Car
		err = cursor.Decode(&car)
		if err != nil {
			return []entity.Car{}, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

func CreateCar(car entity.NewCarRequest, userId string) (string, error) {
	ctx, client, cancel, err := DbSetup()
	if err != nil {
		return "", err
	}
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)

	autosDatabase := client.Database("autos")
	carsCollection := autosDatabase.Collection("cars")
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return "", err
	}

	result, err := carsCollection.InsertOne(ctx, bson.D{
		{Key: "name", Value: car.Name},
		{Key: "company", Value: car.Company},
		{Key: "model", Value: car.Model},
		{Key: "userId", Value: userObjId},
	})

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func UpdateCar(car entity.ModifyCarRequest, userId string) (updatedId string, err error) {
	ctx, client, cancel, err := DbSetup()
	if err != nil {
		return "", err
	}
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)

	autosDatabase := client.Database("autos")
	carsCollection := autosDatabase.Collection("cars")
	carObjId, err := primitive.ObjectIDFromHex(car.Id)
	if err != nil {
		return "", err
	}
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return "", err
	}

	filter := bson.D{
		{Key: "_id", Value: carObjId},
	}

	carResult := carsCollection.FindOne(ctx, filter)
	var dbCar entity.Car
	err = carResult.Decode(&dbCar)
	if err != nil {
		return "", &autoserror.CustomError{
			Message: fmt.Sprintf("No car with id %s", car.Id),
			Status:  http.StatusNotFound,
		}
	}
	if dbCar.UserId != userId {
		return "", &autoserror.CustomError{
			Message: fmt.Sprintf("No edit access for car %s", car.Id),
			Status:  http.StatusForbidden,
		}
	}

	carName := car.Name
	carCompany := car.Company
	carModel := car.Model

	if carName == "" {
		carName = dbCar.Name
	}
	if carCompany == "" {
		carCompany = dbCar.Company
	}
	if carModel == 0 {
		carModel = dbCar.Model
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: carName},
			{Key: "company", Value: carCompany},
			{Key: "model", Value: carModel},
		}},
	}

	filter = bson.D{
		{Key: "_id", Value: carObjId},
		{Key: "userId", Value: userObjId},
	}

	_, err = carsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return dbCar.Id, nil
}

func DeleteCar(carId string, userId string) (err error) {
	ctx, client, cancel, err := DbSetup()
	if err != nil {
		return err
	}
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	autosDatabase := client.Database("autos")
	carsCollection := autosDatabase.Collection("cars")
	carObjId, err := primitive.ObjectIDFromHex(carId)
	if err != nil {
		return err
	}
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: carObjId},
	}

	carResult := carsCollection.FindOne(ctx, filter)
	var dbCar entity.Car
	err = carResult.Decode(&dbCar)
	if err != nil {
		fmt.Println("Error:", err)
		return &autoserror.CustomError{
			Message: fmt.Sprintf("No car with id %s", carId),
			Status:  http.StatusNotFound,
		}
	}
	if dbCar.UserId != userId {
		return &autoserror.CustomError{
			Message: fmt.Sprintf("No delete access for car %s", carId),
			Status:  http.StatusForbidden,
		}
	}

	filter = bson.D{
		{Key: "_id", Value: carObjId},
		{Key: "userId", Value: userObjId},
	}
	_, err = carsCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func getCarForId(carId string, userId string) (dbCar entity.Car, err error) {
	ctx, client, cancel, err := DbSetup()
	if err != nil {
		return dbCar, err
	}
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return dbCar, err
	}
	defer client.Disconnect(ctx)

	autosDatabase := client.Database("autos")
	carsCollection := autosDatabase.Collection("cars")
	carObjId, err := primitive.ObjectIDFromHex(carId)
	if err != nil {
		return dbCar, err
	}
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return dbCar, err
	}

	filter := bson.D{
		{Key: "_id", Value: carObjId},
		{Key: "userId", Value: userObjId},
	}

	carResult := carsCollection.FindOne(ctx, filter)
	err = carResult.Decode(&dbCar)
	return dbCar, nil
}
