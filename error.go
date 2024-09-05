package main

import "fmt"

type ErrProfileNotFound string

func (e ErrProfileNotFound) Error() string {
	return fmt.Sprintf("store %q not found", string(e))
}
