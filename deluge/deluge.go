package deluge

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"

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
		log.Printf("de;uge.go: Failed to connect to Deluge: ", err)
		return err
	}

	var buffer bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buffer)
	encoder.Write(data)
	encoder.Close()

	options := map[string]interface{}{
		"add_paused": false,
	}

	id, err := d.CoreAddTorrentFile(fileName, buffer.String(), options)
	if err != nil {
		fmt.Println("deluge.go: Error adding torrent file :", err)
	}

	log.Printf("deluge.go: Added torrent with id: %s", id)

	return nil

}
