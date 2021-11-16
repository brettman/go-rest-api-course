package http

import (
	"encoding/json"
	"errors"
	"github.com/brettman/go-rest-api-course/internal/comment"
	"github.com/gorilla/mux"
	"net/http"
	log "github.com/sirupsen/logrus"
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
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
		}
	}
}


// SetupRoutes - sets up all the routes for the app
func (h *Handler) SetupRoutes() {
	log.Info("Setting up routes.")
	h.Router = mux.NewRouter()
	h.Router.Use(LoggingMiddleware)

	h.Router.HandleFunc("/api/comment", h.GetAllComments).Methods("GET")
	h.Router.HandleFunc("/api/comment", BasicAuth(h.PostComment)).Methods("POST")

	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", BasicAuth(h.UpdateComment)).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", BasicAuth(h.DeleteComment)).Methods("DELETE")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := sendOkResponse(w, Response{Message: "I am alive"}); err != nil {
			panic(err)
		}
	})
}

func sendOkResponse(w http.ResponseWriter, resp interface{}) error{
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil{
		panic(err)
	}
}