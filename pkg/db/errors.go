package db

import "fmt"

type ErrAliasExists struct {
	alias string
}

func (e ErrAliasExists) Error() string {
	return fmt.Sprintf("alias %v is already taken", e.alias)
}

type ErrNotExists struct {
	alias string
}

func (e ErrNotExists) Error() string {
	return fmt.Sprintf("alias %v does not exist", e.alias)
}
