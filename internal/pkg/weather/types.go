package weather

// Weather struct
type Weather struct {
	Key string
}

// CurWeather defines JSON structure for current weather response
type CurWeather struct {
	Location *Location `json:"location"`
	Current  *Current  `json:"current"`
}

// ForeWeather defines JSON structure for forecast response
type ForeWeather struct {
}

// Location sub struct for CurWeather
type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Tzid           string  `json:"tz_id"`
	LocaltimeEpoch int64   `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

// Current sub struct fot CurWeather
type Current struct {
	LastUpdatedEpoch int64      `json:"last_updated_epoch"`
	LastUpdated      string     `json:"last_updated"`
	TempC            float64    `json:"temp_c"`
	TempF            float64    `json:"temp_f"`
	IsDay            int64      `json:"is_day"`
	Condition        *Condition `json:"condition"`
	WindMph          float64    `json:"wind_mph"`
	WindKph          float64    `json:"wind_kph"`
	WindDegree       int64      `json:"wind_degree"`
	WindDir          string     `json:"wind_dir"`
	PressureMb       float64    `json:"pressure_mb"`
	PressureIn       float64    `json:"pressure_in"`
	PrecipMm         float64    `json:"precip_mm"`
	PrecipIn         float64    `json:"precip_in"`
	Humidity         int64      `json:"humidity"`
	Cloud            int64      `json:"cloud"`
	FeelslikeC       float64    `json:"feelslike_c"`
	FeelslikeF       float64    `json:"feelslike_f"`
	VisKm            float64    `json:"vis_km"`
	VisMiles         float64    `json:"vis_miles"`
}

// Condition sub struct for Current
type Condition struct {
	Text string
	Icon string
	Code int64
}
