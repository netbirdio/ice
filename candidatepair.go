// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package ice

import (
	"fmt"
	"time"

	"github.com/pion/stun/v2"
)

func newCandidatePair(local, remote Candidate, controlling bool) *CandidatePair {
	return &CandidatePair{
		iceRoleControlling: controlling,
		Remote:             remote,
		Local:              local,
		state:              CandidatePairStateWaiting,
	}
}

type TransactionID [stun.TransactionIDSize]byte

// CandidatePair is a combination of a
// local and remote candidate
type CandidatePair struct {
	iceRoleControlling       bool
	Remote                   Candidate
	Local                    Candidate
	latency                  time.Duration
	lastBindingRequest       time.Time
	lastBindingTransactionID TransactionID
	bindingRequestCount      uint16
	state                    CandidatePairState
	nominated                bool
	nominateOnBindingSuccess bool
}

func (p *CandidatePair) String() string {
	if p == nil {
		return ""
	}

	return fmt.Sprintf("prio %d (local, prio %d) %s <-> %s (remote, prio %d), state: %s, nominated: %v, nominateOnBindingSuccess: %v",
		p.priority(), p.Local.Priority(), p.Local, p.Remote, p.Remote.Priority(), p.state, p.nominated, p.nominateOnBindingSuccess)
}

func (p *CandidatePair) equal(other *CandidatePair) bool {
	if p == nil && other == nil {
		return true
	}
	if p == nil || other == nil {
		return false
	}
	return p.Local.Equal(other.Local) && p.Remote.Equal(other.Remote)
}

// RFC 5245 - 5.7.2.  Computing Pair Priority and Ordering Pairs
// Let G be the priority for the candidate provided by the controlling
// agent.  Let D be the priority for the candidate provided by the
// controlled agent.
// pair priority = 2^32*MIN(G,D) + 2*MAX(G,D) + (G>D?1:0)
func (p *CandidatePair) priority() uint64 {
	var g, d uint32
	if p.iceRoleControlling {
		g = p.Local.Priority()
		d = p.Remote.Priority()
	} else {
		g = p.Remote.Priority()
		d = p.Local.Priority()
	}

	// Just implement these here rather
	// than fooling around with the math package
	min := func(x, y uint32) uint64 {
		if x < y {
			return uint64(x)
		}
		return uint64(y)
	}
	max := func(x, y uint32) uint64 {
		if x > y {
			return uint64(x)
		}
		return uint64(y)
	}
	cmp := func(x, y uint32) uint64 {
		if x > y {
			return uint64(1)
		}
		return uint64(0)
	}

	// 1<<32 overflows uint32; and if both g && d are
	// maxUint32, this result would overflow uint64
	return (1<<32-1)*min(g, d) + 2*max(g, d) + cmp(g, d)
}

func (p *CandidatePair) Write(b []byte) (int, error) {
	return p.Local.writeTo(b, p.Remote)
}

func (a *Agent) sendSTUN(msg *stun.Message, local, remote Candidate) {
	_, err := local.writeTo(msg.Raw, remote)
	if err != nil {
		a.log.Tracef("Failed to send STUN message: %s", err)
	}
}

func (p *CandidatePair) markBindingRequest(transactionID TransactionID) {
	p.lastBindingRequest = time.Now()
	p.lastBindingTransactionID = transactionID
}

func (p *CandidatePair) markBindingResponse(transactionID TransactionID) bool {
	if p.lastBindingRequest.IsZero() || transactionID != p.lastBindingTransactionID {
		return false
	}

	p.latency = time.Since(p.lastBindingRequest)
	return true
}

func (p *CandidatePair) Latency() time.Duration {
	return p.latency
}
