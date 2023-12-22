// Package entity contains the main domain entities (data models) of the
// core business logic
package entity

// Link is the main entity which represents a shortened link
type Link struct {
	// Alias is by specification a 10-character string which is used in the
	// short versin of the link
	Alias string `json:"alias"`

	// Original represents the original URL which has a corresponding short
	// version
	Original string `json:"original"`
}
