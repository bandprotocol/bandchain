package executor

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

// MutliExec is a higher-order executor that utlizes the underlying executors to perform Exec.
type MultiExec struct {
	execs    []Executor // The list of underlying executors.
	strategy string     // Execution strategy. Can be "order" order "round-robin".
	rIndex   int64      // Round-robin index, only applicable when strategy is "round-robin"
}

// MultiError encapsulates error messages from the underlying executors into one error.
type MultiError struct {
	errs []error
}

// NewMultiExec creates a new MultiExec instance.
func NewMultiExec(execs []Executor, strategy string) (*MultiExec, error) {
	if strategy != "order" && strategy != "round-robin" {
		return &MultiExec{}, fmt.Errorf("unknown MultiExec strategy: %s", strategy)
	}
	return &MultiExec{execs: execs, strategy: strategy, rIndex: -1}, nil
}

// Error implements error interface for MultiError by returning all error messages concatenated.
func (e *MultiError) Error() string {
	var s strings.Builder
	s.WriteString("MultiError: ")
	for idx, each := range e.errs {
		if idx != 0 {
			s.WriteString(", ")
		}
		s.WriteString(each.Error())
	}
	return s.String()
}

// nextExecOrder returns the next order of executors to be used by MultiExec.
func (e *MultiExec) nextExecOrder() []Executor {
	switch e.strategy {
	case "order":
		return e.execs
	case "round-robin":
		rIndex := atomic.AddInt64(&e.rIndex, 1) % int64(len(e.execs))
		return append(append([]Executor{}, e.execs[rIndex:]...), e.execs[:rIndex]...)
	default:
		panic("unknown MultiExec strategy") // We should never reach here.
	}
}

// Exec implements Executor interface for MultiExec.
func (e *MultiExec) Exec(timeout time.Duration, code []byte, arg string) (ExecResult, error) {
	errs := []error{}
	for _, each := range e.nextExecOrder() {
		res, err := each.Exec(timeout, code, arg)
		if err == nil || err == ErrExecutionimeout {
			return res, err
		} else {
			errs = append(errs, err)
		}
	}
	return ExecResult{}, &MultiError{errs: errs}
}
