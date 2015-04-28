package session

import (
	"net/http"
)

type Store interface {
	//initiate a new session (normally saved it into the cookie of users)
	New(*http.Request, http.ResponseWriter) (*Session, error)
	//save the session into backend (redis)
	Save(http.ResponseWriter, *Cache) error
	//get the session out of the backend (redis)
	Get(*http.Request) (*Cache, error)
	//purge the storage of session
	Delete(http.ResponseWriter, *Cache) error
}
