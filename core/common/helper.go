package common

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/gabriel-vasile/mimetype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DateTimeNullableToString(dateTimeNullable *time.Time) string {
	if dateTimeNullable == nil {
		return ""
	}

	if dateTimeNullable.IsZero() {
		return ""
	}
	return dateTimeNullable.Format(time.RFC3339)
}

func DToMap(d primitive.D) map[string]interface{} {
	result := make(map[string]interface{})
	for _, elem := range d {
		result[elem.Key] = elem.Value
	}
	return result
}

func IsSlugValid(slug string) bool {
	// Check if slug is empty
	if slug == "" {
		return false
	}

	// Check if slug contains only lowercase letters, numbers, or hyphens
	for _, char := range slug {
		if !unicode.IsLower(char) && !unicode.IsDigit(char) && char != '-' {
			return false
		}
	}

	// Check if slug starts or ends with a hyphen
	if slug[0] == '-' || slug[len(slug)-1] == '-' {
		return false
	}

	// Check if slug has consecutive hyphens
	for i := 0; i < len(slug)-1; i++ {
		if slug[i] == '-' && slug[i+1] == '-' {
			return false
		}
	}

	return true
}

func GetFileMetaData(fileName string) (*GotaroFileMetaData, error) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the file size
	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Reset the read pointer to the start of the file
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// Get the MIME type
	mimeData, err := mimetype.DetectReader(file)
	if err != nil {
		return nil, err
	}

	return &GotaroFileMetaData{
		FileExt:  mimeData.Extension(),
		FileMime: mimeData.String(),
		FileSize: uint64(fileStat.Size()),
	}, nil
}

func GetFileNameUnique(filename string) string {
	//split filename and extension
	extension := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, extension)
	// place all special character with _
	filename = strings.Map(func(r rune) rune {
		switch r {
		case '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '+', '=', '{', '}', '[', ']', '|', '\\', '<', '>', '~', '`':
			return '_'
		}
		return r
	}, filename)
	// add random string to make it unique
	filename = filename + "_" + RandomString(10)

	return filename + extension
}

func IsMimeValid(allowedMimes []string, fileMime string) bool {
	for _, allowedMime := range allowedMimes {
		if fileMime == allowedMime {
			return true
		}
	}
	return false
}

func CreateFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
