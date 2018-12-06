package alltogether

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// ErrorArray is simply []error
type ErrorArray []error

// AllNil returns true if all the errors are nil
func (ea ErrorArray) AllNil() bool {
	for _, v := range ea {
		if v != nil {
			return false
		}
	}
	return true
}

func (ea ErrorArray) Error() string {
	accumulator := "errors: "
	nonNilErrorsCounter := 0
	for i := range ea {
		if ea[i] == nil {
			continue
		}
		accumulator = accumulator + ea[i].Error()
		nonNilErrorsCounter++
		if i < len(ea)-1 {
			accumulator = accumulator + ", "
		}
	}
	if nonNilErrorsCounter == 0 {
		accumulator = ""
	}
	return accumulator
}

// errors
var (
	ErrProcessorIsNotFunc = errors.New("processor is not a func")
)

// NewConcurrentProcessor creates a new concurrent processor
func NewConcurrentProcessor(tasks interface{}, processor interface{}) (*Processor, error) {

	// processor must be a func
	if reflect.TypeOf(processor).Kind() != reflect.Func {
		return nil, ErrProcessorIsNotFunc
	}

	ff := reflect.ValueOf(processor)
	// processor must have one parameter
	if ff.Type().NumIn() != 1 {
		return nil, fmt.Errorf("wrong number of parameters for func: got %v, but want 1", ff.Type().NumIn())
	}
	// processor must have one return value
	if ff.Type().NumOut() != 1 {
		return nil, fmt.Errorf("wrong number of return for func: got %v, but want 1", ff.Type().NumOut())
	}
	// the processor return type must be error
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !ff.Type().Out(0).Implements(errorInterface) {
		return nil, errors.New("does not implement error")
	}

	// tasks must be a slice or array
	if reflect.TypeOf(tasks).Kind() != reflect.Slice && reflect.TypeOf(tasks).Kind() != reflect.Array {
		return nil, errors.New("tasks is not an array or slice")
	}

	// the type of task, and the processor param must be of the same type
	if !reflect.DeepEqual(reflect.TypeOf(tasks).Elem(), ff.Type().In(0)) {
		return nil, errors.New("the task type and the param of the processor func do not match")
	}

	pp := &Processor{
		tasks:     tasks,
		processor: processor,
	}

	return pp, nil
}

// Processor is a concurrent processor
type Processor struct {
	tasks     interface{}
	processor interface{}
}

// Do executes a processor function passing each element from tasks, and returning a map of errors;
// each error (or nil) has the same index as the task in the tasks array.
func (p *Processor) Do() ErrorArray {
	s := reflect.ValueOf(p.tasks)
	pr := newBookKeeper(s.Len())

	for i := 0; i < s.Len(); i++ {
		task := s.Index(i).Interface()
		pr.WG.Add(1)
		go pr.singleDo(i, task, p.processor)
	}

	pr.WG.Wait()
	return pr.Errors
}

type bookKeeper struct {
	mu     *sync.RWMutex
	Errors []error
	WG     *sync.WaitGroup
}

func newBookKeeper(size int) *bookKeeper {
	var mapOfErrors = &bookKeeper{
		mu:     &sync.RWMutex{},
		Errors: make([]error, size),
		WG:     &sync.WaitGroup{},
	}
	return mapOfErrors
}

// singleDo calls the user-defined callback and collects the eventual error
func (em bookKeeper) singleDo(index int, task interface{}, processor interface{}) {
	ff := reflect.ValueOf(processor)
	returned := ff.Call([]reflect.Value{reflect.ValueOf(task)})

	em.mu.Lock()
	defer em.mu.Unlock()

	if returned != nil {
		// Get the first returned element (which is an error, or nil)
		e := returned[0]
		if e.Interface() == nil {
			em.Errors[index] = nil
		} else {
			em.Errors[index] = e.Interface().(error)
		}
	}

	em.WG.Done()
}
