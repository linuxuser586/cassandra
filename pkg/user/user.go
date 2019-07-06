package user

import (
	"os"
	"strconv"
)

const defaultID = 9042

// UID is the user ID
var UID int

// GID is the group ID
var GID int

func init() {
	initUID()
	initGID()
}

func initUID() {
	id := os.Getenv("PUID")
	if id == "" {
		UID = defaultID
	}
	i, err := strconv.Atoi(id)
	if err != nil {
		UID = defaultID
	} else {
		UID = i
	}
}

func initGID() {
	id := os.Getenv("PGID")
	if id == "" {
		GID = defaultID
	}
	i, err := strconv.Atoi(id)
	if err != nil {
		GID = defaultID
	} else {
		GID = i
	}
}
