package projections

import (
	"testing"

	"bitbucket.org/jonathanoliver/docpile/app/events"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestProjectorFixture(t *testing.T) {
	gunit.Run(new(ProjectorFixture), t)
}

type ProjectorFixture struct {
	*gunit.Fixture

	projector *Projector
}

func (this *ProjectorFixture) Setup() {
	this.projector = NewProjector()
}

func (this *ProjectorFixture) TestSearchTagsNotAssociatedWithDocumentNotReturned() {
	this.addTag(1, "a")

	result := this.search("a")

	this.So(result, should.BeEmpty)
}
func (this *ProjectorFixture) TestSearchReturnsMatchingTag() {
	this.addTag(1, "a")
	this.defineDocument(1)

	result := this.search("a")

	this.So(result, should.Resemble, []MatchingTag{
		{TagID: 1, TagText: "a", Synonym: false, Indexes: []int{0}},
	})
}
func (this *ProjectorFixture) TestSearchFindSynonyms() {
	this.addTag(1, "a")
	this.defineSynonym(1, "b")
	this.defineDocument(1)

	result := this.search("b")

	this.So(result, should.Resemble, []MatchingTag{
		{TagID: 1, TagText: "b", Synonym: true, Indexes: []int{0}},
	})
}
func (this *ProjectorFixture) TestSearchIgnoresTagsAlreadyProvided() {
	this.addTag(1, "aa")
	this.addTag(2, "ab")
	this.defineDocument(1, 2)

	result := this.search("a", 1)

	this.So(result, should.Resemble, []MatchingTag{
		{TagID: 2, TagText: "ab", Synonym: false, Indexes: []int{0}},
	})
}
func (this *ProjectorFixture) TestSearchIgnoresSynonymsAlreadyProvided() {
	this.addTag(1, "a")
	this.defineSynonym(1, "ab")
	this.addTag(2, "c")
	this.defineSynonym(2, "ad")
	this.defineDocument(1, 2)

	result := this.search("a", 1)

	this.So(result, should.Resemble, []MatchingTag{
		{TagID: 2, TagText: "ad", Synonym: true, Indexes: []int{0}},
	})
}
func (this *ProjectorFixture) TestSearchReturnsTagNameAndSynonym() {
	this.addTag(1, "ab")
	this.defineSynonym(1, "ac")
	this.defineDocument(1)

	result := this.search("a")

	this.So(result, should.Resemble, []MatchingTag{
		{TagID: 1, TagText: "ac", Synonym: true, Indexes: []int{0}},
		{TagID: 1, TagText: "ab", Synonym: false, Indexes: []int{0}},
	})
}

func (this *ProjectorFixture) addTag(tagID uint64, name string) {
	this.apply(events.TagAdded{TagID: tagID, TagName: name})
}
func (this *ProjectorFixture) defineSynonym(tagID uint64, name string) {
	this.apply(events.TagSynonymDefined{TagID: tagID, Synonym: name})
}
func (this *ProjectorFixture) defineDocument(tagIDs ...uint64) {
	this.apply(events.DocumentDefined{Tags: tagIDs})
}

func (this *ProjectorFixture) apply(messages ...interface{}) {
	this.projector.Apply(messages)
}
func (this *ProjectorFixture) search(text string, tags ...uint64) []MatchingTag {
	return this.projector.MatchingTags.Search(text, tags)
}
