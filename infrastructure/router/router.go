package router

import "net/http"

func New() http.Handler {
	mux := http.NewServeMux()

	registerV1Routes(mux)

	return mux
}
