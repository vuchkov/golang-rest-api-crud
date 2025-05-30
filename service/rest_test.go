package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/devskiller-tasks/rest-api-blog-golang/model"
	"gitlab.com/devskiller-tasks/rest-api-blog-golang/repository"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var (
	testDate = time.Date(2018, time.September, 16, 12, 0, 0, 0, time.UTC)
)

func TestAddPost(t *testing.T) {
	tests := []struct {
		testName           string
		commentRepository  repository.CommentRepository
		postRepository     repository.PostRepository
		post               interface{}
		expectedHttpStatus int
		expectedResponse   interface{}
	}{
		{
			testName:           "testSuccessfullyAddPost",
			post:               model.Post{Id: 256, Title: "title", Content: "cntnt", CreationDate: time.Now()},
			commentRepository:  repository.CustomCommentRepository(make([]model.Comment, 0)),
			postRepository:     repository.CustomPostRepository(make([]model.Post, 0)),
			expectedHttpStatus: 200,
			expectedResponse:   AckJsonResponse{Message: "Post with id: 256 successfully added", Status: http.StatusOK},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// GIVEN
			data, _ := json.Marshal(tc.post)
			req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewReader(data))
			w := httptest.NewRecorder()
			svc := RestApiService{&tc.postRepository, &tc.commentRepository}

			// WHEN
			handleAddPost(&svc)(w, req)
			response := w.Result()
			body, _ := io.ReadAll(response.Body)
			var ackResponse AckJsonResponse
			err := json.Unmarshal(body, &ackResponse)
			if err != nil {
				t.Fail()
			}

			// THEN
			assert.Equal(t, tc.expectedHttpStatus, response.StatusCode)
			assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
			assert.Equal(t, tc.expectedResponse, ackResponse)
		})
	}
}

var validComments = []model.Comment{
	{Id: 123, PostId: 3, Comment: "abc", Author: "cool author", CreationDate: testDate},
	{Id: 321, PostId: 3, Comment: "def", Author: "cool author2", CreationDate: testDate},
	{Id: 543, PostId: 3, Comment: "ghi", Author: "cool author3", CreationDate: testDate},
}

func TestGetComments(t *testing.T) {
	tests := []struct {
		testName           string
		commentRepository  repository.CommentRepository
		postRepository     repository.PostRepository
		postId             int
		expectedHttpStatus int
		expectedResponse   interface{}
	}{
		{
			testName:           "testSuccessfullyGetComments",
			commentRepository:  repository.CustomCommentRepository(validComments),
			postRepository:     repository.CustomPostRepository(make([]model.Post, 0)),
			postId:             3,
			expectedHttpStatus: 200,
			expectedResponse:   validComments,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// GIVEN
			svc := RestApiService{commentRepository: &tc.commentRepository,
				postRepository: &tc.postRepository}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/comments?postId=%d", tc.postId), nil)
			w := httptest.NewRecorder()

			// WHEN
			handleGetCommentsByPostId(&svc)(w, req)
			response := w.Result()
			body, _ := io.ReadAll(response.Body)
			var commentsList []model.Comment
			err := json.Unmarshal(body, &commentsList)

			if err != nil {
				t.Fail()
			}

			// THEN
			assert.Equal(t, tc.expectedHttpStatus, response.StatusCode)
			assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
			assert.ElementsMatch(t, tc.expectedResponse, commentsList)
		})
	}
}

func TestAddComment(t *testing.T) {
	tests := []struct {
		testName           string
		commentRepository  repository.CommentRepository
		postRepository     repository.PostRepository
		comment            model.Comment
		expectedHttpStatus int
		expectedResponse   interface{}
	}{
		{
			testName:           "testSuccessfullyAddComment",
			commentRepository:  repository.CustomCommentRepository(make([]model.Comment, 0)),
			postRepository:     repository.CustomPostRepository(make([]model.Post, 0)),
			comment:            model.Comment{Id: 123, PostId: 3, Comment: "cool cmnt", Author: "cool auth", CreationDate: time.Now()},
			expectedHttpStatus: 200,
			expectedResponse:   AckJsonResponse{Message: "Comment with id: 123 successfully added", Status: http.StatusOK},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// GIVEN
			svc := RestApiService{commentRepository: &tc.commentRepository,
				postRepository: &tc.postRepository}

			data, _ := json.Marshal(&tc.comment)
			req := httptest.NewRequest(http.MethodPost, "/api/comments", bytes.NewReader(data))
			w := httptest.NewRecorder()

			// WHEN
			handleAddComment(&svc)(w, req)

			response := w.Result()
			body, _ := io.ReadAll(response.Body)
			var resp AckJsonResponse
			err := json.Unmarshal(body, &resp)

			if err != nil {
				t.Fail()
			}

			// THEN
			assert.Equal(t, tc.expectedHttpStatus, response.StatusCode)
			assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}

var validPost = model.Post{
	Id:           34,
	Title:        "happy post",
	Content:      "test content",
	CreationDate: testDate,
}

func TestGetPost(t *testing.T) {
	tests := []struct {
		testName           string
		commentRepository  repository.CommentRepository
		postRepository     repository.PostRepository
		postId             string
		expectedHttpStatus int
		expectedResponse   interface{}
	}{
		{
			testName:           "testSuccessfullyGetPost",
			commentRepository:  repository.CustomCommentRepository(make([]model.Comment, 0)),
			postRepository:     repository.CustomPostRepository([]model.Post{validPost}),
			postId:             strconv.Itoa(int(validPost.Id)),
			expectedHttpStatus: 200,
			expectedResponse:   validPost,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// GIVEN
			svc := RestApiService{commentRepository: &tc.commentRepository,
				postRepository: &tc.postRepository}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/posts/%s", tc.postId), nil)
			req.SetPathValue("postId", tc.postId)
			w := httptest.NewRecorder()

			// WHEN
			handleGetPostByPostId(&svc)(w, req)
			response := w.Result()
			body, _ := io.ReadAll(response.Body)
			var post model.Post
			err := json.Unmarshal(body, &post)

			if err != nil {
				t.Fail()
			}

			// THEN
			assert.Equal(t, tc.expectedHttpStatus, response.StatusCode)
			assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
			assert.Equal(t, tc.expectedResponse, post)
		})
	}
}
