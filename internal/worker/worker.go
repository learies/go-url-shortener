package worker

import (
	"github.com/learies/go-url-shortener/internal/models"
)

func GenerateShortURL(deleteUserUrls ...models.UserURL) chan models.UserURL {
	ch := make(chan models.UserURL, len(deleteUserUrls))
	go func() {
		defer close(ch)
		for _, userUrl := range deleteUserUrls {
			ch <- userUrl
		}
	}()
	return ch
}
