package util

// Importer is a data import interface
type Importer interface {
	ImportFolder(folder string) error
	ImportFile(file, output string) error
}
