package dbware

// gorm had packaged db features but more business functions is needed.
// dbware.go provides them.
// we dont wrap all the functions bcz that is an abstract layer over another.

/*
func (p *ProviderDB) getLatestID() (uint, error) {
	var maxid uint
	err := p.db.Table("users").
		Select("COALESCE(MAX(id), 0)").
		Scan(&maxid).Error
	return maxid, err
}*/
