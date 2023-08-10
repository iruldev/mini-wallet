package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/iruldev/mini-wallet/src/config"

	muxHandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/iruldev/mini-wallet/engine/rest/routes"
	"github.com/spf13/viper"
)

func main() {
	r := mux.NewRouter()
	routes.AppRoutes(r)

	fmt.Printf("REST server listening at %v", viper.GetString("REST_PORT"))
	err := http.ListenAndServe(":"+viper.GetString("REST_PORT"), muxHandler.CORS(
		muxHandler.AllowedHeaders([]string{"Content-Type", "Authorization", "Accept"}),
		muxHandler.AllowedMethods([]string{http.MethodGet, http.MethodOptions, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}),
		muxHandler.AllowedOrigins([]string{"*"}),
	)(muxHandler.CompressHandler(muxHandler.LoggingHandler(os.Stdout, r))))
	if err != nil {
		panic(err)
	}
}
