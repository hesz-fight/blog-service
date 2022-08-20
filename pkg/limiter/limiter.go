package limiter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

type LimiterIface interface {
	// get key of buckets
	Key(c *gin.Context) string
	// get bucket named key
	GetBucket(key string) (*ratelimit.Bucket, bool)
	// add buckets
	AddBuckets(rules ...LimiterBucketRule) LimiterIface
}

// Limiter limiter instance
type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}

type LimiterBucketRule struct {
	// key value name
	Key string
	// interval of put token bucket
	FillInterval time.Duration
	// capacity of token bucket
	Capacity int64
	// number of buckets putting when time is up
	Quantum int64
}
