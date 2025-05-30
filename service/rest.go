package service

import (
	"encoding/json"
	"fmt"
	"gitlab.com/devskiller-tasks/rest-api-blog-golang/model"
	"gitlab.com/devskiller-tasks/rest-api-blog-golang/repository"
	"net/http"
	"strconv"
)

type RestApiService struct {
	postRepository    *repository.PostRepository
	commentRepository *repository.CommentRepository
}

type AckJsonResponse struct {
	Message string
	Status  int
}

func NewRestApiService() RestApiService {
	return RestApiService{postRepository: repository.NewPostRepository(), commentRepository: repository.NewCommentRepository()}
}

func (svc *RestApiService) ServeContent(port int) error {
	http.HandleFunc("POST /api/posts", handleAddPost(svc))
	http.HandleFunc("GET /api/posts/{postId}", handleGetPostByPostId(svc))
	http.HandleFunc("POST /api/comments", handleAddComment(svc))
	http.HandleFunc("GET /api/comments", handleGetCommentsByPostId(svc))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handleAddPost(svc *RestApiService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var post model.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}
		if err := svc.postRepository.Insert(post); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(&AckJsonResponse{Message: fmt.Sprintf("Post with id: %d successfully added", post.Id), Status: http.StatusOK})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleGetPostByPostId(svc *RestApiService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Example: GET /api/posts/42

		// Every response should have the Content-Type=application/json header set.
		w.Header().Set("Content-Type", "application/json")

		// If an invalid ID is given, the response should be in the format of `AckJsonResponse` with a status of 400 and a message:
		// { "Message": "Wrong id path variable: PATH_VARIABLE", "Status": 400 }
		// The HTTP response code should also be set to 400.
		var postId uint64
		_, err := fmt.Sscanf(r.URL.Path, "/api/posts/%d", &postId)
		if err != nil {
			resp := AckJsonResponse{
				Message: fmt.Sprintf("Wrong id path variable: %s", r.URL.Path),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// If the given postID does not exist, the response should be in the format of `AckJsonResponse` with a status of 404 and a message:
		// { "Message": "Post with id: [POST_ID] does not exist", "Status": 404 }
		// The HTTP response code should also be set to 404.
		post, err := svc.postRepository.GetById(postId)
		if err != nil {
			resp := AckJsonResponse{
				Message: fmt.Sprintf("Post with id: %d does not exist", postId),
				Status:  http.StatusNotFound,
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// If the post with the given ID exists, the response should be a valid JSON representation of the post entity:
		// { "Id": 2, "Title": "test title", "Content": "this is a post content", "CreationDate": "1970-01-01T03:46:40+01:00" }
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
}

func handleGetCommentsByPostId(svc *RestApiService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Example: GET /api/comments?postId=4

		// Every response should have the Content-Type=application/json header set.
		w.Header().Set("Content-Type", "application/json")

		// If an invalid ID is given, the response should be in the format of `AckJsonResponse` with a status of 400 and a message:
		// { "Message": "Wrong id path variable: PATH_VARIABLE", "Status": 400 }
		// The HTTP response code should also be set to 400.
		query := r.URL.Query()
		postIdStr := query.Get("postId")
		if postIdStr == "" {
			resp := AckJsonResponse{
				Message: "Wrong id path variable: postId is missing",
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// If a valid postId is provided, the response should be a JSON array of comments.
		// If there are no comments for the given postId, the response should be an empty list.
		postId, err := strconv.ParseUint(postIdStr, 10, 64)
		if err != nil {
			resp := AckJsonResponse{
				Message: fmt.Sprintf("Wrong id path variable: %s", postIdStr),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		comments := svc.commentRepository.GetAllByPostId(postId)

		// Example JSON response:
		// [
		//     {"Id": 1, "PostId": 101, "Comment": "comment1", "Author": "author5", "CreationDate": "1970-01-01T03:46:40+01:05"},
		//     {"Id": 3, "PostId": 101, "Comment": "comment2", "Author": "author4", "CreationDate": "1970-01-01T03:46:40+01:10"},
		//     {"Id": 5, "PostId": 101, "Comment": "comment3", "Author": "author13", "CreationDate": "1970-01-01T03:46:40+01:15"}
		// ]

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(comments)
	}
}

func handleAddComment(svc *RestApiService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Example:
		// POST /api/comments
		// { "Id": 1, "PostId": 101, "Comment": "comment1", "Author": "author1", "CreationDate": "1970-01-01T03:46:40+01:00" }

		// Every response should have the Content-Type=application/json header set.
		w.Header().Set("Content-Type", "application/json")

		// If invalid or incomplete data is posted, the response should be in the format of `AckJsonResponse` with a status code of 400 and a message:
		// { "Message": "Could not deserialize comment JSON payload", "Status": 400 }
		// Data is considered incomplete when the payload misses any member property of the model.
		// The HTTP response code should also be 400.
		// Example:
		// POST /api/comments
		// { "weird_payload": "weird value" }
		// Response:
		// { "Message": "Could not deserialize comment JSON payload", "Status": 400 }
		var comment model.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			resp := AckJsonResponse{
				Message: "Could not deserialize comment JSON payload",
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// If a comment with the given ID already exists in the database, the response should be in the format of `AckJsonResponse` with a status code of 400 and a message:
		// { "Message": "Comment with id: COMMENT_ID already exists", "Status": 400 }
		// Example:
		// POST /api/comments
		// { "Id": 30, "PostId": 23123, "Comment": "comment1", "Author": "author1", "CreationDate": "1970-01-01T03:46:40+01:00" }
		// Response:
		// { "Message": "Comment with id: 30 already exists", "Status": 400 }
		if comment.Id == 0 || comment.PostId == 0 || comment.Comment == "" || comment.Author == "" || comment.CreationDate.IsZero() {
			resp := AckJsonResponse{
				Message: "Could not deserialize comment JSON payload",
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// If the data is posted successfully, the response should be in the format of `AckJsonResponse` with a status code of 200 and a message:
		// { "Message": "Comment with id: COMMENT_ID successfully added", "Status": 200 }
		// Example:
		// POST /api/comments
		// { "Id": 123, "PostId": 663, "Comment": "this is a comment", "Author": "blogger", "CreationDate": "1970-01-01T03:46:40+01:00" }
		// Response:
		// { "Message": "Comment with id: 123 successfully added", "Status": 200 }
		if _, err := svc.commentRepository.GetById(comment.Id); err == nil {
			resp := AckJsonResponse{
				Message: fmt.Sprintf("Comment with id: %d already exists", comment.Id),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		if err := svc.commentRepository.Insert(comment); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := AckJsonResponse{
			Message: fmt.Sprintf("Comment with id: %d successfully added", comment.Id),
			Status:  http.StatusOK,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
