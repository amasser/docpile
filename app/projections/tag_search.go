package projections

import "github.com/sahilm/fuzzy"

type TagSearch struct {
	allDocs  []Document
	allTags  map[uint64]Tag
	synonyms map[string]bool

	searchText         string
	selectedTagIDs     []uint64
	candidateTagIDs    []uint64
	candidateTagValues []string
	fuzzyMatches       []fuzzy.Match
	results            []MatchingTag
}

func NewTagSearch(searchText string, selectedTagIDs []uint64) *TagSearch {
	return &TagSearch{searchText: searchText, selectedTagIDs: selectedTagIDs}
}

func (this *TagSearch) Search() []MatchingTag {
	this.gatherCandidates()
	this.conductSearch()
	this.renderResults()
	return this.results
}
func (this *TagSearch) gatherCandidates() {
	candidates := make(map[uint64]struct{})
	criteria := NewDocumentCriteria(nil, nil, nil, nil, this.selectedTagIDs)
	for _, doc := range this.allDocs {
		if !criteria.Match(doc) {
			continue
		}

		for _, tagID := range doc.Tags {
			candidates[tagID] = struct{}{}
		}
	}

	for _, selectedTagID := range this.selectedTagIDs {
		delete(candidates, selectedTagID)
	}

	this.candidateTagIDs = make([]uint64, 0, len(candidates))
	this.candidateTagValues = make([]string, 0, len(candidates))

	for tagID := range candidates {
		this.candidateTagIDs = append(this.candidateTagIDs, tagID)
		this.candidateTagValues = append(this.candidateTagValues, this.allTags[tagID].TagName)
	}
}
func (this *TagSearch) conductSearch() {
	this.fuzzyMatches = fuzzy.Find(this.searchText, this.candidateTagValues)
}
func (this *TagSearch) renderResults() {
	for _, fuzzyMatch := range this.fuzzyMatches {
		this.results = append(this.results, MatchingTag{
			TagID:   this.candidateTagIDs[fuzzyMatch.Index],
			TagText: fuzzyMatch.Str,
			Synonym: this.synonyms[fuzzyMatch.Str],
			Indexes: fuzzyMatch.MatchedIndexes,
		})
	}
}
