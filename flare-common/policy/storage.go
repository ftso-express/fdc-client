package policy

import (
	"cmp"
	"fmt"
	"sort"
)

// Does not lock the structure, should be called from a function that does lock.
// We assume that the list is sorted by rewardEpochID and also by startVotingRoundID.
func (s *SigningPolicyStorage) findByVotingRoundID(votingRoundID uint32) *SigningPolicy {
	i, found := sort.Find(len(s.spList), func(i int) int {
		return cmp.Compare(votingRoundID, s.spList[i].startVotingRoundID)
	})
	if found {
		return s.spList[i]
	}
	if i == 0 {
		return nil
	}
	return s.spList[i-1]
}

func (s *SigningPolicyStorage) Add(sp *SigningPolicy) error {
	s.Lock()
	defer s.Unlock()

	if len(s.spList) > 0 {
		// check consistency, previous epoch should be already added
		if s.spList[len(s.spList)-1].rewardEpochID != sp.rewardEpochID-1 {
			return fmt.Errorf("missing signing policy for reward epoch ID %d", sp.rewardEpochID-1)
		}
		// should be sorted by voting round ID, should not happen
		if sp.startVotingRoundID < s.spList[len(s.spList)-1].startVotingRoundID {
			return fmt.Errorf("signing policy for reward epoch ID %d has larger start voting round ID than previous policy",
				sp.rewardEpochID)
		}
	}

	s.spList = append(s.spList, sp)
	return nil
}

// Return the signing policy for the voting round, or nil if not found.
// Also returns true if the policy is the last one or false otherwise.
func (s *SigningPolicyStorage) ForVotingRound(votingRoundID uint32) (*SigningPolicy, bool) {
	s.Lock()
	defer s.Unlock()

	sp := s.findByVotingRoundID(votingRoundID)
	if sp == nil {
		return nil, false
	}
	return sp, sp.rewardEpochID == s.spList[len(s.spList)-1].rewardEpochID
}

// Removes all signing policies with start voting round ID <= than the provided one.
// Returns the list of removed reward epoch ids.
func (s *SigningPolicyStorage) RemoveByVotingRound(votingRoundID uint32) []uint32 {
	s.Lock()
	defer s.Unlock()

	var removedRewardEpochIDs []uint32
	for len(s.spList) > 0 && s.spList[0].startVotingRoundID <= votingRoundID {
		removedRewardEpochIDs = append(removedRewardEpochIDs, uint32(s.spList[0].rewardEpochID))
		s.spList[0] = nil
		s.spList = s.spList[1:]
	}
	return removedRewardEpochIDs
}

// RemoveBefore removes all signing policies that ended strictly before votingRoundID.
// Returns the list of removed reward epoch ids.
func (s *SigningPolicyStorage) RemoveBefore(votingRoundID uint32) []uint32 {
	s.Lock()
	defer s.Unlock()

	var removedRewardEpochIDs []uint32
	for len(s.spList) > 1 && s.spList[1].startVotingRoundID < votingRoundID {
		removedRewardEpochIDs = append(removedRewardEpochIDs, uint32(s.spList[0].rewardEpochID))
		s.spList[0] = nil
		s.spList = s.spList[1:]
	}
	return removedRewardEpochIDs
}
