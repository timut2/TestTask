package jsonutil

import (
	"encoding/json"
	"log"
	"net/http"
)

type Wrap map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Wrap, header http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')
	for key, value := range header {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1 MB

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
