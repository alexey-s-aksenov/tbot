package joke

import (
	"errors"
	"html"
	"io/ioutil"
	"net/http"
	"strings"
)

const site string = "http://nextjoke.net/Random"

type nextjoke struct {
	text string
}

// GetJoke получает html код по ссылке, вырезаем нужный кусок, конвертируем html спецсимволы
func (j *nextjoke) GetJoke() (string, error) {
	site := "http://nextjoke.net/Random"
	resp, err := http.Get(site)
	if err != nil {
		return "", errors.New("joke.go: Error in getJoke: " + err.Error())
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
