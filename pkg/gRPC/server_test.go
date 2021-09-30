package gRPC

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "src/http/api/proto"
	"src/http/pkg/model"
	"src/http/storage/mocks"
	"testing"
)

func TestStorageServer_Get(t *testing.T) {
	s := new(mocks.Storage)
	id := uuid.New()
	idStr := id.String()
	p := model.Post{id, "stas", "the first"}
	s.On("Get", id).Return(p, nil)

	s2 := new(mocks.Storage)
	s2.On("Get", mock.Anything).Return(model.Post{}, errors.New("post with Id %d not found"))

	tests := []struct {
		name    string
		stor    *mocks.Storage
		param   *pb.PostId
		want    *pb.PostObj
		wantErr codes.Code
	}{
		{
			name:  "Get everything good",
			stor:  s,
			param: &pb.PostId{Id: idStr},
			want:  &pb.PostObj{Id: idStr, Author: "stas", Message: "the first"},
		},
		{
			name:    "Get everything err parse",
			stor:    s2,
			param:   &pb.PostId{Id: "123"},
			want:    nil,
			wantErr: codes.Internal,
		},
		{
			name:    "Get everything err not found",
			stor:    s2,
			param:   &pb.PostId{Id: idStr},
			want:    nil,
			wantErr: codes.NotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewGRPCStore(tc.stor)
			got, err := u.Get(context.Background(), tc.param)
			if err != nil && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageServer_GetAll(t *testing.T) {
	s := new(mocks.Storage)
	id1 := uuid.New()
	idStr1 := id1.String()
	id2 := uuid.New()
	idStr2 := id2.String()
	x := []model.Post{
		{id1, "Stas", "first"},
		{id2, "Stas", "second"},
	}

	p1 := &pb.PostObj{Id: idStr1, Author: "Stas", Message: "first"}
	p2 := &pb.PostObj{Id: idStr2, Author: "Stas", Message: "second"}
	all := []*pb.PostObj{p1, p2}
	aa := pb.AllPosts{
		Allposts: all,
	}
	s.On("GetAll").Return(x, nil)

	s2 := new(mocks.Storage)
	s2.On("Get", mock.Anything).Return(model.Post{}, errors.New("post with Id %d not found"))

	tests := []struct {
		name    string
		stor    *mocks.Storage
		want    *pb.AllPosts
		wantErr codes.Code
	}{
		{
			name: "GetAll everything good",
			stor: s,
			want: &aa,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewGRPCStore(tc.stor)
			got, err := u.GetAll(context.Background(), &emptypb.Empty{})
			if err != nil && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageServer_Create(t *testing.T) {
	s := new(mocks.Storage)
	id := uuid.New()
	idStr := id.String()
	p := model.Post{id, "stas", "the first"}
	s.On("Create", model.Post{Author: "stas", Message: "the first"}).Return(p, nil)

	s2 := new(mocks.Storage)
	s2.On("Create", mock.Anything).Return(model.Post{}, errors.New("post with Id %d not found"))

	tests := []struct {
		name    string
		stor    *mocks.Storage
		param   *pb.NewPost
		want    *pb.PostObj
		wantErr codes.Code
	}{
		{
			name:  "Create everything good",
			stor:  s,
			param: &pb.NewPost{Author: "stas", Message: "the first"},
			want:  &pb.PostObj{Id: idStr, Author: "stas", Message: "the first"},
		},
		{
			name:    "Create err",
			stor:    s2,
			param:   &pb.NewPost{Author: "stas", Message: "the first"},
			want:    nil,
			wantErr: codes.Internal,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewGRPCStore(tc.stor)
			got, err := u.Create(context.Background(), tc.param)
			if err != nil && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageServer_Update(t *testing.T) {
	s := new(mocks.Storage)
	id := uuid.New()
	idStr := id.String()
	p := model.Post{id, "stas", "the first"}
	s.On("Update", model.Post{Id: id, Author: "stas", Message: "the first"}).Return(p, nil)

	s2 := new(mocks.Storage)
	s2.On("Update", mock.Anything).Return(model.Post{}, errors.New("post with Id %d not found"))

	tests := []struct {
		name    string
		stor    *mocks.Storage
		param   *pb.PostObj
		want    *pb.PostObj
		wantErr codes.Code
	}{
		{
			name:  "Update everything good",
			stor:  s,
			param: &pb.PostObj{Id: idStr, Author: "stas", Message: "the first"},
			want:  &pb.PostObj{Id: idStr, Author: "stas", Message: "the first"},
		},
		{
			name:    "Update err",
			stor:    s2,
			param:   &pb.PostObj{Id: idStr, Author: "stas", Message: "the first"},
			want:    nil,
			wantErr: codes.InvalidArgument,
		},
		{
			name:    "Update err",
			stor:    s2,
			param:   &pb.PostObj{Id: "123", Author: "stas", Message: "the first"},
			want:    nil,
			wantErr: codes.Internal,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewGRPCStore(tc.stor)
			got, err := u.Update(context.Background(), tc.param)
			if err != nil && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageServer_Delete(t *testing.T) {
	s := new(mocks.Storage)
	id := uuid.New()
	idStr := id.String()
	s.On("Delete", id).Return(nil)

	s2 := new(mocks.Storage)
	s2.On("Delete", mock.Anything).Return(errors.New("post with Id %d not found"))

	tests := []struct {
		name    string
		stor    *mocks.Storage
		param   *pb.PostId
		want    *emptypb.Empty
		wantErr codes.Code
	}{
		{
			name:  "Create everything good",
			stor:  s,
			param: &pb.PostId{Id: idStr},
			want:  &emptypb.Empty{},
		},
		{
			name:    "Delete err",
			stor:    s2,
			param:   &pb.PostId{Id: "123"},
			want:    nil,
			wantErr: codes.Internal,
		},
		{
			name:    "Delete err",
			stor:    s2,
			param:   &pb.PostId{Id: idStr},
			want:    nil,
			wantErr: codes.NotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewGRPCStore(tc.stor)
			got, err := u.Delete(context.Background(), tc.param)
			if err != nil && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageServer_CreateFromFile(t *testing.T) {
	s := new(mocks.Storage)
	id := uuid.New()
	idStr := id.String()
	s.On("CreateFromFile", model.Post{Id: id, Author: "stas", Message: "the first"}).Return(nil)

	s2 := new(mocks.Storage)
	s2.On("CreateFromFile", mock.Anything).Return(errors.New("post with Id %d not found"))

	tests := []struct {
		name    string
		stor    *mocks.Storage
		param   *pb.PostObj
		want    *emptypb.Empty
		wantErr codes.Code
	}{
		{
			name:  "Create everything good",
			stor:  s,
			param: &pb.PostObj{Id: idStr, Author: "stas", Message: "the first"},
			want:  &emptypb.Empty{},
		},
		{
			name:    "Create err",
			stor:    s2,
			param:   &pb.PostObj{Id: idStr, Author: "stas", Message: "the first"},
			want:    nil,
			wantErr: codes.Internal,
		},
		{
			name:    "Create err",
			stor:    s2,
			param:   &pb.PostObj{Id: "123", Author: "stas", Message: "the first"},
			want:    nil,
			wantErr: codes.Internal,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewGRPCStore(tc.stor)
			got, err := u.CreateFromFile(context.Background(), tc.param)
			if err != nil && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
