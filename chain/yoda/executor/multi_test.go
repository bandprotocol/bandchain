package executor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockExec struct {
	called int
	result ExecResult
	err    error
}

func newMockExec(output []byte, code uint32, err error) *mockExec {
	return &mockExec{
		called: 0,
		result: ExecResult{Output: output, Code: code},
		err:    err,
	}
}

func (e *mockExec) Exec(code []byte, arg string, env interface{}) (ExecResult, error) {
	e.called++
	return e.result, e.err
}

func TestMultiExecOrderStrategy(t *testing.T) {
	exec1 := newMockExec([]byte("output1"), 1, nil)
	exec2 := newMockExec([]byte("output2"), 2, nil)
	exec, err := NewMultiExec([]Executor{exec1, exec2}, "order")
	require.NoError(t, err)
	// Exec for the first time. Should go to exec1.
	result, err := exec.Exec(nil, "", nil)
	require.NoError(t, err)
	require.Equal(t, result, ExecResult{Output: []byte("output1"), Code: 1})
	require.Equal(t, 1, exec1.called)
	require.Equal(t, 0, exec2.called)
	// Doing it again. Should still go to exec1.
	result, err = exec.Exec(nil, "", nil)
	require.NoError(t, err)
	require.Equal(t, result, ExecResult{Output: []byte("output1"), Code: 1})
	require.Equal(t, 2, exec1.called)
	require.Equal(t, 0, exec2.called)
}

func TestMultiExecRoundRobinStrategy(t *testing.T) {
	exec1 := newMockExec([]byte("output1"), 1, nil)
	exec2 := newMockExec([]byte("output2"), 2, nil)
	exec, err := NewMultiExec([]Executor{exec1, exec2}, "round-robin")
	require.NoError(t, err)
	// Exec for the first time. Should go to exec1.
	result, err := exec.Exec(nil, "", nil)
	require.NoError(t, err)
	require.Equal(t, result, ExecResult{Output: []byte("output1"), Code: 1})
	require.Equal(t, 1, exec1.called)
	require.Equal(t, 0, exec2.called)
	// Doing it again. Should go to exec2.
	result, err = exec.Exec(nil, "", nil)
	require.NoError(t, err)
	require.Equal(t, result, ExecResult{Output: []byte("output2"), Code: 2})
	require.Equal(t, 1, exec1.called)
	require.Equal(t, 1, exec2.called)
	// Doing it again. Should go to exec1.
	result, err = exec.Exec(nil, "", nil)
	require.NoError(t, err)
	require.Equal(t, result, ExecResult{Output: []byte("output1"), Code: 1})
	require.Equal(t, 2, exec1.called)
	require.Equal(t, 1, exec2.called)
}

func TestMultiExecBadStrategy(t *testing.T) {
	_, err := NewMultiExec([]Executor{}, "bad")
	require.EqualError(t, err, "unknown MultiExec strategy: bad")
}

func TestMultiExecOneWorking(t *testing.T) {
	exec1 := newMockExec(nil, 0, errors.New("error1"))
	exec2 := newMockExec([]byte("output"), 0, nil)
	exec3 := newMockExec(nil, 0, errors.New("error3"))
	exec, err := NewMultiExec([]Executor{exec1, exec2, exec3}, "round-robin")
	result, err := exec.Exec(nil, "", nil)
	require.NoError(t, err)
	require.Equal(t, ExecResult{Output: []byte("output"), Code: 0}, result)
	result, err = exec.Exec(nil, "", nil)
	require.NoError(t, err)
	require.Equal(t, ExecResult{Output: []byte("output"), Code: 0}, result)
	result, err = exec.Exec(nil, "", nil)
	require.NoError(t, err)
	require.Equal(t, ExecResult{Output: []byte("output"), Code: 0}, result)
}

func TestMultiExecAllErrors(t *testing.T) {
	exec1 := newMockExec(nil, 0, errors.New("error1"))
	exec2 := newMockExec(nil, 0, errors.New("error2"))
	exec3 := newMockExec(nil, 0, errors.New("error3"))
	exec, err := NewMultiExec([]Executor{exec1, exec2, exec3}, "round-robin")
	_, err = exec.Exec(nil, "", nil)
	require.EqualError(t, err, "MultiError: error1, error2, error3")
	_, err = exec.Exec(nil, "", nil)
	require.EqualError(t, err, "MultiError: error2, error3, error1")
	_, err = exec.Exec(nil, "", nil)
	require.EqualError(t, err, "MultiError: error3, error1, error2")
}
