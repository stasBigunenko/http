package services

import (
	"encoding/json"

	"src/http/pkg/model"
)

func Response(message string) []byte {

	info := model.MyResponse{
		Msg: message,
	}

	res, _ := json.Marshal(info)
	return res

}
