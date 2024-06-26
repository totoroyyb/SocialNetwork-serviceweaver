package main

import (
	"context"
	"math/rand"
	"sync"

	"github.com/ServiceWeaver/weaver"
)

type IUrlShortenService interface {
	ComposeUrl(context.Context, []string) ([]Url, error)
	GetExtendedUrls(context.Context, []string) ([]string, error)
	RemoveUrls(context.Context, []string) error
}

type UrlShortenService struct {
	weaver.Implements[IUrlShortenService]

	storage weaver.Ref[IStorage]
}

func (us *UrlShortenService) ComposeUrl(ctx context.Context, urls []string) ([]Url, error) {
	targetUrls := make([]Url, 0)
	for _, url := range urls {
		shortUrl := SHORTEN_URL_HOSTNAME + us.GenRandomStr(10)
		targetUrls = append(targetUrls, Url{
			ShortenedUrl: shortUrl,
			ExpandedUrl:  url,
		})
	}
	var wg sync.WaitGroup
	storage := us.storage.Get()
	for _, url := range targetUrls {
		wg.Add(1)
		go func(url Url) {
			defer wg.Done()
			storage.PutShortenUrl(ctx, url.ShortenedUrl, url.ExpandedUrl)
		}(url)
	}
	wg.Wait()
	return targetUrls, nil
}

func (us *UrlShortenService) GetExtendedUrls(ctx context.Context, shortUrls []string) ([]string, error) {
	var wg sync.WaitGroup
	urlChannel := make(chan string)
	storage := us.storage.Get()
	for _, url := range shortUrls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			extendedUrl, exist, _ := storage.GetShortenUrl(ctx, url)
			if exist {
				urlChannel <- extendedUrl
			}
		}(url)
	}
	wg.Wait()

	var result []string
	for range shortUrls {
		url := <-urlChannel
		result = append(result, url)
	}
	close(urlChannel)
	return result, nil
}

func (us *UrlShortenService) RemoveUrls(ctx context.Context, shortUrls []string) error {
	var wg sync.WaitGroup
	storage := us.storage.Get()
	for _, url := range shortUrls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			storage.RemoveShortenUrl(ctx, url)
		}(url)
	}
	wg.Wait()
	return nil
}

func (us *UrlShortenService) GenRandomStr(length int) string {
	const charMap string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`

	var result string
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charMap))
		result += string(charMap[randomIndex])
	}

	return result
}
