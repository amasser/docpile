package search

import "bitbucket.org/jonathanoliver/docpile/app/projections"

type DocumentSearcher struct {
	projection *projections.AllDocuments
}

func NewDocumentSearcher(projection *projections.AllDocuments) *DocumentSearcher {
	return &DocumentSearcher{projection: projection}
}

func (this *DocumentSearcher) Search(spec DocumentSpecification) (matching []projections.Document) {
	for _, document := range this.projection.List() {
		if spec.IsSatisfiedBy(document) {
			matching = append(matching, document)
		}
	}

	return matching
}
