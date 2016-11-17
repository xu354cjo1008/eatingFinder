package main

import "io"

const (
	ALG_HIGHEST_RATE   = iota
	ALG_HIGHEST_SELECT = iota
)

type algUserData struct {
	lat float64
	lng float64
}

type algInterface interface {
	findRestaurantList(int, algUserData, int)
}

type algorithm struct {
	algHandler algInterface
}

func (alg *algorithm) findRestaurant(lat float64, lng float64) {

}

func (alg *algorithm) findRestaurantList(mode int, userData algUserData, size int) {

	alg.algHandler.findRestaurantList(mode, userData, size)

}

func (alg *algorithm) selectRestaurant(index int) {

}

func NewAlgorithm(logFile io.Writer) *algorithm {

	alg := algorithm{
		algHandler: newCCAlgorithm(logFile),
	}
	return &alg
}
