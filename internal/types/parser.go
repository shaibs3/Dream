package types

// ProcessEntry is a marker interface for process entries (Linux, Windows, etc.)
type ProcessEntry interface{}

type Parser interface {
	Parse(raw string) ([]ProcessEntry, error)
}
