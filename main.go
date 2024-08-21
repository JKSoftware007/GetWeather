package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Content format not specified.  Went with JSON since it's less verbose.
// Other can be added if needed.

// A better logger should be used to support various log levels, etc.

// https should be used

// Should add an error response body to let the end user know more about what failed

func main() {
	// Gracefully handle ctrl-c termination
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Ctrl-C detected.  Exiting...")
		os.Exit(1)
	}()

	router := mux.NewRouter()
	router.HandleFunc("/currentweather", currentWeatherHandler).Methods("POST").Schemes("http")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln("There's an error with the server,", err.Error())
	}
}

// currentWeatherHandler handles the post to get current weather
// It will return the forecast for the indicated longitude and latitude
func currentWeatherHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var location Coordinates
	if err := json.NewDecoder(request.Body).Decode(&location); err != nil {
		log.Default().Println(`Unable to process request, decode of coordinates failed: `, err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	pointsDetail, err := getLocationDetails(location)
	if err != nil {
		writer.WriteHeader(http.StatusFailedDependency)
		return
	}

	forecast, err := getForecast(pointsDetail)
	if err != nil {
		writer.WriteHeader(http.StatusFailedDependency)
		return
	}

	// Get the short forecast for today/tonight
	period := forecast.Properties.Periods[0]

	CharacterizeTemp(&period)

	if err := json.NewEncoder(writer).Encode(&period); err != nil {
		log.Default().Println(`Unable encode the response: `, err.Error())
		// Might be a better error code, but this works well enough for a quick exercise.
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

// CharacterizeTemp will add a descriptive label to the temperature: hot, cold, or moderate
func CharacterizeTemp(period *Period) {
	// Ideally code would check unit and use values correctly based on that.  This assumes Fahrenheit
	switch {
	case period.Temperature > 85:
		period.Characterization = "hot"
	case period.Temperature < 45:
		period.Characterization = "cold"
	default:
		period.Characterization = "moderate"
	}
}

// getForecast will get the forecast for the given location
func getForecast(locationDetails LocationDetails) (Forecast, error) {
	response, err := http.Get(locationDetails.Properties.Forecast)

	if err != nil {
		log.Println("Error retrieving forecast from national weather service: ", err.Error())
		return Forecast{}, err
	}

	var forecast Forecast
	if err := json.NewDecoder(response.Body).Decode(&forecast); err != nil {
		log.Println("Error processing forecast detail: ", err.Error())
		return Forecast{}, err
	}

	return forecast, nil
}

// getLocationDetails will retrieve the location details given the coordinates received
func getLocationDetails(location Coordinates) (LocationDetails, error) {
	// Should avoid the hard coded string.  Should pull from a config.
	// The endpoint takes a maximum of 4 decimal points, so we have to format the values we receive
	url := fmt.Sprintf("https://api.weather.gov/points/%.4f,%.4f", location.Latitude, location.Longitude)
	response, err := http.Get(url)

	if err != nil {
		log.Println("Error retrieving location from national weather service: ", err.Error())
		return LocationDetails{}, err
	}

	var locationDetails LocationDetails
	if err := json.NewDecoder(response.Body).Decode(&locationDetails); err != nil {
		log.Println("Error processing points detail: ", err.Error())
		return LocationDetails{}, err
	}

	return locationDetails, nil
}

// In Golang, the property must be capitalized to be public.  Since the JSON properties are lower case, we tell
// the JSON parser the field in the data will be lowercase.

// Coordinates are the input from the user we get the forecast for
type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// Location information received back from NWS
type LocationDetails struct {
	Properties LocationDetailProperties `json:"properties"`
}

// LocationDetailProperties contains the forecast link used to get the forecast from NWS
type LocationDetailProperties struct {
	Forecast string `json:"forecast"`
}

// Forecast is the response object from NWS
type Forecast struct {
	Properties ForecastProperties `json:"properties"`
}

// ForecastProperties contains the periods for the forecast from NWS
type ForecastProperties struct {
	Periods []Period `json:"periods"`
}

// Period the forecast period with the forecast details
type Period struct {
	Name                       string                   `json:"name"`
	Temperature                int                      `json:"temperature"`
	TemperatureUnit            string                   `json:"temperatureUnit"`
	Characterization           string                   `json:"characterization"`
	WindSpeed                  string                   `json:"windSpeed"`
	WindDirection              string                   `json:"windDirection"`
	ShortForecast              string                   `json:"shortForecast"`
	DetailedForecast           string                   `json:"detailedForecast"`
	ProbabilityOfPrecipitation PrecipitationProbability `json:"probabilityOfPrecipitation"`
}

// PrecipitationProabability a subobject of the period containing precipitation chance
type PrecipitationProbability struct {
	UnitCode string `json:"unitCode"`
	Value    int    `json:"value"`
}
