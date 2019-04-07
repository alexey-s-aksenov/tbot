package saveusers

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// EncConfig struct for encoding package to store filename
type EncConfig struct {
	File string
}

// Encode encodes given map into bytes.Buffer
func encode(m *map[string]int64) (bytes.Buffer, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(m)
	if err != nil {
		return buff, errors.New("saveusers.go: Failed to encode: " + err.Error())
	}
	return buff, nil
}

// Decoder decodec given bytes.Buffer into map[string]int64
func decode(b *bytes.Buffer) (*map[string]int64, error) {
	m := make(map[string]int64)
	dec := gob.NewDecoder(b)
	err := dec.Decode(&m)
	if err != nil {
		log.Fatal(err)
	}
	return &m, nil
}

// Load tries to load saved users from file if exists.
func (c *EncConfig) Load() (*map[string]int64, error) {
	m := make(map[string]int64)
	if _, err := os.Stat(c.File); os.IsNotExist(err) {
		// No file, return empty list
		return &m, nil
	}
	dat, err := ioutil.ReadFile(c.File)
	if err != nil {
		return &m, errors.New("saveusers.go: Failed to read file, error: " + err.Error())
	}
	buff := bytes.NewBuffer(dat)
	return decode(buff)

}

// Save writes encoded registered users to file
func (c *EncConfig) Save(m *map[string]int64) error {
	buff, err := encode(m)
	if err != nil {
		return errors.New("saveusers.go: Failed to save file, error: " + err.Error())
	}
	return ioutil.WriteFile(c.File, buff.Bytes(), 0600)
}
