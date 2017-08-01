package meander

import "strings"

type j struct {
	Name string
	PlaceTypes []string
}

func (j j) Public() interface{} {
	return map[string]interface{}{
		"name": j.Name,
		"journey": strings.Join(j.PlaceTypes, "|"),
	}
}

var Journeys = []interface{}{
	j{Name: "Romantic", PlaceTypes: []string{"park", "bar", "movie_theater", "restaurant", "florist", "taxi_stand"}},
	j{Name: "Shopping", PlaceTypes: []string{"department_store", "cafe", "clothing_store", "jewlery_store", "shoe_store"}},
	j{Name: "Night out", PlaceTypes: []string{"bar", "cafe", "casino", "food", "night_club"}},
	j{Name: "Culture", PlaceTypes: []string{"museum", "cafe", "library", "art_gallery"}},
	j{Name: "Pamper", PlaceTypes: []string{"hair_care", "beauty_salon", "cafe", "spa"}},
}
