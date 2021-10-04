package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"src/http/pkg/graphQL/graph/generated"
	"src/http/pkg/graphQL/graph/model"
	mymodel "src/http/pkg/model"

	"github.com/google/uuid"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	var p mymodel.Post
	p.Author = input.Author
	p.Message = input.Message

	newPost, err := r.service.CreateId(&p)
	if err != nil {
		return nil, errors.New("Internal problems")
	}

	idStr := newPost.Id.String()

	var post model.Post
	post.ID = idStr
	post.Author = newPost.Author
	post.Message = newPost.Message

	return &post, nil
}

func (r *mutationResolver) UpdatePost(ctx context.Context, input model.UpdatePost) (*model.Post, error) {
	id := input.ID
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("uuid parse problem")
	}
	var p mymodel.Post
	p.Id = idUUID
	p.Author = *input.Author
	p.Message = *input.Message

	res, err := r.service.UpdateId(&p)
	if err != nil {
		return nil, errors.New("services problem")
	}

	var post model.Post
	post.ID = id
	post.Author = res.Author
	post.Message = res.Message

	return &post, nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (*string, error) {
	var res string
	err := r.service.DeleteId(id)
	if err != nil {
		res = "post can't be deleted"
		return &res, errors.New("storage problem")
	}
	res = "post have been deleted"
	return &res, nil
}

func (r *queryResolver) GetPosts(ctx context.Context) ([]*model.Post, error) {
	var posts []*model.Post

	res := r.service.GetALL()

	for _, r := range *res {
		var p mymodel.Post
		p = r
		idStr := p.Id.String()
		var post model.Post
		post.ID = idStr
		post.Author = p.Author
		post.Message = p.Message
		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *queryResolver) GetPost(ctx context.Context, id string) (*model.Post, error) {
	res, err := r.service.GetId(id)
	if err != nil {
		return nil, errors.New("services problem")
	}
	idUUID := res.Id.String()
	post := model.Post{
		ID:      idUUID,
		Author:  res.Author,
		Message: res.Message,
	}

	return &post, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
