package utils

import (
	"sync"

	"github.com/silenceper/wechat/v2/cache"
)

var memCache *cache.Memory
var once sync.Once

func GetCache() *cache.Memory {
	once.Do(func() {
		memCache = cache.NewMemory()
	})
	return memCache
}
