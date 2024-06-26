package structs

type WeatherData struct {
	Coord      coordinates `json:"coord"`
	Base       string      `json:"base"`
	Main       main        `json:"main"`
	Visibility int         `json:"visibility"`
	Wind       wind        `json:"wind"`
	Clouds     clouds      `json:"clouds"`
	Dt         int         `json:"dt"`
	Sys        sys         `json:"sys"`
	Timezone   int         `json:"timezone"`
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Cod        int         `json:"cod"`
}

type coordinates struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

type wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type clouds struct {
	All int `json:"all"`
}

type sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

type City struct {
	Name string `json:"city"`
}

type CityData []struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}
