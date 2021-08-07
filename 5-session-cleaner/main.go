//////////////////////////////////////////////////////////////////////
//
// Given is a SessionManager that stores session information in
// memory. The SessionManager itself is working, however, since we
// keep on adding new sessions to the manager our program will
// eventually run out of memory.
//
// Your task is to implement a session cleaner routine that runs
// concurrently in the background and cleans every session that
// hasn't been updated for more than 5 seconds (of course usually
// session times are much longer).
//
// Note that we expect the session to be removed anytime between 5 and
// 7 seconds after the last update. Also, note that you have to be
// very careful in order to prevent race conditions.
//

package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

const MAX_TIME = 5

var wg sync.WaitGroup

// SessionManager keeps track of all sessions from creation, updating
// to destroying.
type SessionManager struct {
	sessions map[string]*Session
	mu       sync.Mutex
}

// Session stores the session's data
type Session struct {
	Data                map[string]interface{}
	totalTimePerSession int
}

// NewSessionManager creates a new sessionManager
func NewSessionManager() *SessionManager {
	m := &SessionManager{
		sessions: make(map[string]*Session),
		mu:       sync.Mutex{},
	}

	return m
}

func (m *SessionManager) SessionCleaner(sID string) {
	defer wg.Done()
	tick := time.Tick(time.Second)
	for {
		m.mu.Lock()
		currSession, ok := m.sessions[sID]
		m.mu.Unlock()
		if !ok {
			return
		}
		<-tick
		m.mu.Lock()
		currSession.totalTimePerSession += 1
		m.mu.Unlock()

		if currSession.totalTimePerSession >= MAX_TIME {
			fmt.Println("Deleting the session: ", sID)
			m.mu.Lock()
			delete(m.sessions, sID)
			m.mu.Unlock()
			return
		}
	}
}

// CreateSession creates a new session and returns the sessionID
func (m *SessionManager) CreateSession() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sessionID, err := MakeSessionID()
	if err != nil {
		return "", err
	}

	m.sessions[sessionID] = &Session{
		Data:                make(map[string]interface{}),
		totalTimePerSession: 0,
	}
	wg.Add(1)
	go m.SessionCleaner(sessionID)

	return sessionID, nil
}

// ErrSessionNotFound returned when sessionID not listed in
// SessionManager
var ErrSessionNotFound = errors.New("SessionID does not exists")

// GetSessionData returns data related to session if sessionID is
// found, errors otherwise
func (m *SessionManager) GetSessionData(sessionID string) (map[string]interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	session, ok := m.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session.Data, nil
}

// UpdateSessionData overwrites the old session data with the new one
func (m *SessionManager) UpdateSessionData(sessionID string, data map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	// Hint: you should renew expiry of the session here
	m.sessions[sessionID].Data = data
	m.sessions[sessionID].totalTimePerSession = 0

	return nil
}

func checkIfSessionIsOpen(m *SessionManager, sID string) {
	fmt.Println()
	updatedData, err := m.GetSessionData(sID)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Get session data:", updatedData)
	}
	log.Println("Sleeping for 3 seconds")
	fmt.Println()
	time.Sleep(3 * time.Second)
}

func test(m *SessionManager, sID string) {
	fmt.Println()
	log.Println("@CSB testing")
	fmt.Println()

	checkIfSessionIsOpen(m, sID)

	data := make(map[string]interface{})
	data["website"] = "longhoang.de"
	log.Println("Updating the session")
	err := m.UpdateSessionData(sID, data)
	if err != nil {
		log.Fatal(err)
	}

	checkIfSessionIsOpen(m, sID)

	checkIfSessionIsOpen(m, sID)

	checkIfSessionIsOpen(m, sID)
	wg.Wait()
}

func main() {
	// Create new sessionManager and new session
	m := NewSessionManager()
	sID, err := m.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created new session with ID", sID)

	// Update session data
	data := make(map[string]interface{})
	data["website"] = "longhoang.de"

	err = m.UpdateSessionData(sID, data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Update session data, set website to longhoang.de")

	test(m, sID)
}
