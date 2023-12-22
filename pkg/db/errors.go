package db

import "fmt"

// ErrAliasExists represents an error in case when there is an attempt to save
// a short link with alias that already exists in the database
type ErrAliasExists struct {
	alias string
}

func (e ErrAliasExists) Error() string {
	return fmt.Sprintf("alias %v is already taken", e.alias)
}

// ErrNotExists represents an error in case when there's an attempt to fetch a
// short link which does not exist in the database
type ErrNotExists struct {
	alias string
}

func (e ErrNotExists) Error() string {
	return fmt.Sprintf("alias %v does not exist", e.alias)
}
