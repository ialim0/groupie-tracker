package help

import (
	"encoding/json"
	"net/http"
	"strings"
)

func FetchDataFromAPI(apiURL string, data interface{}) error {

	response, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func FindStringIndex(slice []string, target string) int {
	for i, s := range slice {
		if s == target {
			return i
		}
	}
	return -1
}

// Function to check if a given path matches any of the defined paths
func IsMatch(path string, tab []string) bool {

	for _, p := range tab {
		if path == p {
			return true
		}
	}

	return false
}

func TrimStart(tab []string) []string {
	for i, _ := range tab {
		tab[i] = strings.Trim(tab[i], "*")
	}
	return tab
}
