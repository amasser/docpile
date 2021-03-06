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

/*
As a tag is selected, the id is added to the tag search so that:
1. that tag (and synonyms) are no longer returned in the suggested result set
2. only tags/synonyms associated with documents containing the selected tags are returned

this has the effect of rapidly narrowing the window of available tags instead of making
the user swim through hundreds of useless tags that obviously don't apply and could never
be met as a search condition.

For example:
GET /search/tags?text=search+query&tag=previouslySelectedTagID1&tag=previouslySelectedTagID2

In this search above, "search query" would not return any tag names/synonyms for
previouslySelectedTagID1 and previouslySelectedTagID2. Further, "search query" would only match
names and synonyms for tags for documents that contain by previouslySelectedTagID1 and 2.
*/

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
		tag := this.getTag(tagID)
		this.addCandidate(tagID, tag.TagName)
		for synonym := range tag.Synonyms {
			this.addCandidate(tagID, synonym)
		}
	}
}
func (this *TagSearch) addCandidate(id uint64, value string) {
	this.candidateTagIDs = append(this.candidateTagIDs, id)
	this.candidateTagValues = append(this.candidateTagValues, value)
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
