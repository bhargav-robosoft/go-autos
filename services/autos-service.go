package service

import "autos/entity"

type AutosService interface {
	GetAutos() ([]entity.Car, error)
	InsertAuto(entity.Car) (string, error)
}

type autosService struct{}

func New() AutosService {
	return &autosService{}
}

func (service *autosService) GetAutos() ([]entity.Car, error) {
	return []entity.Car{}, nil
}

func (service *autosService) InsertAuto(car entity.Car) (string, error) {
	return "", nil
}
