package models

import "time"

// JournalEntry represents a journal entry in Netbox
type JournalEntry struct {
	ID          int       `json:"id,omitempty"`
	AssignedObj string    `json:"assigned_object_type,omitempty"`
	AssignedID  int       `json:"assigned_object_id,omitempty"`
	Created     time.Time `json:"created,omitempty"`
	CreatedBy   string    `json:"created_by,omitempty"`
	Kind        string    `json:"kind,omitempty"`
	Comments    string    `json:"comments"`
}
