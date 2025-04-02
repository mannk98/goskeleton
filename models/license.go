package models

// License represents a software license agreement, containing the Name of
// the license, its possible matches (on the command line as given to cobra),
// the header to be used with each file on the file's creating, and the text
// of the license
type License struct {
	Name            string   // The type of license in use
	PossibleMatches []string // Similar names to guess
	Text            string   // License text data
	Header          string   // License header for source files
}
