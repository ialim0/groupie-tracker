package help

type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type DateData struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type LocationData struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	DatesURL  string   `json:"dates"`
}
type RelationsData struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type LocationDate struct {
	Location string
	Date     string
}

type AllocationEntry struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

// Structure pour représenter l'index complet
type AllLocationIndex struct {
	AllLocationIndex []AllocationEntry `json:"index"`
}

type AllDatesEntry struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type AllDatesIndex struct {
	AllDatesIndex []AllDatesEntry `json:"index"`
}

type LocationList struct {
	ListOfLocation []string
	ListOfArtist   [][]Artists
}

type DateList struct {
	ListOfDate   []string
	ListOfArtist [][]Artists
}
type ErrorStruct struct {
	ErrorName    string
	ErrorMessage string
}

type Context struct {
	Artists       []Artists
	LocationData  LocationData
	DateData      DateData
	RelationData  RelationsData
	LocationDates []LocationDate
	LocationList  LocationList
	DateList      DateList
	ErrorStruct   ErrorStruct
	Image         string
	Name          string
}
