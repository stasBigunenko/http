package grpccontroller

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "src/http/api/proto"
	"src/http/pkg/model"
)

type gRPCClient struct {
	client pb.PostServiceClient
}

func New(store pb.PostServiceClient) gRPCClient {
	return gRPCClient{
		client: store,
	}
}

func (gc gRPCClient) Get(id uuid.UUID) (model.Post, error) {
	val := id.String()

	pp, err := gc.client.Get(context.Background(), &pb.PostId{
		Id: val,
	})
	if err != nil {
		return model.Post{}, status.Error(codes.NotFound, "post not found")
	}

	res, err := uuid.Parse(pp.Id)

	return model.Post{
		Id:      res,
		Author:  pp.Author,
		Message: pp.Message,
	}, nil
}

func (gc gRPCClient) GetAll() []model.Post {
	ap, _ := gc.client.GetAll(context.Background(), &emptypb.Empty{})

	posts := []model.Post{}

	for _, val := range ap.Allposts {
		res, err := uuid.Parse(val.Id)
		if err != nil {
			return nil
		}

		posts = append(posts, model.Post{
			Id:      res,
			Author:  val.Author,
			Message: val.Message,
		})
	}

	return posts
}

func (gc gRPCClient) Create(pp model.Post) (model.Post, error) {

	post, err := gc.client.Create(context.Background(), &pb.NewPost{
		Author:  pp.Author,
		Message: pp.Message,
	})
	if err != nil {
		return model.Post{}, status.Error(codes.Internal, "couldn't create post")
	}

	val, err := uuid.Parse(post.Id)
	if err != nil {
		return model.Post{}, status.Error(codes.Internal, "couldn't parse id")
	}

	return model.Post{
		Id:      val,
		Author:  post.Author,
		Message: post.Message,
	}, nil

}

func (gc gRPCClient) Update(pp model.Post) (model.Post, error) {

	val := pp.Id.String()

	newPost, err := gc.client.Update(context.Background(), &pb.PostObj{
		Id:      val,
		Author:  pp.Author,
		Message: pp.Message,
	})
	if err != nil {
		return model.Post{}, status.Errorf(codes.NotFound, "post not found")
	}

	res, err := uuid.Parse(newPost.Id)
	if err != nil {
		return model.Post{}, status.Error(codes.Internal, "couldn't parse id")
	}

	post := model.Post{
		Id:      res,
		Author:  newPost.Author,
		Message: newPost.Message,
	}
	return post, nil
}

func (gc gRPCClient) Delete(id uuid.UUID) error {

	val := id.String()

	_, err := gc.client.Delete(context.Background(), &pb.PostId{
		Id: val,
	})
	if err != nil {
		return status.Errorf(codes.NotFound, "post not forund")
	}
	return nil
}

func (gc gRPCClient) CreateFromFile(pp model.Post) error {
	val := pp.Id.String()

	_, err := gc.client.CreateFromFile(context.Background(), &pb.PostObj{
		Id:      val,
		Author:  pp.Author,
		Message: pp.Message,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "couldn't create post from file")
	}

	return nil
}
