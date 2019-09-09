package hexcore

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicClean(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	password := "my-super-secret-password"
	router.Use(RecoveryWithWriter(buffer))
	router.Handle("/recovery", func(c *Context) {
		//c.AbortWithStatus(http.StatusBadRequest)
		panic("Oupps, Houston, we have a problem")
	})
	// RUN
	performRequest(router, "/recovery")
	// TEST

	// Check the buffer does not have the secret key
	assert.NotContains(t, buffer.String(), password)
}

// TestPanicInHandler assert that panic has been recovered.
func TestPanicInHandler(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(RecoveryWithWriter(buffer))
	router.Handle("/recovery", func(_ *Context) {
		panic("Oupps, Houston, we have a problem")
	})
	// RUN
	performRequest(router, "/recovery")
	// TEST
	assert.Contains(t, buffer.String(), "panic recovered")
	assert.Contains(t, buffer.String(), "Oupps, Houston, we have a problem")
	assert.Contains(t, buffer.String(), "TestPanicInHandler")
}

// TestPanicWithAbort assert that panic has been recovered even if context.Abort was used.
func TestPanicWithAbort(t *testing.T) {
	router := New()
	router.Use(RecoveryWithWriter(nil))
	router.Handle("/recovery", func(c *Context) {
		c.Abort()
		panic("Oupps, Houston, we have a problem")
	})
	// RUN
	performRequest(router, "/recovery")
	// TEST
	// assert.Contains(t, buffer.String(), "Oupps, Houston, we have a problem")
}

func TestSource(t *testing.T) {
	bs := source(nil, 0)
	assert.Equal(t, []byte("???"), bs)

	in := [][]byte{
		[]byte("Hello world."),
		[]byte("Hi, gin.."),
	}
	bs = source(in, 10)
	assert.Equal(t, []byte("???"), bs)

	bs = source(in, 1)
	assert.Equal(t, []byte("Hello world."), bs)
}

func TestFunction(t *testing.T) {
	bs := function(1)
	assert.Equal(t, []byte("???"), bs)
}
