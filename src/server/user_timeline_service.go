package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

type IUserTimelineService interface {
	WriteUserTimeline(context.Context, int64, int64, int64) error
	ReadUserTimeline(context.Context, int64, int, int) ([]Post, error)
	RemovePost(context.Context, int64, int64, int64) error
}

type UserTimelineService struct {
	weaver.Implements[IUserTimelineService]
	storage            weaver.Ref[IStorage]
	postStorageService weaver.Ref[PostStorageServicer]
}

func (uts *UserTimelineService) WriteUserTimeline(ctx context.Context, postId, userId, timestamp int64) error {
	storage := uts.storage.Get()
	storage.PutPostTimeline(ctx, userId, postId, timestamp)
	return nil
}

func (uts *UserTimelineService) ReadUserTimeline(ctx context.Context, userId int64, start int, stop int) ([]Post, error) {
	storage := uts.storage.Get()
	postStorageService := uts.postStorageService.Get()
	postIds, err := storage.GetPostTimeline(ctx, userId, start, stop)
	if err != nil {
		return make([]Post, 0), err
	}
	return postStorageService.ReadPosts(ctx, postIds)
}

func (uts *UserTimelineService) RemovePost(ctx context.Context, userId int64, postId int64, timestamp int64) error {
	storage := uts.storage.Get()
	storage.RemovePostTimeline(ctx, userId, postId, timestamp)
	return nil
}
