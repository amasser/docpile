package projections

import "github.com/sahilm/fuzzy"

type TagSearch struct {
	allDocs      []Document
	tagIndexByID map[uint64]int
	allTags      []Tag

	searchText         string
	selectedTagIDs     []uint64
	candidateTagIDs    []uint64
	candidateTagValues []string
	fuzzyMatches       []fuzzy.Match
	results            []MatchingTag
}

func NewTagSearch(documents []Document, tagIndexByID map[uint64]int, allTags []Tag) *TagSearch {
	return &TagSearch{
		allDocs:      documents,
		tagIndexByID: tagIndexByID,
		allTags:      allTags,
	}
}

func (this *TagSearch) Search(searchText string, selectedTagIDs []uint64) []MatchingTag {
	this.searchText = searchText
	this.selectedTagIDs = selectedTagIDs

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

	this.removeSelectedTags(candidates)
	this.addCandidates(candidates)
}
func (this *TagSearch) removeSelectedTags(candidates map[uint64]struct{}) {
	for _, selectedTagID := range this.selectedTagIDs {
		delete(candidates, selectedTagID)
	}
}
func (this *TagSearch) addCandidates(candidates map[uint64]struct{}) {
	this.candidateTagIDs = make([]uint64, 0, len(candidates))
	this.candidateTagValues = make([]string, 0, len(candidates))
	for tagID := range candidates {
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
		tagID := this.candidateTagIDs[fuzzyMatch.Index]
		isSynonym := this.isSynonym(tagID, fuzzyMatch.Str)
		this.results = append(this.results, MatchingTag{
			TagID:   tagID,
			TagText: fuzzyMatch.Str,
			Synonym: isSynonym,
			Indexes: fuzzyMatch.MatchedIndexes,
		})
	}
}
func (this *TagSearch) isSynonym(id uint64, value string) bool {
	return this.getTag(id).TagName != value
}

func (this *TagSearch) getTag(id uint64) Tag {
	return this.allTags[this.tagIndexByID[id]]
}
