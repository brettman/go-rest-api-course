package http

import (
	"encoding/json"
	"fmt"
	"github.com/brettman/go-rest-api-course/internal/comment"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

// SetupRoutes - sets up all the routes for the app
func (h *Handler) SetupRoutes() {
	log.Info("Setting up routes.")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/comment", h.GetAllComments).Methods("GET")
	h.Router.HandleFunc("/api/comment", h.PostComment).Methods("POST")

	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", h.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", h.DeleteComment).Methods("DELETE")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := sendOkResponse(w, Response{Message: "I am alive"}); err != nil {
			panic(err)
		}
	})
}

// GetComment - retrieve a comment by id
func(h *Handler) GetComment(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10,64)

	if err !=nil{
		sendErrorResponse(w, "unable to parse UINT from ID", err)
		return
	}
	comment, err := h.Service.GetComment(uint(i))
	if err != nil {
		sendErrorResponse(w, "Error retriving comment by id", err)
		return
	}

	if err = sendOkResponse(w, comment); err !=nil{
		panic(err)
	}
}

// GetAllComments - retrieve all comments from db
func(h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request){
	comments, err := h.Service.GetAllComments()
	if err != nil{
		sendErrorResponse(w, "Error retrieving all comments", err)
		return
	}
	if err = sendOkResponse(w, comments); err !=nil{
		panic(err)
	}
}

// PostComments - retrieve all comments from db
func(h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var comment comment.Comment
	if err:= json.NewDecoder(r.Body).Decode(&comment); err != nil{
		sendErrorResponse(w, "Failed to decode request body", err)
		return
	}
	comment, err := h.Service.PostComment(comment)
	if err != nil{
		sendErrorResponse(w, "Failed to create comment", err)
		return
	}

	if err = sendOkResponse(w, comment); err !=nil{
		panic(err)
	}
}

// UpdateComments - update comment by id with new comment
func(h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10,64)

	if err !=nil{
		sendErrorResponse(w, "unable to parse UINT from ID", err)
		return
	}

	var comment comment.Comment
	if err:= json.NewDecoder(r.Body).Decode(&comment); err != nil{
		sendErrorResponse(w, "Failed to decode request body", err)
		return
	}

	comment, err = h.Service.UpdateComment(uint(i), comment)

	if err != nil{
		sendErrorResponse(w, "Failed to update comment", err)
		return
	}

	if err = sendOkResponse(w, comment); err !=nil{
		panic(err)
	}
}

// DeleteComments - delete comment by id
func(h *Handler) DeleteComment(w  http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10,64)

	if err !=nil{
		sendErrorResponse(w, "unable to parse UINT from ID", err)
		return
	}

	err = h.Service.DeleteComment(uint(i))
	if err !=nil{
		sendErrorResponse(w, fmt.Sprintf("Failed to delete comment with id %s", id), err)
		return
	}

	if err = sendOkResponse(w, &Response{Message: "Successfully deleted comment."}); err !=nil{
		panic(err)
	}
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