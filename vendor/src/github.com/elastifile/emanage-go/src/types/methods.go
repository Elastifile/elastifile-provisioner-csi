package types

import (
	"fmt"
	"math/rand"
	"time"
)

const defaultTimeout = 5 * time.Second

var randContext = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewContext() *Context {
	// TODO: Use a monotonically increasing number? Does this matter?
	id := randContext.Uint32()

	return &Context{
		ID:      id,
		Timeout: defaultTimeout,
	}
}

func (c *Context) String() string {
	return fmt.Sprintf("<Context: %#08x, Timeout: %v>", c.ID, c.Timeout)
}

func NewContextWithTimeout(timeout time.Duration) *Context {
	context := NewContext()
	context.Timeout = timeout
	return context
}

func (f File) String() string {
	const maxLen = 100
	var summary string
	if len(f.Content) < maxLen {
		summary = string(f.Content)
	} else {
		summary = string(f.Content[:maxLen]) + "..."
	}
	return fmt.Sprintf("'%v': '%s'", f.Name, summary)
}

// Support sort.Interface
type ByJobID []Job

func (p ByJobID) Len() int           { return len(p) }
func (p ByJobID) Less(i, j int) bool { return p[i].ID < p[j].ID }
func (p ByJobID) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Support sort.Interface
type ByStartTime []JobInfo

func (p ByStartTime) Len() int           { return len(p) }
func (p ByStartTime) Less(i, j int) bool { return p[i].StartTime.Before(p[j].StartTime) }
func (p ByStartTime) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
