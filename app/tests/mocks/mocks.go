package mocks

import (
	"context"
	"url-shortener/internal/storage"
	"url-shortener/internal/utils"
)

func GetUrlMocks(ctx context.Context, storage *storage.Repositories) *map[string]string {
	urls := [...]string{"https://youtube.com", "https://google.com"}
	shortUrls := make(map[string]string, len(urls))

	for ind := range urls {
		var attempts int
		var err error
		for attempts < 5 {
			shortUrls[urls[ind]], err = utils.GenerateShortURL(urls[ind], 10)
			if err == nil {
				break
			}
			attempts++
		}
		storage.Url.SaveUrl(ctx, urls[ind], shortUrls[urls[ind]])
	}
	return &shortUrls
}