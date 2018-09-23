package model

import (
	"regexp"
	"time"

	"github.com/hyperjiang/gallery-service/app/provider"
)

// File - file table
type File struct {
	ID        uint32    `db:"id"`
	Name      string    `db:"name"`
	Path      string    `db:"path"`
	Type      string    `db:"type"`
	Size      uint32    `db:"size"`
	Checksum  string    `db:"checksum"`
	UserID    uint32    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Files - file list
type Files []*File

// Create - create a file record
func (f *File) Create() error {

	now := time.Now()
	f.CreatedAt = now
	f.UpdatedAt = now

	res, err := provider.DI().DBInsert("file", f)
	if err != nil {
		return err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	f.ID = uint32(lastInsertID)

	return err
}

// GetByChecksum - get file by checksum
func (f *File) GetByChecksum(checksum string) error {
	return provider.DI().DB().Get(f, `
        SELECT * FROM file WHERE checksum = ? LIMIT 1
    `, checksum)
}

// IsImage - return true when the file is an image
func (f *File) IsImage() bool {
	reg := regexp.MustCompile("^(.+)\\/.+$")
	matches := reg.FindStringSubmatch(f.Type)
	return len(matches) == 2 && matches[1] == "image"
}

// IsVideo - return true when the file is a video
func (f *File) IsVideo() bool {
	reg := regexp.MustCompile("^(.+)\\/.+$")
	matches := reg.FindStringSubmatch(f.Type)
	return len(matches) == 2 && matches[1] == "video"
}

// Get - get file list
func (fs *Files) Get() error {
	return provider.DI().DB().Select(fs, `
        SELECT * FROM file ORDER BY ID DESC;
    `)
}
