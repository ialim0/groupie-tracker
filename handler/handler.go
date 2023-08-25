package handler

import (
	"fmt"
	help "groupie-tracker/help"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var allocation help.AllLocationIndex
var alldate help.AllDatesIndex
var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index"] = template.Must(template.ParseFiles("./templates/index.html", "./templates/base.html"))
	templates["about"] = template.Must(template.ParseFiles("./templates/about.html", "./templates/base.html"))
	templates["dates"] = template.Must(template.ParseFiles("./templates/dates.html", "./templates/base.html"))
	templates["date"] = template.Must(template.ParseFiles("./templates/date.html", "./templates/base.html"))
	templates["error"] = template.Must(template.ParseFiles("./templates/error.html", "./templates/base.html"))
	templates["location"] = template.Must(template.ParseFiles("./templates/location.html", "./templates/base.html"))
	templates["locations"] = template.Must(template.ParseFiles("./templates/locations.html", "./templates/base.html"))
	templates["relation"] = template.Must(template.ParseFiles("./templates/relation.html", "./templates/base.html"))
}

func RenderTemplate(w http.ResponseWriter, name string, templateName string, viewModel interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "The template does not exist.", http.StatusInternalServerError)
	}
	err := tmpl.ExecuteTemplate(w, templateName, viewModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {

	paths := []string{
		"/",
		"/location/",
		"/date/",
		"/relation/",
		"/locations/",
		"/locations-artist/",
		"/dates/",
		"/dates-artist/",
		"/about/",
		"/contact/",
		"/artists",
	}

	if !help.IsMatch(r.URL.Path, paths) {

		w.WriteHeader(404)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "404",
			ErrorMessage: "The page you are looking for cannot be found. It might have been moved or doesn't exist. ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(404)
		RenderTemplate(w, "error", "base", context)
		return

	}

	// Fetch artists data
	artistsURL := "https://groupietrackers.herokuapp.com/api/artists"
	var artists []help.Artists
	err := help.FetchDataFromAPI(artistsURL, &artists)
	if err != nil {
		w.WriteHeader(500)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "500",
			ErrorMessage: "Erreur lors de la récupération des données de l'API ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(500)
		RenderTemplate(w, "error", "base", context)
		return
	}

	context := help.Context{
		Artists: artists,
	}

	RenderTemplate(w, "index", "base", context)
}

func About(w http.ResponseWriter, r *http.Request) {

	paths := []string{
		"/",
		"/location/",
		"/date/",
		"/relation/",
		"/locations/",
		"/locations-artist/",
		"/dates/",
		"/dates-artist/",
		"/about/",
		"/contact/",
	}

	if !help.IsMatch(r.URL.Path, paths) {

		w.WriteHeader(404)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "404",
			ErrorMessage: "The page you are looking for cannot be found. It might have been moved or doesn't exist. ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(404)
		RenderTemplate(w, "error", "base", context)
		return

	}
	context := help.Context{}
	RenderTemplate(w, "about", "base", context)
}

func Dates(w http.ResponseWriter, r *http.Request) {
	paths := []string{
		"/",
		"/location/",
		"/date/",
		"/relation/",
		"/locations/",
		"/locations-artist/",
		"/dates/",
		"/dates-artist/",
		"/about/",
		"/contact/",
	}

	if !help.IsMatch(r.URL.Path, paths) {

		w.WriteHeader(404)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "404",
			ErrorMessage: "The page you are looking for cannot be found. It might have been moved or doesn't exist. ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(404)
		RenderTemplate(w, "error", "base", context)
		return

	}
	// Fetch artists data
	artistsURL := "https://groupietrackers.herokuapp.com/api/artists"
	var artists []help.Artists
	err := help.FetchDataFromAPI(artistsURL, &artists)
	if err != nil {
		w.WriteHeader(500)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "500",
			ErrorMessage: "Erreur lors de la récupération des données de l'API ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(500)
		RenderTemplate(w, "error", "base", context)
		return
	}

	// Fetch alldate data
	err = help.FetchDataFromAPI("https://groupietrackers.herokuapp.com/api/dates", &alldate)
	if err != nil {
		w.WriteHeader(500)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "500",
			ErrorMessage: "Erreur lors de la récupération des données de l'API ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(500)
		RenderTemplate(w, "error", "base", context)
		return
	}

	// Process data and render template
	uniqueDates := make(map[string][]int)
	for _, entry := range alldate.AllDatesIndex {
		for _, date := range entry.Dates {
			uniqueDates[date] = append(uniqueDates[date], entry.ID)
		}
	}

	var dateList help.DateList
	for date, artistIDs := range uniqueDates {
		var tab []help.Artists
		dateList.ListOfDate = append(dateList.ListOfDate, date)
		for _, id := range artistIDs {
			tab = append(tab, artists[id-1])
		}
		dateList.ListOfArtist = append(dateList.ListOfArtist, tab)
	}

	trimDateStruct := help.DateList{
		ListOfDate:   help.TrimStart(dateList.ListOfDate),
		ListOfArtist: dateList.ListOfArtist,
	}

	context := help.Context{
		DateList: trimDateStruct,
	}

	RenderTemplate(w, "dates", "base", context)
}

func Date(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de la location depuis l'URL
	id := r.URL.Path[len("/date/"):]

	intid, _ := strconv.Atoi(id)

	if intid <= 0 || intid > 52 {
		w.WriteHeader(404)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "404",
			ErrorMessage: "The page you are looking for cannot be found. It might have been moved or doesn't exist. ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(404)
		RenderTemplate(w, "error", "base", context)
		return

	}

	// Fetch dates data
	datesURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%s", id)
	var datesData help.DateData
	err := help.FetchDataFromAPI(datesURL, &datesData)
	if err != nil {
		w.WriteHeader(500)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "500",
			ErrorMessage: "Erreur lors de la récupération des données de l'API ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(500)
		RenderTemplate(w, "error", "base", context)
		return
	}

	// Fetch artists data
	artistsURL := "https://groupietrackers.herokuapp.com/api/artists"
	var artists []help.Artists
	err = help.FetchDataFromAPI(artistsURL, &artists)
	if err != nil {
		w.WriteHeader(500)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "500",
			ErrorMessage: "Erreur lors de la récupération des données de l'API ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(500)
		RenderTemplate(w, "error", "base", context)
		return
	}

	// Process data and render template
	n := datesData.ID - 1

	artist := artists[n]
	trimeDateStruct := help.DateData{
		ID:    datesData.ID,
		Dates: help.TrimStart(datesData.Dates),
	}

	context := help.Context{
		DateData: trimeDateStruct,
		Image:    artist.Image,
		Name:     artist.Name,
	}

	RenderTemplate(w, "date", "base", context)
}

func Locations(w http.ResponseWriter, r *http.Request) {

	paths := []string{
		"/",
		"/location/",
		"/date/",
		"/relation/",
		"/locations/",
		"/locations-artist/",
		"/dates/",
		"/dates-artist/",
		"/about/",
		"/contact/",
	}

	if !help.IsMatch(r.URL.Path, paths) {

		w.WriteHeader(404)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "404",
			ErrorMessage: "The page you are looking for cannot be found. It might have been moved or doesn't exist. ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(404)
		RenderTemplate(w, "error", "base", context)
		return

	}
	// Fetch artists data
	artistsURL := "https://groupietrackers.herokuapp.com/api/artists"
	var artists []help.Artists
	err := help.FetchDataFromAPI(artistsURL, &artists)
	if err != nil {
		w.WriteHeader(500)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "500",
			ErrorMessage: "Erreur lors de la récupération des données de l'API ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(500)
		RenderTemplate(w, "error", "base", context)
		return
	}

	// Fetch allocation data
	err = help.FetchDataFromAPI("https://groupietrackers.herokuapp.com/api/locations", &allocation)
	if err != nil {
		w.WriteHeader(500)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "500",
			ErrorMessage: "Erreur lors de la récupération des données de l'API ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(500)
		RenderTemplate(w, "error", "base", context)
		return
	}

	// Process data and render template
	uniqueLocations := make(map[string][]int)
	for _, entry := range allocation.AllLocationIndex {
		for _, location := range entry.Locations {
			location = strings.ToLower(location)
			uniqueLocations[location] = append(uniqueLocations[location], entry.ID)
		}
	}

	var locationList help.LocationList
	for location, artistIDs := range uniqueLocations {
		var tab []help.Artists
		locationList.ListOfLocation = append(locationList.ListOfLocation, location)
		for _, id := range artistIDs {
			tab = append(tab, artists[id-1])
		}
		locationList.ListOfArtist = append(locationList.ListOfArtist, tab)
	}

	context := help.Context{
		LocationList: locationList,
	}

	RenderTemplate(w, "locations", "base", context)
}

func Location(w http.ResponseWriter, r *http.Request) {

	// Récupérer l'ID de la location depuis l'URL
	id := r.URL.Path[len("/location/"):]
	intid, _ := strconv.Atoi(id)
	if intid <= 0 || intid > 52 {
		w.WriteHeader(404)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "404",
			ErrorMessage: "The page you are looking for cannot be found. It might have been moved or doesn't exist. ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(404)
		RenderTemplate(w, "error", "base", context)
		//help.RenderTemplate(w, context, "templates/error.html")
		return

	}

	// Fetch location data
	locationURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%s", id)
	var locationData help.LocationData
	err := help.FetchDataFromAPI(locationURL, &locationData)
	if err != nil {
		w.WriteHeader(404)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "404",
			ErrorMessage: "The page you are looking for cannot be found. It might have been moved or doesn't exist. ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(404)
		RenderTemplate(w, "error", "base", context)
		//help.RenderTemplate(w, context, "templates/error.html")
		return

	}

	// Fetch artists data
	artistsURL := "https://groupietrackers.herokuapp.com/api/artists"
	var artists []help.Artists
	err = help.FetchDataFromAPI(artistsURL, &artists)
	if err != nil {
		w.WriteHeader(500)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "500",
			ErrorMessage: "Erreur lors de la récupération des données de l'API ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(500)
		RenderTemplate(w, "error", "base", context)
		//help.RenderTemplate(w, context, "templates/error.html")
		return

	}

	// Process data and render template
	n := locationData.ID - 1

	artist := artists[n]

	context := help.Context{
		Artists:      artists,
		LocationData: locationData,
		Image:        artist.Image,
		Name:         artist.Name,
	}
	RenderTemplate(w, "location", "base", context)
}

func Relation(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de la location depuis l'URL
	id := r.URL.Path[len("/relation/"):]
	intid, _ := strconv.Atoi(id)

	if intid <= 0 || intid > 52 {
		w.WriteHeader(404)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "404",
			ErrorMessage: "The page you are looking for cannot be found. It might have been moved or doesn't exist. ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(404)
		RenderTemplate(w, "error", "base", context)
		return

	}
	// Fetch relations data
	relationURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%s", id)
	var relationData help.RelationsData
	err := help.FetchDataFromAPI(relationURL, &relationData)

	var locationDates []help.LocationDate

	for location, dates := range relationData.DatesLocations {
		for _, date := range dates {
			locationDate := help.LocationDate{
				Location: location,
				Date:     date,
			}
			locationDates = append(locationDates, locationDate)
		}
	}

	if err != nil {
		w.WriteHeader(500)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "500",
			ErrorMessage: "Erreur lors de la récupération des données de l'API ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(500)
		RenderTemplate(w, "error", "base", context)
		return
	}

	// Fetch artists data
	artistsURL := "https://groupietrackers.herokuapp.com/api/artists"
	var artists []help.Artists
	err = help.FetchDataFromAPI(artistsURL, &artists)
	if err != nil {
		w.WriteHeader(500)
		ErrMessage := help.ErrorStruct{
			ErrorName:    "500",
			ErrorMessage: "Erreur lors de la récupération des données de l'API ",
		}
		context := help.Context{
			ErrorStruct: ErrMessage,
		}
		w.WriteHeader(500)

		RenderTemplate(w, "error", "base", context)
		return
	}

	// Process data and render template
	n := relationData.ID - 1
	artist := artists[n]

	context := help.Context{
		Artists:       artists,
		RelationData:  relationData,
		LocationDates: locationDates,
		Image:         artist.Image,
		Name:          artist.Name,
	}

	RenderTemplate(w, "relation", "base", context)
}
