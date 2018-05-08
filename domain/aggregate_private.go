package domain

import (
	"fmt"
	"strings"

	"bitbucket.org/jonathanoliver/docpile/events"
)

type managedKey events.SHA256Hash
type cloudKey string

func newCloudAssetKey(provider, resource string) cloudKey {
	return cloudKey(fmt.Sprintf("%s.%s", strings.ToLower(provider), resource))
}
