package saveusers

import (
	"bytes"
	"encoding/gob"
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
		return buff, err
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
		log.Println("SAVEUSERS.Load: Nothing to load")
		return &m, nil
	}
	dat, err := ioutil.ReadFile(c.File)
	if err != nil {
		return &m, err
	}
	buff := bytes.NewBuffer(dat)
	log.Println("SAVEUSERS.Load: Loaded users")
	return decode(buff)

}

// Save writes encoded registered users to file
func (c *EncConfig) Save(m *map[string]int64) error {
	buff, err := encode(m)
	if err != nil {
		return err
	}
	log.Println("SAVEUSERS.Save: users saved")
	return ioutil.WriteFile(c.File, buff.Bytes(), 0600)
}
