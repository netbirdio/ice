package ice

import (
	"fmt"
	"strings"

	"github.com/pion/stun/v2"
)

const (
	AttrCandidatePairID stun.AttrType = 0x8100 // Custom attribute for candidate ID
)

type CandidatePairID struct {
	source      string
	destination string
}

func CandidatePairIDFromSTUN(msg *stun.Message) (*CandidatePairID, bool, error) {
	candidatePairIDBytes, err := msg.Get(AttrCandidatePairID)
	if err != nil {
		return nil, false, nil
	}
	candidatePairID, err := ParseCandidatePairID(candidatePairIDBytes)
	if err != nil {
		return nil, false, err
	}
	return &candidatePairID, true, nil
}

// NewCandidatePairID creates a CandidatePairID from local and remote candidates.
func NewCandidatePairID(local, remote Candidate) CandidatePairID {
	localID := strings.TrimPrefix(local.ID(), "candidate:")
	remoteID := strings.TrimPrefix(remote.ID(), "candidate:")

	return CandidatePairID{
		source:      localID,
		destination: remoteID,
	}
}

// ParseCandidatePairID parses a CandidatePairID from its string representation.
func ParseCandidatePairID(id []byte) (CandidatePairID, error) {
	parts := strings.SplitN(string(id), ":", 2)
	if len(parts) != 2 {
		return CandidatePairID{}, fmt.Errorf("invalid candidates ID format: %s", id)
	}
	return CandidatePairID{
		source:      parts[0],
		destination: parts[1],
	}, nil
}

func (cp CandidatePairID) String() string {
	return fmt.Sprintf("%s:%s", cp.source, cp.destination)
}

func (cp CandidatePairID) SourceCandidateID() string {
	return fmt.Sprintf("candidate:%s", cp.source)
}

func (cp CandidatePairID) TargetCandidateID() string {
	return fmt.Sprintf("candidate:%s", cp.destination)
}
