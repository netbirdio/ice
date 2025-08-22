package ice

const (
	// used in RDP for candidate ID extension
	extensionKeyCandidateID = "cid"
)

func candidateIDFromExtensions(extensions []CandidateExtension) string {
	for _, ext := range extensions {
		if ext.Key == extensionKeyCandidateID {
			return ext.Value
		}
	}

	return ""
}

func newCandidateIDExtension(candidateID string) CandidateExtension {
	return CandidateExtension{
		Key:   extensionKeyCandidateID,
		Value: candidateID,
	}
}
