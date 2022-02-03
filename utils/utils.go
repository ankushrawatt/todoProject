package utils

//func Decode(writer http.ResponseWriter){
//	err:=json.NewEncoder(writer).Encode()
//	CheckError(err)
//}

func CheckError(err error) {
	if err != nil {

		panic(err)
	}
	return
}
