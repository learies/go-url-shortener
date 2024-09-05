package shortener

import (
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"sync"
)

type URLShortener struct {
	hasher   hash.Hash
	hasherMu sync.Mutex
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		hasher: sha256.New(),
	}
}

func (us *URLShortener) Get() hash.Hash {
	us.hasherMu.Lock()
	defer us.hasherMu.Unlock()
	return us.hasher
}

func (us *URLShortener) Set(newHasher hash.Hash) {
	us.hasherMu.Lock()
	defer us.hasherMu.Unlock()
	us.hasher = newHasher
}

func (us *URLShortener) GenerateShortURL(url string) string {
	hasher := us.Get()

	hasher.Reset()
	hasher.Write([]byte(url))
	hash := hasher.Sum(nil)
	shortURL := base64.URLEncoding.EncodeToString(hash)[:8]

	us.Set(hasher)

	return shortURL
}
