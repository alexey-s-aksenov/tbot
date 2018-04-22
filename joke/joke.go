package joke

import (
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// GetJoke получает html код по ссылке, вырезаем нужный кусок, конвертируем html спецсимволы
func GetJoke() (string, error) {
	site := "http://nextjoke.net/Random"
	resp, err := http.Get(site)
	if err != nil {
		log.Printf("joke.go: Error in getJoke func: %s", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	response := string(body)
	first := strings.Index(response, "<div class=\"joke-text-div\">")
	response = response[first:]
	first = strings.Index(response, "<h1>")
	response = response[first+4:]

	second := strings.Index(response, "</h1>")
	response = response[0:second]
	return html.UnescapeString(response), nil

}
