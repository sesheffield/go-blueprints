package thesaurus

// Thesaurus contract
type Thesaurus interface {
	Synonyms(string) ([]string, error)
}
