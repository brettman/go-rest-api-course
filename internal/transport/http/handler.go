package http

import (
	"encoding/json"
	"fmt"
	"github.com/brettman/go-rest-api-course/internal/comment"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Handler - stores pointer to comment service
type Handler struct {
	Router *mux.Router
	Service *comment.Service
}

//Response - object to store responses
type Response struct {
	Message string
}

// NewHandler - returns a pointer to a Handler
func NewHandler(service *comment.Service) *Handler{
	return &Handler{
		Service: service,
	}
}

// SetupRoutes - sets up all the routes for the app
func (h *Handler) SetupRoutes() {
	fmt.Println("Setting up routes.")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/comment", h.GetAllComments).Methods("GET")
	h.Router.HandleFunc("/api/comment", h.PostComment).Methods("POST")

	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", h.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", h.DeleteComment).Methods("DELETE")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{Message: "I am alive"}); err != nil {
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
		fmt.Fprintf(w, "unable to parse UINT from ID")
	}
	comment, err := h.Service.GetComment(uint(i))
	if err != nil {
		fmt.Fprintf(w, "Error retriving comment by id")
	}

	fmt.Fprintf(w, "%+v", comment)
}

// GetAllComments - retrieve all comments from db
func(h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request){
	comments, err := h.Service.GetAllComments()
	if err != nil{
		fmt.Fprintf(w, "Error retrieving all comments")
	}
	fmt.Fprintf(w, "%+v", comments)
}

// PostComments - retrieve all comments from db
func(h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	comment, err := h.Service.PostComment(comment.Comment{
		Slug: "/",
	})
	if err != nil{
		fmt.Fprintf(w, "Failed to create comment")
	}

	fmt.Fprintf(w, "%+v", comment)
}

// UpdateComments - update comment by id with new comment
func(h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10,64)

	if err !=nil{
		fmt.Fprintf(w, "unable to parse UINT from ID")
	}

	comment, err := h.Service.UpdateComment(uint(i), comment.Comment{
		Slug: "/new",
	})

	if err != nil{
		fmt.Fprintf(w, "Failed to update comment")
	}

	fmt.Fprintf(w, "%+v", comment)
}

// DeleteComments - delete comment by id
func(h *Handler) DeleteComment(w  http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10,64)

	if err !=nil{
		fmt.Fprintf(w, "unable to parse UINT from ID")
	}

	err = h.Service.DeleteComment(uint(i))
	if err !=nil{
		fmt.Fprintf(w, "Failed to delete comment with id %s", id)
	}

	fmt.Fprintf(w, "Sucessfully deleted comment")
}
