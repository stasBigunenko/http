package gRPC

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "src/http/api/proto"
	"src/http/pkg/model"
	"src/http/storage"
)

type StorageServer struct {
	pb.UnimplementedPostServiceServer

	Store storage.Storage
}

func NewGRPCStore(store storage.Storage) *StorageServer {
	return &StorageServer{
		Store: store,
	}
}

func (s *StorageServer) Get(_ context.Context, in *pb.PostId) (*pb.PostObj, error) {
	val, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't parse id")
	}

	post, err := s.Store.Get(val)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to get post")
	}

	res := post.Id.String()

	return &pb.PostObj{
		Id:      res,
		Author:  post.Author,
		Message: post.Message,
	}, nil
}

func (s *StorageServer) GetAll(_ context.Context, in *emptypb.Empty) (*pb.AllPosts, error) {
	posts := s.Store.GetAll()

	pbPosts := []*pb.PostObj{}

	for _, val := range posts {
		res := val.Id.String()
		pbPosts = append(pbPosts, &pb.PostObj{
			Id:      res,
			Author:  val.Author,
			Message: val.Message,
		})
	}

	return &pb.AllPosts{
		Allposts: pbPosts,
	}, nil
}

func (s *StorageServer) Create(_ context.Context, in *pb.NewPost) (*pb.PostObj, error) {
	p := model.Post{
		Author:  in.Author,
		Message: in.Message,
	}

	post, err := s.Store.Create(p)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal storage problem")
	}

	res := post.Id.String()

	return &pb.PostObj{
		Id:      res,
		Author:  post.Author,
		Message: post.Message,
	}, nil
}

func (s *StorageServer) Update(_ context.Context, in *pb.PostObj) (*pb.PostObj, error) {
	val, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create post")
	}

	post := model.Post{
		Id:      val,
		Author:  in.Author,
		Message: in.Message,
	}

	res, err := s.Store.Update(post)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to update post")
	}

	resId := res.Id.String()

	return &pb.PostObj{
		Id:      resId,
		Author:  res.Author,
		Message: res.Message,
	}, nil
}

func (s *StorageServer) Delete(_ context.Context, in *pb.PostId) (*emptypb.Empty, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create post")
	}

	err = s.Store.Delete(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to delete post")
	}
	return &emptypb.Empty{}, nil
}

func (s *StorageServer) CreateFromFile(_ context.Context, in *pb.PostObj) (*emptypb.Empty, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create post")
	}

	post := model.Post{
		Id:      id,
		Author:  in.Author,
		Message: in.Message,
	}

	err = s.Store.CreateFromFile(post)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed due to internal problem")
	}
	return &emptypb.Empty{}, nil
}
