package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
	"sync/atomic"
)

var attempts uint64

func GenerateShortURL(url string, length int) (string, error) {
	attempts := atomic.AddUint64(&attempts, 1)

	saltUrl := fmt.Sprintf("%s%d%d", url, attempts, length)
	hasher := sha1.New()
	if _, err := hasher.Write([]byte(saltUrl)); err != nil {
		return "", fmt.Errorf("failed to hash URL: %w", err)
	}
	hash := hasher.Sum(nil)

	encoded := base64.URLEncoding.EncodeToString(hash)
    encoded = strings.TrimRight(encoded, "=")
    if len(encoded) > length {
        encoded = encoded[:length]
    }

	return encoded, nil
}