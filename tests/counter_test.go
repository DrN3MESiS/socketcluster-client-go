package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/daominah/socketcluster-client-go/scclient/utils"
)

func TestShouldReturnIncrementedValue(t *testing.T) {
	counter := utils.AtomicCounter{
		Counter: 0,
	}

	actualValue := counter.IncrementAndGet()

	assert.Equal(t, uint64(1), actualValue)
	assert.Equal(t, uint64(1), counter.Value())
}

func TestShouldGetAndIncrementValue(t *testing.T) {
	counter := utils.AtomicCounter{
		Counter: 0,
	}

	actualValue := counter.GetAndIncrement()

	assert.Equal(t, uint64(0), actualValue)
	assert.Equal(t, uint64(1), counter.Value())
}
