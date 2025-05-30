package repository

import (
	"fmt"
	"gitlab.com/devskiller-tasks/rest-api-blog-golang/model"
)

type CommentRepository struct {
	repository []model.Comment
}

func NewCommentRepository() *CommentRepository {
	repo := CustomCommentRepository(make([]model.Comment, 0))
	return &repo
}

func CustomCommentRepository(mockStorage []model.Comment) CommentRepository {
	return CommentRepository{repository: mockStorage}
}

type CommentAlreadyExistsError struct {
	id uint64
}

func (e CommentAlreadyExistsError) Error() string {
	return fmt.Sprintf("Comment with id: %v already exists", e.id)
}

type CommentNotFoundError struct {
	id uint64
}

func (e CommentNotFoundError) Error() string {
	return fmt.Sprintf("Comment with id: %v does not exist", e.id)
}

func (c *CommentRepository) Insert(comment model.Comment) error {
	// Insert should insert a comment passed as an argument to the persistent in memory repository.
	// The method should return an error as an instance of `CommentAlreadyExistsError` struct
	// when a comment with given id already exists in the repository.
	_, err := c.GetById(comment.Id)
	if err == nil {
		return CommentAlreadyExistsError{id: comment.Id}
	}

	if _, ok := err.(CommentNotFoundError); !ok {
		return err
	}

	c.repository = append(c.repository, comment)
	return nil
}

func (c *CommentRepository) GetById(id uint64) (*model.Comment, error) {
	// GetById should return a comment from a repository that has a given id.
	// If there's no comment with given id, this function should return a (nil, CommentNotFoundError) pair
	// with CommentNotFound instance having id member variable set with id passed to this method.
	for i := range c.repository {
		if c.repository[i].Id == id {
			return &c.repository[i], nil
		}
	}

	return nil, CommentNotFoundError{id}
}

func (c *CommentRepository) GetAllByPostId(id uint64) []model.Comment {
	// GetAllByPostId should return a slice of all comments that have PostId member variable
	// equal to given id.
	// The method should return an empty slice when there are no comments with given id in the repository.
	var result []model.Comment

	for _, comment := range c.repository {
		if comment.PostId == id {
			result = append(result, comment)
		}
	}

	return result
}

type PostRepository struct {
	repository []model.Post
}

func CustomPostRepository(mockStorage []model.Post) PostRepository {
	return PostRepository{repository: mockStorage}
}

func NewPostRepository() *PostRepository {
	repo := CustomPostRepository(make([]model.Post, 0))
	return &repo
}

type PostAlreadyExistsError struct {
	id uint64
}

func (e PostAlreadyExistsError) Error() string {
	return fmt.Sprintf("Post with id: %v does not exist", e.id)
}

type PostNotFoundError struct {
	id uint64
}

func (e PostNotFoundError) Error() string {
	return fmt.Sprintf("Post with id: %v already exists", e.id)
}

func (c *PostRepository) Insert(post model.Post) error {
	// Insert should insert a post passed as an argument to the persistent in memory repository.
	// The method should return an error as an instance of `PostAlreadyExistsError` struct
	// when a post with given id already exists in the repository.
	_, err := c.GetById(post.Id)
	if err == nil {
		return PostAlreadyExistsError{id: post.Id}
	}

	if _, ok := err.(PostNotFoundError); !ok {
		return err
	}

	c.repository = append(c.repository, post)
	return nil
}

func (c *PostRepository) GetById(id uint64) (*model.Post, error) {
	// GetById should return a post from a repository that has a given id.
	// If there's no post with given id, this function should return a (nil, PostNotFoundError) pair
	// with PostNotFoundError instance having id member variable set with id passed to this method.
	for i := range c.repository {
		if c.repository[i].Id == id {
			return &c.repository[i], nil
		}
	}

	return nil, PostNotFoundError{id}
}
