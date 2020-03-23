package db

func (b *BandDB) handleMessageEvent(action string) error {
	return b.tx.Create(&Event{Name: action}).Error
}
