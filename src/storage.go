package main

import (
	"fmt"
	"github.com/google/uuid"
)

// Storing element IDs in the hash map to make sure they are unique
type IdStorage struct {
	Storage map[string]struct{}
}

func newIdStorage() *IdStorage {
	return &IdStorage{
		Storage: make(map[string]struct{}, 0),
	}
}

func (storage *IdStorage) push(new_id string) ([]byte, error) {

	if len(new_id) == 0 {
		new_id = uuid.New().String()
	}

	_, already_exists := storage.Storage[new_id]

	if already_exists {
		return nil, fmt.Errorf("ID '%s' declared twice!", new_id)
	}

	storage.Storage[new_id] = struct{}{}

	return []byte(new_id), nil
}

// Storing local references to make sure we don't reference the IDs we don't have
type LinkStorage struct {
	Storage map[string]struct{}
}

func newLinkStorage() *LinkStorage {
	return &LinkStorage{
		Storage: make(map[string]struct{}, 0),
	}
}

func (storage *LinkStorage) push(new_link string) {
	storage.Storage[new_link] = struct{}{}
}

// Storing image paths to copy them to destination
type ImageStorage struct {
	Storage map[string]struct{}
}

func newImageStorage() *ImageStorage {
	return &ImageStorage{
		Storage: make(map[string]struct{}, 0),
	}
}

func (images *ImageStorage) dump() []string {

	result := make([]string, 0, len(images.Storage))

	for k, _ := range images.Storage {
		result = append(result, k)
	}
	return result
}

func (storage *ImageStorage) push(new_link string) {
	storage.Storage[new_link] = struct{}{}
}

type HeaderStorage struct {
	Headers          []*HeaderInfo
	LastHeaderNumber []int
}

type HeaderInfo struct {
	Text   []byte
	Level  int
	Number []int
	ID     string
}

func (storage *HeaderStorage) push(level int, header []byte, id string) ([]int, error) {

	if level < 1 {
		return nil, fmt.Errorf("Invalid level '%d', levels <1 not allowed", level)
	}

	if level <= len(storage.LastHeaderNumber) {
		storage.LastHeaderNumber[level-1] += 1
	} else {
		for i := len(storage.LastHeaderNumber); i <= level; i++ {
			storage.LastHeaderNumber = append(storage.LastHeaderNumber, 1)
		}
	}
	storage.LastHeaderNumber = storage.LastHeaderNumber[:level]

	new_number := append([]int(nil), storage.LastHeaderNumber...)

	header_info := &HeaderInfo{
		Text:   header,
		Level:  level,
		Number: new_number,
		ID:     id,
	}

	storage.Headers = append(storage.Headers, header_info)

	return new_number, nil
}

func checkAllLinksValid(links *LinkStorage, ids *IdStorage) error {
	return nil
}

func (storage *HeaderStorage) generateTOC() ([]byte, error) {

	return serve(storage.Headers, "<div class=\"toc\">{{ range . }}<p class=\"level{{ printf \"%d\" .Level }}\">{{ range .Number }}{{ printf \"%d\" . }}.{{ end }} <a href=\"#{{ .ID }}\">{{ printf \"%s\" .Text }}</a></p>{{ end }}</div>")
}
