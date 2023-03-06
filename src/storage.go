package main

import (
	"fmt"
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

func (storage *IdStorage) push(new_id string) error {

	_, already_exists := storage.Storage[new_id]

	if already_exists {
		return fmt.Errorf("ID '%s' declared twice!", new_id)
	}

	storage.Storage[new_id] = struct{}{}

	return nil
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

func (images *ImageStorage) copy() error {
	return nil
}

func (storage *ImageStorage) push(new_link string) {
	storage.Storage[new_link] = struct{}{}
}

type HeaderStorage []*HeaderInfo

type HeaderInfo struct {
	Text  []byte
	Level int
}

func (storage *HeaderStorage) push(level int, header []byte) error {

	if level < 1 {
		return fmt.Errorf("Invalid level '%d', levels <1 not allowed", level)
	}

	header_info := &HeaderInfo{
		Text:  header,
		Level: level,
	}

	*storage = append(*storage, header_info)

	return nil
}

func checkAllLinksValid(links *LinkStorage, ids *IdStorage) error {
	return nil
}
