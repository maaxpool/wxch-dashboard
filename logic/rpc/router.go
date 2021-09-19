package rpc

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"reflect"
)

func registerRouter() http.Handler {
	r := mux.NewRouter()

	r.Handle("/rpc/partner", handlers.MethodHandler{
		"GET":  http.HandlerFunc(simpleUrlQueryWrap(reflect.TypeOf(getPartnerListRequest{}), getPartnerListHandler)),
		"POST": http.HandlerFunc(simpleJsonBodyWrap(reflect.TypeOf(createPartnerRequest{}), createPartnerHandler)),
	})

	r.Handle("/rpc/transaction", handlers.MethodHandler{
		"GET": http.HandlerFunc(simpleUrlQueryWrap(reflect.TypeOf(getTransactionListRequest{}), getTransactionListHandler)),
	})

	return handlers.LoggingHandler(os.Stdout, r)
}
