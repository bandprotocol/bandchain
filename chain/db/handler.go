package db

func (b *BandDB) handleMessageEvent(action string) {
	b.tx.Create(&Event{Name: action})
}
