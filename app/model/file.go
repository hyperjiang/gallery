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
	CreatedAt time.Time `db:"created_at"`
}

// Files - file list
type Files []*File

// Create - create a file record
func (f *File) Create() error {

	f.CreatedAt = time.Now()

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

// CountFilesByType - count files by type
func CountFilesByType(t string) (uint, error) {
	var count uint
	err := provider.DI().DB().Get(&count, `
		SELECT COUNT(0) FROM file
		WHERE type LIKE ?
	`, t+"%")
	return count, err
}

// GetByType - get file list by type
func (fs *Files) GetByType(t string, limit, offset uint) error {
	return provider.DI().DB().Select(fs, `
		SELECT * FROM file
		WHERE type LIKE ?
		ORDER BY ID DESC
		LIMIT ? OFFSET ?
    `, t+"%", limit, offset)
}
