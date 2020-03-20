package db

const (
	UptimeLookBackDuration = 20
)

func (b *BandDB) AddValidator(
	operatorAddress string,
	consensusAddress string,
) error {
	return b.tx.Create(&ValidatorStatus{
		OperatorAddress:  operatorAddress,
		ConsensusAddress: consensusAddress,
	}).Error
}

func (b *BandDB) UpdateValidatorUpTime(
	consensusAddress string,
	height int64,
	voted bool,
) error {
	err := b.tx.Create(&ValidatorVote{
		ConsensusAddress: consensusAddress,
		BlockHeight:      height,
		Voted:            voted,
	}).Error

	if err != nil {
		return err
	}

	var validator ValidatorStatus
	err = b.tx.Where(ValidatorStatus{ConsensusAddress: consensusAddress}).First(&validator).Error
	if err != nil {
		return err
	}

	validator.ElectedCount++
	if voted {
		validator.VotedCount++
	} else {
		validator.MissedCount++
	}

	// Find old vote
	if height > UptimeLookBackDuration {
		var vote ValidatorVote
		err = b.tx.Where(ValidatorVote{
			ConsensusAddress: consensusAddress,
			BlockHeight:      height - UptimeLookBackDuration,
		}).First(&vote).Error
		if err == nil {
			validator.ElectedCount--
			if vote.Voted {
				validator.VotedCount--
			} else {
				validator.MissedCount--
			}
			b.tx.Delete(&vote)
		}
	}
	b.tx.Save(&validator)
	return nil
}

func (b *BandDB) ClearOldVotes(currentHeight int64) error {
	if currentHeight > UptimeLookBackDuration {
		return b.tx.Delete(
			ValidatorVote{},
			"block_height <= ?",
			currentHeight-UptimeLookBackDuration,
		).Error
	}
	return nil
}
