package deluge

import (
	"bytes"
	"encoding/base64"
	"errors"

	deluge "github.com/brunoga/go-deluge"
)

// Deluge struct to provide with Host, Port and Password to connect to Deluge
type Deluge struct {
	Host     string
	Port     string
	Password string
}

// AddTorrent func to add new torrent
func AddTorrent(data []byte, fileName string, config *Deluge) error {
	url := "http://" + config.Host + ":" + config.Port + "/json"

	d, err := deluge.New(url, config.Password)
	if err != nil {
		return errors.New("deluge.go: Failed to connect to Deluge: " + err.Error())
	}

	var buffer bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buffer)
	encoder.Write(data)
	encoder.Close()

	options := map[string]interface{}{
		"add_paused": false,
	}

	_, err = d.CoreAddTorrentFile(fileName, buffer.String(), options)
	if err != nil {
		return errors.New("deluge.go: Error adding torrent file :" + err.Error())
	}

	return nil

}
