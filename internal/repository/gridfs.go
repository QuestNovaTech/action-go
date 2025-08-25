package repository

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

var (
	gridfsOnce   sync.Once
	gridfsBucket *gridfs.Bucket
	gridfsErr    error
)

// GridFS 返回全局 GridFS 桶实例（懒加载）。
func GridFS() (*gridfs.Bucket, error) {
	gridfsOnce.Do(func() {
		gridfsBucket, gridfsErr = gridfs.NewBucket(DB())
	})
	return gridfsBucket, gridfsErr
} 