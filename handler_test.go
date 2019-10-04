package bucket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrepareBucket(t *testing.T) {
	const durationInSec = 1
	h, err := New(&Configuration{
		Username:       "Administrator",
		Password:       "password",
		BucketName:     bucketName,
		BucketPassword: "",
		Separator:      "::",
		Opts: Opts{
			OperationTimeout:      NullTimeoutSec(durationInSec),
			BulkOperationTimeout:  NullTimeoutSec(durationInSec),
			DurabilityTimeout:     NullTimeoutSec(durationInSec),
			DurabilityPollTimeout: NullTimeoutSec(durationInSec),
			ViewTimeout:           NullTimeoutSec(durationInSec),
			N1qlTimeout:           NullTimeoutSec(durationInSec),
			AnalyticsTimeout:      NullTimeoutSec(durationInSec),
		},
	})

	assert.Nil(t, err)
	assert.NotNil(t, h)

	assert.Equal(t, durationInSec*time.Second, h.state.bucket.OperationTimeout())
	assert.Equal(t, durationInSec*time.Second, h.state.bucket.BulkOperationTimeout())
	assert.Equal(t, durationInSec*time.Second, h.state.bucket.DurabilityTimeout())
	assert.Equal(t, durationInSec*time.Second, h.state.bucket.DurabilityPollTimeout())
	assert.Equal(t, durationInSec*time.Second, h.state.bucket.ViewTimeout())
	assert.Equal(t, durationInSec*time.Second, h.state.bucket.N1qlTimeout())
	assert.Equal(t, durationInSec*time.Second, h.state.bucket.AnalyticsTimeout())
}

func TestNewStateError(t *testing.T) {
	_, err := New(&Configuration{ConnectionString: ""})

	assert.NotNil(t, err)
	assert.Equal(t, "no access", err.Error())
}
