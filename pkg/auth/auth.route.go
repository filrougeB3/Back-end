package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
)

var AuthRoute = func(router *mux.Router) {
	r := chi.NewRouter()

	// Routes non protégées par le middleware
	r.Post("/auth/register", CreateUser)
	r.Post("/auth/login", LoginUser)

	// Routes protégées par le middleware
	r.With(AuthMiddleware).Post("/auth/logout", LogoutUser)

	http.ListenAndServe(":8080", r)
}
