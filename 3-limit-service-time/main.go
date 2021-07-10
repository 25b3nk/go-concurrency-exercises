//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int // in seconds
	mu        sync.Mutex
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	var t1 time.Time
	if !u.IsPremium {
		if u.TimeUsed > 10 {
			// Kill the process if the total usage already exceeded the quota
			return false
		}
		t1 = time.Now()
	}
	process()
	if !u.IsPremium {
		t2 := time.Since(t1)
		u.TimeUsed += int(t2.Seconds())
		if u.TimeUsed > 10 {
			// Kill the process if the current process time or the time accumulated
			// after the current process exceeded the quota
			return false
		}
	}
	return true
}

func main() {
	RunMockServer()
}
