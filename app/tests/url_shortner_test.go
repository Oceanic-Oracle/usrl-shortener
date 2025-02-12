package tests

import (
	"testing"
	"url-shortener/internal/api/gen"
	"url-shortener/tests/mocks"
	"url-shortener/tests/suite"
)

func TestRegUrl(t *testing.T) {
	ctx, suite := suite.New(t)

	testURL := [...]string{"https://example.com", "https://youtube.com"}

	for _, val := range testURL {
		resp, err := suite.UrlClient.RegUrl(ctx,
			&gen.RegUrlReq{
				Url: val,
			})
		if err != nil {
			t.Fatalf("RegUrl failed: %v", err)
		}

		if resp.ShortUrl == "" {
			t.Error("Expected non-empty short URL")
		}

		t.Logf("Registered URL: %s -> %s", val, resp.ShortUrl)
	}
}

func TestGetUrl(t *testing.T) {
	ctx, suite := suite.New(t)

	tests := mocks.GetUrlMocks(ctx, suite.Storage)

	for key, val := range *tests {
		resp, err := suite.UrlClient.GetUrl(ctx,
			&gen.GetUrlReq{
				ShortUrl: val,
			})
		if err != nil {
			t.Fatalf("RegUrl failed: %v", err)
		}

		if resp.Url != key {
			t.Error("Expected non-empty short URL")
		}

		t.Logf("Registered URL: %s -> %s", val, resp.Url)
	}
}
