package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Onyekachukwu-Nweke/piko-blog/backend/internal/comment"
	"github.com/Onyekachukwu-Nweke/piko-blog/backend/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)


type CommentService interface {
	PostComment(context.Context, comment.Comment) (comment.Comment, error)
	GetComment(ctx context.Context, ID string) (comment.Comment, error)
	UpdateComment(ctx context.Context, ID string, newCmt comment.Comment) (comment.Comment, error)
	DeleteComment(ctx context.Context, ID string) error
}

type PostCommentRequest struct {
	PostID string `json:"post_id" validate:"required"`
	UserID string `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func convertPostCommentRequestToComment(c PostCommentRequest) comment.Comment {
	return comment.Comment{
		PostID: c.PostID,
		UserID: c.UserID,
		Content: c.Content,
	}
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	postID := vars["id"]
	userID, err := utils.GetUserIDFromContext(r) // This function should extract the user ID from the request context, set during authentication

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var cmt PostCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		return
	}

	cmt.PostID = postID
	cmt.UserID = userID

	// Debugging
	// fmt.Println(cmt)

	validate := validator.New()
	err = validate.Struct(cmt)
	if err != nil {
		http.Error(w, "not a valid comment", http.StatusBadRequest)
		return
	}

	convertedComment := convertPostCommentRequestToComment(cmt)
	postedCmt, err := h.CommentService.PostComment(r.Context(), convertedComment)
	if err != nil {
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(postedCmt); err != nil {
		panic(err)
	}
}

func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmt, err := h.CommentService.GetComment(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cmt); err != nil {
		panic(err)
	}
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Message: "Not Authorized"})
		return
	}

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var cmt comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		return
	}

	cmt.UserID = userID

	// Use the centralized authorization service
	if !h.AuthorizationService.IsUserAuthorized(r.Context(), userID, id, "comment") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(Response{Message: "Forbidden"})
		return
	}

	updatedCmt, err := h.CommentService.UpdateComment(r.Context(), id, cmt)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(updatedCmt); err != nil {
		panic(err)
	}
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Message: "Not Authorized"})
		return
	}

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Use the centralized authorization service
	if !h.AuthorizationService.IsUserAuthorized(r.Context(), userID, id, "comment") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(Response{Message: "Forbidden"})
		return
	}

	err = h.CommentService.DeleteComment(r.Context(), id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(Response{Message: "Successfully deleted"}); err != nil {
		panic(err)
	}
}