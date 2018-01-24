package msg

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
)

type MessageType int

const (
	Json MessageType = iota
	Xml
)

type ResponseBody struct {
	ContentType MessageType
	Data        interface{}
}

func (body *ResponseBody) Execute(writer http.ResponseWriter) error {
	var data []byte
	var err error
	switch body.ContentType {
	case Json:
		data, err = json.Marshal(body.Data)
	case Xml:
		data, err = xml.Marshal(body.Data)
	default:
		data, err = nil, errors.New("do not support content-type ")
	}
	if nil != err {
		return err
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
	return nil
}
