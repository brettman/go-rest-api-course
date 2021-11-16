package http

import (
	"encoding/json"
	"fmt"
	"github.com/brettman/go-rest-api-course/internal/comment"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

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

