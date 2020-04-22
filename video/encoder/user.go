package encoder

import (
	"math/rand"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/video"
)

// UserCallback is a callback function for when a video gets updated
type UserCallback func(*video.Status)

// UserRequests keeps track of all the user's requests and their callbacks
type UserRequests struct {
	user      *sbvision.User
	m         *VideoRequestManager
	requests  map[string]*videoRequest
	callbacks map[int64]UserCallback
}

// GetUserRequests gets the object tracking this user's video requests
func (m *VideoRequestManager) GetUserRequests(user *sbvision.User) *UserRequests {
	if _, exists := m.userRequests[user.Email]; !exists {
		m.userRequests[user.Email] = &UserRequests{
			m:         m,
			user:      user,
			requests:  make(map[string]*videoRequest),
			callbacks: make(map[int64]UserCallback),
		}
	}
	return m.userRequests[user.Email]
}

// AddListener adds the callback function to the list of user callbacks
func (u *UserRequests) AddListener(callback UserCallback) int64 {
	id := rand.Int63()
	u.callbacks[id] = callback
	for _, request := range u.requests {
		callback(&request.Status)
	}
	return id
}

// RemoveListener deletes the given id from the list of callbacks
func (u *UserRequests) RemoveListener(id int64) {
	delete(u.callbacks, id)
}
