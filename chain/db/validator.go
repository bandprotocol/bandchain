package db

func (b *BandDB) AddValidator(validatorAddress string) error {
	return b.tx.Create(ValidatorStatus{
		ValidatorAddress: validatorAddress,
	}).Error
}

func (b *BandDB) AddNewValidatorVote(validatorAddress string, height int64, voted bool) error {
	return b.tx.Create(ValidatorVote{
		ValidatorAddress: validatorAddress,
		BlockHeight:      height,
		Voted:            voted,
	}).Error
}
