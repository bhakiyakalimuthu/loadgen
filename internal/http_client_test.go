package internal

import (
	"go.uber.org/zap"
	"testing"
)

func BenchmarkHttpClient_Notify(b *testing.B) {
	b.Run("benchHttpClient", func(b *testing.B) {
		n := NewHttpClient(zap.NewNop(), "https://httpbin.org/")
		benchHttpClient(n)
	})
}

func benchHttpClient(client HttpClient) {
	s := []string{"hello", "hi", "how are you", "beautiful"}
	for _, i := range s {
		client.GenerateLoad(i)
	}

}
