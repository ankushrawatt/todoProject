package utils

import (
	"encoding/json"
	"net/http"
)

//Encode all data  into JSON format
func Encoder(writer http.ResponseWriter, i interface{}) {
	err := json.NewEncoder(writer).Encode(i)
	CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
	return
}
