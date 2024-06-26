package main

import (
	"context"
	"sync"

	"github.com/ServiceWeaver/weaver"
)

type IHomeTimelineService interface {
	ReadHomeTimeline(context.Context, int64, int, int) ([]Post, error)
	WriteHomeTimeline(context.Context, int64, int64, int64, []int64) error
	RemovePost(context.Context, int64, int64, int64) error
}

type HomeTimelineService struct {
	weaver.Implements[IHomeTimelineService]

	postStorageService weaver.Ref[PostStorageServicer]
	socialGraphService weaver.Ref[ISocialGraphService]
	storage            weaver.Ref[IStorage]
}

func (hts *HomeTimelineService) ReadHomeTimeline(ctx context.Context, userId int64, start int, stop int) ([]Post, error) {
	if stop <= start || start < 0 {
		return make([]Post, 0), nil
	}
	storage := hts.storage.Get()
	postStorageService := hts.postStorageService.Get()

	postIds, _ := storage.GetPostTimeline(ctx, userId, start, stop)
	return postStorageService.ReadPosts(ctx, postIds)
}

func (hts *HomeTimelineService) WriteHomeTimeline(ctx context.Context, postId int64, userId int64, timestamp int64, userMentionIds []int64) error {
	storage := hts.storage.Get()
	socialGraphService := hts.socialGraphService.Get()
	ids, _ := socialGraphService.GetFollowers(ctx, userId)
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(ctx context.Context, id, postId, timestamp int64) {
			defer wg.Done()
			storage.PutPostTimeline(ctx, id, postId, timestamp)
		}(ctx, id, postId, timestamp)
	}
	wg.Wait()
	return nil
}

func (hts *HomeTimelineService) RemovePost(ctx context.Context, userId int64, postId int64, timestamp int64) error {
	storage := hts.storage.Get()
	storage.RemovePostTimeline(ctx, userId, postId, timestamp)
	return nil
}
