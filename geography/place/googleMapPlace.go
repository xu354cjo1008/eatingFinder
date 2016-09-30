package nearPlace

import (
	"context"
	"errors"
	"fmt"
	"googlemaps.github.io/maps"
	"time"
)

type gMapNearbySearchBase struct {
	client *maps.Client
	resp   map[int]*maps.PlacesSearchResponse
	req    *maps.NearbySearchRequest
}

func (base *gMapNearbySearchBase) requireBy(lat float64, lng float64, rad uint, lan string) (err error) {
	base.req.Location = &maps.LatLng{Lat: lat, Lng: lng}
	base.req.Radius = rad
	base.req.Language = lan
	return
}
func (base *gMapNearbySearchBase) requireTo() (err error) {
	page := 1
	for {
		resp, err := base.client.NearbySearch(context.Background(), base.req)
		if err != nil {
			break
		}

		base.resp[page] = &resp
		if base.resp[page].NextPageToken != "" {
			time.Sleep(3000 * time.Millisecond)
			base.req.PageToken = base.resp[page].NextPageToken
			page++
		} else {
			break
		}
	}
	return
}
func (base *gMapNearbySearchBase) parsing() (res map[int]map[int]map[string]interface{}, err error) {
	res = make(map[int]map[int]map[string]interface{})
	allRes := base.resp
	for page, paper := range allRes {
		if _, ok := res[page]; !ok {
			res[page] = make(map[int]map[string]interface{})
		}
		for rank, value := range paper.Results {
			if _, ok := res[page][rank]; !ok {
				res[page][rank] = make(map[string]interface{})
			}
			res[page][rank]["Location"] = fmt.Sprintf("lat: %f, lng: %f", value.Geometry.Location.Lat, value.Geometry.Location.Lng) // float64, float64
			res[page][rank]["name"] = value.Name                                                                                    // string
			res[page][rank]["open_now"] = value.OpeningHours.OpenNow                                                                // bool
			res[page][rank]["place_id"] = value.PlaceID                                                                             // string
			res[page][rank]["rating"] = value.Rating                                                                                // float32
			res[page][rank]["vicinity"] = value.Vicinity                                                                            // string
		}
	}
	return
}

func Init_gMapNS(apiKey string, clientID string, signature string, lat float64, lng float64, rad uint, keyword string, lang string, minPrice string, maxPrice string, name string, opennow bool, rankBy string, types string, pageToken string) (base *gMapNearbySearchBase, err error) {
	// Initialize
	base = &gMapNearbySearchBase{
		client: nil,
		resp:   make(map[int]*maps.PlacesSearchResponse),
		req:    new(maps.NearbySearchRequest),
	}
	// Client
	if apiKey != "" {
		base.client, err = maps.NewClient(maps.WithAPIKey(apiKey))
	} else if clientID != "" || signature != "" {
		base.client, err = maps.NewClient(maps.WithClientIDAndSignature(clientID, signature))
	} else {
		err = errors.New("Please specify an API Key, or Client ID and Signature.")
	}
	// Request
	base.req.Location = &maps.LatLng{Lat: lat, Lng: lng}
	base.req.Radius = rad
	base.req.Keyword = keyword
	base.req.Language = lang
	base.req.MinPrice, err = parsePriceLevel(minPrice)
	if err != nil {
		return
	}
	base.req.MaxPrice, err = parsePriceLevel(maxPrice)
	if err != nil {
		return
	}
	base.req.Name = name
	base.req.OpenNow = opennow
	base.req.RankBy, err = parseRankBy(rankBy)
	if err != nil {
		return
	}
	base.req.Type, err = parsePlaceType(types)
	if err != nil {
		return
	}
	base.req.PageToken = ""
	return
}

func parsePriceLevel(priceLevel string) (level maps.PriceLevel, err error) {
	switch priceLevel {
	case "0":
		return maps.PriceLevelFree, nil
	case "1":
		return maps.PriceLevelInexpensive, nil
	case "2":
		return maps.PriceLevelModerate, nil
	case "3":
		return maps.PriceLevelExpensive, nil
	case "4":
		return maps.PriceLevelVeryExpensive, nil
	default:
		return "", nil
	}
	return "", errors.New(fmt.Sprintf("Not handle price level : '%s'", priceLevel))
}
func parseRankBy(rankBy string) (res maps.RankBy, err error) {
	switch rankBy {
	case "prominence":
		res = maps.RankByProminence
	case "distance":
		res = maps.RankByDistance
	case "":
		res = maps.RankByProminence
	default:
		err = errors.New(fmt.Sprintf("Unknown rank by: \"%v\"", rankBy))
	}
	return
}
func parsePlaceType(placeType string) (res maps.PlaceType, err error) {
	if placeType != "" {
		res, err = maps.ParsePlaceType(placeType)
		if err != nil {
			err = errors.New(fmt.Sprintf("Unknown place type \"%v\"", placeType))
		}
	}
	return
}
