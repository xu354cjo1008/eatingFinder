package nearPlace

import ()

var placeUrl = map[string]string{
	"nearbySearch": "https://maps.googleapis.com/maps/api/place/nearbysearch",
}

type nearbySearch struct {
	// Essential parameter
	Key      paraFormat
	Location paraFormat // ${${latitude}, ${longtitude}}
	Radius   paraFormat // Unit: meter, maximum = 50,000 meter
	// Optional parameter
	Keyword       paraFormat
	Language      paraFormat // en, zh-TW
	Minprice      paraFormat // 0 ~ 4
	Maxprice      paraFormat // 0 ~ 4
	Name          paraFormat // split by " "
	Opennow       paraFormat // if enable only return the restaurant now is open
	Rankby        paraFormat // prominence(default), distance
	Type          paraFormat // please refer google api doc
	Pagetoken     paraFormat // search again and show the following 20 items of last result
	Zagatselected paraFormat // for Zagat
}

// level 4
type paraFormat struct {
	para  string
	value string
}

func addParameter(str string, parameter paraFormat) (res string) {
	if parameter.value != NSP {
		if str[len(str)-1] == '?' {
			res = str + parameter.para + "=" + parameter.value
		} else {
			res = str + "&" + parameter.para + "=" + parameter.value
		}
	} else {
		res = str
	}
	return
}
