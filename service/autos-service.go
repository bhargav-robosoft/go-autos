package service

import (
	"autos/db"
	"autos/entity"
)

type AutosService interface {
	ReadCars(userId string) (car []entity.Car, err error)
	CreateCar(car entity.NewCarRequest, userId string) (id string, err error)
	UpdateCar(car entity.ModifyCarRequest, userId string) (updatedCarId string, err error)
	DeleteCar(carId string, userId string) (err error)
}

type autosService struct{}

func NewAutosService() AutosService {
	return &autosService{}
}

func (service *autosService) ReadCars(userId string) ([]entity.Car, error) {
	cars, err := db.ReadCars()
	if err != nil {
		return cars, err
	}

	for i := range cars {
		if cars[i].UserId != "" && cars[i].UserId == userId {
			cars[i].IsAdmin = true
		}
	}

	return cars, nil
}

func (service *autosService) CreateCar(car entity.NewCarRequest, userId string) (string, error) {
	return db.CreateCar(car, userId)
}

func (service *autosService) UpdateCar(car entity.ModifyCarRequest, userId string) (string, error) {
	return db.UpdateCar(car, userId)
}

func (service *autosService) DeleteCar(carId string, userId string) error {
	return db.DeleteCar(carId, userId)
}
