package repository

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/devskiller-tasks/rest-api-blog-golang/model"
	"testing"
	"time"
)

var (
	comment1          = model.Comment{Id: 1, PostId: 101, Comment: "comment2", Author: "author2", CreationDate: time.Unix(10011, 0)}
	comment2          = model.Comment{Id: 2, PostId: 101, Comment: "comment2", Author: "author2", CreationDate: time.Unix(10011, 0)}
	comment3          = model.Comment{Id: 3, PostId: 100, Comment: "comment3", Author: "author3", CreationDate: time.Unix(10022, 0)}
	NonExistentPostId = uint64(10101010)
)

func TestSimpleGetAllByPostId(t *testing.T) {
	c := CommentRepository{}
	c.Insert(comment1)
	assert.ElementsMatch(t, c.GetAllByPostId(comment1.PostId), []model.Comment{comment1})
}

func TestGetAllByPostId(t *testing.T) {
	c := CommentRepository{}
	c.Insert(comment1)
	c.Insert(comment2)
	c.Insert(comment3)

	expectedResult := []model.Comment{comment1, comment2}
	result := c.GetAllByPostId(comment1.PostId)
	assert.ElementsMatch(t, expectedResult, result)
}

func TestNoComments(t *testing.T) {
	c := CommentRepository{}
	c.Insert(comment1)
	c.Insert(comment2)
	c.Insert(comment3)

	assert.ElementsMatch(t, c.GetAllByPostId(NonExistentPostId), make([]*model.Comment, 0))
}

func TestInsertExistingComment(t *testing.T) {
	c := CommentRepository{}
	if err := c.Insert(comment1); err != nil {
		t.Fail()
	}
	err := c.Insert(comment1)
	assert.EqualErrorf(t, err, "Comment with id: 1 already exists", "test failed because of wrong error msg: %+v", err)
}
