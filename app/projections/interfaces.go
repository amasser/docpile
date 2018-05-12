package projections

type DocumentSpecification interface {
	IsSatisfiedBy(Document) bool
}
