package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"net/http"
	"reflect"
)

var validate = validator.New()
var urlDecode = schema.NewDecoder()

type SimpleWrapFunc func(r *http.Request) (resp interface{}, err error)
type SimpleInterfaceWrapFunc func(typ interface{}, r *http.Request) (resp interface{}, err error)
type RespFunc func(writer http.ResponseWriter, request *http.Request)

type successJson struct {
	Success bool        `json:"success"`
	Msg     interface{} `json:"msg"`
}

type failJson struct {
	Success bool   `json:"success"`
	ErrCode uint32 `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

func simpleWrap(api SimpleWrapFunc) RespFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var successJson = &successJson{Success: true}
		var failJson = &failJson{Success: false}

		jsonWriter := json.NewEncoder(w)
		successJson.Msg, err = api(r)
		if err != nil {
			if httpErr, ok := err.(*HTTPError); ok {
				failJson.ErrCode = httpErr.errorCode
				failJson.ErrMsg = httpErr.error.Error()
			} else {
				failJson.ErrCode = 500
				failJson.ErrMsg = err.Error()
			}

			if failJson.ErrCode >= 400 && failJson.ErrCode < 600 {
				w.WriteHeader(int(failJson.ErrCode))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			_ = jsonWriter.Encode(failJson)
		} else {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")

			_ = jsonWriter.Encode(successJson)
		}
	}
}

func simpleJsonBodyWrap(typ reflect.Type, api SimpleInterfaceWrapFunc) RespFunc {
	return simpleWrap(func(r *http.Request) (resp interface{}, err error) {
		if r.ContentLength <= 0 {
			return nil, fmt.Errorf("need a `%s` body", typ.String())
		}

		val := reflect.New(typ).Interface()

		err = json.NewDecoder(r.Body).Decode(&val)
		if err != nil {
			return nil, fmt.Errorf("decode `%s` body err: %v", typ.String(), err)
		}

		err = validate.Struct(val)
		if err != nil {
			return nil, err
		}

		return api(val, r)
	})
}

func simpleUrlQueryWrap(typ reflect.Type, api SimpleInterfaceWrapFunc) RespFunc {
	return simpleWrap(func(r *http.Request) (resp interface{}, err error) {
		val := reflect.New(typ).Interface()

		err = urlDecode.Decode(val, r.URL.Query())
		if err != nil {
			return nil, fmt.Errorf("decode `%s` query err: %v", typ.String(), err)
		}

		err = validate.Struct(val)
		if err != nil {
			return nil, err
		}

		return api(val, r)
	})
}
