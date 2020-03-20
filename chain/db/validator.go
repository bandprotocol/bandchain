package db

func (b *BandDB) AddValidator(
	operatorAddress string,
	consensusAddress string,
) error {
	return b.tx.Create(&ValidatorStatus{
		OperatorAddress:  operatorAddress,
		ConsensusAddress: consensusAddress,
	}).Error
}

func (b *BandDB) AddValidatorUpTime(
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

	b.tx.Save(&validator)
	return nil
}

func (b *BandDB) ClearOldVotes(currentHeight int64) error {
	uptimeLookBackDuration, err := b.GetUptimeLookBackDuration()
	if err != nil {
		return err
	}

	if currentHeight > uptimeLookBackDuration {
		var votes []ValidatorVote
		err := b.tx.Find(
			&votes,
			"block_height <= ?",
			currentHeight-uptimeLookBackDuration,
		).Error

		if err != nil {
			return err
		}
		for _, vote := range votes {
			var validator ValidatorStatus
			err = b.tx.Where(ValidatorStatus{ConsensusAddress: vote.ConsensusAddress}).First(&validator).Error
			if err == nil {
				validator.ElectedCount--
				if vote.Voted {
					validator.VotedCount--
				} else {
					validator.MissedCount--
				}
				b.tx.Save(&validator)
			}

		}
		return b.tx.Delete(
			ValidatorVote{},
			"block_height <= ?",
			currentHeight-uptimeLookBackDuration,
		).Error
	}
	return nil
}
