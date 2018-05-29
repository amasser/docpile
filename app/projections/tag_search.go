package projections

import "github.com/sahilm/fuzzy"

type TagSearch struct {
	allDocs      []Document
	tagIndexByID map[uint64]int
	allTags      []Tag

	searchText     string
	selectedTagIDs []uint64

	candidates         map[uint64]struct{}
	candidateTagIDs    []uint64
	candidateTagValues []string

	fuzzyMatches []fuzzy.Match
	results      []MatchingTag
}

func NewTagSearch(documents []Document, tagIndexByID map[uint64]int, allTags []Tag) *TagSearch {
	return &TagSearch{
		allDocs:      documents,
		tagIndexByID: tagIndexByID,
		allTags:      allTags,
		candidates:   make(map[uint64]struct{}),
		results:      []MatchingTag{},
	}
}

func (this *TagSearch) Search(searchText string, selectedTagIDs []uint64) []MatchingTag {
	this.searchText = searchText
	this.selectedTagIDs = selectedTagIDs

	this.gatherCandidates()
	this.removeSelectedTags()
	this.addCandidates()

	this.conductSearch()
	this.renderResults()
	return this.results
}

func (this *TagSearch) gatherCandidates() {
	criteria := NewDocumentCriteria(nil, nil, nil, nil, this.selectedTagIDs)
	for _, doc := range this.allDocs {
		if criteria.Match(doc) {
			for _, tagID := range doc.Tags {
				this.candidates[tagID] = struct{}{}
			}
		}
	}
}
func (this *TagSearch) removeSelectedTags() {
	for _, selectedTagID := range this.selectedTagIDs {
		delete(this.candidates, selectedTagID)
	}
}
func (this *TagSearch) addCandidates() {
	this.candidateTagIDs = make([]uint64, 0, len(this.candidates))
	this.candidateTagValues = make([]string, 0, len(this.candidates))

	for tagID := range this.candidates {
		this.addCandidate(tagID)
	}
}
func (this *TagSearch) addCandidate(tagID uint64) {
	tag := this.getTag(tagID)
	this.candidateTagIDs = append(this.candidateTagIDs, tag.TagID)
	this.candidateTagValues = append(this.candidateTagValues, tag.TagName)
	for synonym := range tag.Synonyms {
		this.candidateTagValues = append(this.candidateTagValues, synonym)
	}
}

func (this *TagSearch) conductSearch() {
	this.fuzzyMatches = fuzzy.Find(this.searchText, this.candidateTagValues)
}

func (this *TagSearch) renderResults() {
	for _, fuzzyMatch := range this.fuzzyMatches {
		this.addResult(fuzzyMatch)
	}
}
func (this *TagSearch) addResult(value fuzzy.Match) {
	tagID := this.candidateTagIDs[value.Index]
	isSynonym := this.isSynonym(tagID, value.Str)
	this.results = append(this.results, MatchingTag{
		TagID:   tagID,
		TagText: value.Str,
		Synonym: isSynonym,
		Indexes: value.MatchedIndexes,
	})
}
func (this *TagSearch) isSynonym(id uint64, value string) bool {
	return this.getTag(id).TagName != value
}
func (this *TagSearch) getTag(id uint64) Tag {
	return this.allTags[this.tagIndexByID[id]]
}
