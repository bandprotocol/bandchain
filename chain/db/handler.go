package db

func (b *BandDB) handleMessageEvent(action string) {
	b.db.Create(&Event{Name: action})
}
