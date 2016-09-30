package nearPlace

import (
	"fmt"
	"testing"
)

type para_rest struct {
	lat float64
	lng float64
	rad uint
	lan string
}

var test_case_rest = [...]para_rest{
	{25.027228, 121.522637, 500, "en"},
}

func TestNearRestaurants(t *testing.T) {
	for _, tCase := range test_case_rest {
		fmt.Println("=====> test case :", tCase)
		for name, _ := range priority {
			fmt.Println("===@@@@@=== source :", name)
			base, err := InitPlaceNearbySearch(name)
			if err != nil {
				t.Error(err)
			}
			data, err := base.GetNearRestaurants(tCase.lat, tCase.lng, tCase.rad, tCase.lan)
			if err != nil {
				t.Error(err)
			}
			for page, paper := range data {
				fmt.Println("----------Page ", page, "----------")
				for rank, info := range paper {
					fmt.Println("=====Rank=====:", rank)
					for field, value := range info {
						fmt.Println(field, ":", value)
					}
				}
			}
		}
	}
}
