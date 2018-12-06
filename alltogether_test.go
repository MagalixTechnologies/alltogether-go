package alltogether

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrongParams(t *testing.T) {
	arr := []string{"uno", "due", "tre", "quattro", "cinque", "sei", "sette", "otto", "nove", "dieci"}

	// wrong tasks
	_, err := NewConcurrentProcessor(1, func(task string) error {
		fmt.Println(task)
		if task == "due" {
			return errors.New("some error for due")
		}
		return nil
	})
	assert.NotNil(t, err, "error should be not nil")

	// wrong processor func
	_, err = NewConcurrentProcessor(arr, "this is a string")
	assert.NotNil(t, err, "error should be not nil")

	// no return for processor func
	_, err = NewConcurrentProcessor(arr, func(task int) {
		fmt.Println(task)
	})
	assert.NotNil(t, err, "error should be not nil")

	// wrong processor func return type
	_, err = NewConcurrentProcessor(arr, func(task int) int {
		fmt.Println(task)
		return 0
	})
	assert.NotNil(t, err, "error should be not nil")

	// mixed types (task and func param)
	_, err = NewConcurrentProcessor(arr, func(task int) error {
		fmt.Println(task)
		return nil
	})
	assert.NotNil(t, err, "error should be not nil")
}
