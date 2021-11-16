package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brettman/go-rest-api-course/internal/comment"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
	log "github.com/sirupsen/logrus"
	"strings"
)

// Handler - stores pointer to comment service
type Handler struct {
	Router *mux.Router
	Service *comment.Service
}

//Response - object to store responses
type Response struct {
	Message string
	Error string
}

// NewHandler - returns a pointer to a Handler
func NewHandler(service *comment.Service) *Handler{
	return &Handler{
		Service: service,
	}
}

// LoggingMiddleware - Log endpoint activity
func LoggingMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path": r.URL.Path,
			}).Info("Endpoint hit")
		next.ServeHTTP(w, r)
	})
}

// BasicAuth - - middlware function that will provide basic auth around specific endpoints
func BasicAuth(original func(w http.ResponseWriter, r *http.Request)) func (w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		log.Info("basic auth endpoint hit")
		user, pass, ok := r.BasicAuth()
		if user == "admin" && pass == "password" && ok {
			original(w,r)
		}else{
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}
	}
}

// JWTAuth - decorator functio for jwt validatoin of endpoints
func JWTAuth(original func(w http.ResponseWriter, r * http.Request)) func (w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		log.Info("jwt authentiactio hit")
		authHeader := r.Header["Authorization"]

		if authHeader == nil  {
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}
		if len(authHeader) == 0  {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}

		// Bearer jwt-token
		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}
		if validateToken(authHeaderParts[1]){
			original(w,r)
		} else {
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}
	}
}

// SetupRoutes - sets up all the routes for the app
func (h *Handler) SetupRoutes() {
	log.Info("Setting up routes.")
	h.Router = mux.NewRouter()
	h.Router.Use(LoggingMiddleware)

	h.Router.HandleFunc("/api/comment", h.GetAllComments).Methods("GET")
	h.Router.HandleFunc("/api/comment", JWTAuth(h.PostComment)).Methods("POST")

	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", JWTAuth(h.UpdateComment)).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", JWTAuth(h.DeleteComment)).Methods("DELETE")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := sendOkResponse(w, Response{Message: "I am alive"}); err != nil {
			panic(err)
		}
	})
}

func validateToken (accessToken string) bool {
	var mySigningKey = []byte("missionimpossible")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface {}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there has been an error")
		}
		return mySigningKey, nil
	})
	if err != nil{
		return false
	}
	return token.Valid
}

func sendOkResponse(w http.ResponseWriter, resp interface{}) error{
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil{
		log.Error(err)
	}
}