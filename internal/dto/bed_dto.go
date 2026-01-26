package dto

import "time"

// BedResponse represents bed details
type BedResponse struct {
	ID         uint      `json:"id"`
	BedNumber  string    `json:"bed_number"`
	Department string    `json:"department"`
	Ward       string    `json:"ward"`
	BedType    string    `json:"bed_type"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// BedListItem represents simplified bed for list view
type BedListItem struct {
	ID         uint   `json:"id"`
	BedNumber  string `json:"bed_number"`
	Department string `json:"department"`
	Ward       string `json:"ward"`
	BedType    string `json:"bed_type"`
	Status     string `json:"status"`
}
