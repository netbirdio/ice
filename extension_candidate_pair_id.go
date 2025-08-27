package ice

import (
	"fmt"
	"strings"
)

const (
	// used in RDP for candidate ID extension
	extensionKeyCandidateID = "cid"

	candidateIDPrefix = "candidate:"
)

func candidateIDFromExtensions(extensions []CandidateExtension) string {
	for _, ext := range extensions {
		if ext.Key == extensionKeyCandidateID {
			return fmt.Sprintf("candidate:%s", ext.Value)
		}
	}

	return ""
}

func newCandidateIDExtension(candidateID string) CandidateExtension {
	return CandidateExtension{
		Key:   extensionKeyCandidateID,
		Value: strings.TrimPrefix(candidateID, candidateIDPrefix),
	}
}
