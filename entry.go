package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// EntriesService servıce
type EntriesService service

//Entry model
type Entry struct {
	locale string
	Sys    *Sys                   `json:"sys"`
	Fields map[string]interface{} `json:"fields,omitempty"`
}

// GetVersion returns entity version
func (entry *Entry) GetVersion() int {
	version := 1
	if entry.Sys != nil {
		version = entry.Sys.Version
	}

	return version
}

// GetEntryKey returns the entry's keys
func (service *EntriesService) GetEntryKey(entry *Entry, key string) (*EntryField, error) {
	ef := EntryField{
		value: entry.Fields[key],
	}

	col, err := service.c.ContentTypes.List(entry.Sys.Space.Sys.ID).Next()
	if err != nil {
		return nil, err
	}

	for _, ct := range col.ToContentType() {
		if ct.Sys.ID != entry.Sys.ContentType.Sys.ID {
			continue
		}

		for _, field := range ct.Fields {
			if field.ID != key {
				continue
			}

			ef.dataType = field.Type
		}
	}

	return &ef, nil
}

// List returns entries collection
func (service *EntriesService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s%s/entries", spaceID, getEnvPath(service.c))
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Sync returns entries collection
func (service *EntriesService) Sync(spaceID string, initial bool, syncToken ...string) *Collection {
	path := fmt.Sprintf("/spaces/%s%s/sync", spaceID, getEnvPath(service.c))
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	if initial == true {
		col.Query.Initial("true")
	}
	if len(syncToken) == 1 {
		col.SyncToken = syncToken[0]
	}
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single entry
func (service *EntriesService) Get(spaceID, entryID string, locale ...string) (*Entry, error) {
	path := fmt.Sprintf("/spaces/%s%s/entries/%s", spaceID, getEnvPath(service.c), entryID)
	query := url.Values{}
	if len(locale) > 0 {
		query["locale"] = locale
	}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &Entry{}, err
	}

	var entry Entry
	if ok := service.c.do(req, &entry); ok != nil {
		return nil, err
	}

	return &entry, err
}

// Delete the entry
func (service *EntriesService) Delete(spaceID string, entryID string) error {
	path := fmt.Sprintf("/spaces/%s%s/entries/%s", spaceID, getEnvPath(service.c), entryID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}

// Upsert updates or creates a new entry
func (service *EntriesService) Upsert(spaceID string, entry *Entry) error {
	fieldsOnly := map[string]interface{}{
		"fields": entry.Fields,
	}

	bytesArray, err := json.Marshal(fieldsOnly)
	if err != nil {
		return err
	}

	// Creating/updating an entry requires a content type to be provided
	if entry.Sys.ContentType == nil {
		return fmt.Errorf("creating/updating an entry requires a content type")
	}

	var path string
	var method string

	if entry.Sys != nil && entry.Sys.ID != "" {
		path = fmt.Sprintf("/spaces/%s%s/entries/%s", spaceID, getEnvPath(service.c), entry.Sys.ID)
		method = "PUT"
	} else {
		path = fmt.Sprintf("/spaces/%s%s/entries", spaceID, getEnvPath(service.c))
		method = "POST"
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(entry.GetVersion()))
	req.Header.Set("X-Contentful-Content-Type", entry.Sys.ContentType.Sys.ID)

	return service.c.do(req, entry)
}

// Publish the entry
func (service *EntriesService) Publish(spaceID string, entry *Entry) error {
	path := fmt.Sprintf("/spaces/%s%s/entries/%s/published", spaceID, getEnvPath(service.c), entry.Sys.ID)
	method := "PUT"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(entry.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// Unpublish the entry
func (service *EntriesService) Unpublish(spaceID string, entry *Entry) error {
	path := fmt.Sprintf("/spaces/%s%s/entries/%s/published", spaceID, getEnvPath(service.c), entry.Sys.ID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(entry.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}
