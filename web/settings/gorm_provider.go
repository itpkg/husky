package settings

import (
	"github.com/itpkg/husky/web/crypto"
	"github.com/jinzhu/gorm"
)

//GormProvider provider of gorm
type GormProvider struct {
	Db  *gorm.DB
	Enc *crypto.Encryptor
	Ser *crypto.Serial
}

//Set set
func (p *GormProvider) Set(k string, v interface{}, f bool) error {
	buf, err := p.Ser.To(v)
	if err != nil {
		return err
	}
	if f {
		buf, err = p.Enc.Encode(buf)
		if err != nil {
			return err
		}
	}
	var m Model
	null := p.Db.Where("key = ?", k).First(&m).RecordNotFound()
	m.Key = k
	m.Val = buf
	m.Flag = f
	if null {
		err = p.Db.Create(&m).Error
	} else {
		err = p.Db.Save(&m).Error
	}
	return err
}

//Get get
func (p *GormProvider) Get(k string, v interface{}) error {
	var m Model
	err := p.Db.Where("key = ?", k).First(&m).Error
	if err != nil {
		return err
	}
	if m.Flag {
		if m.Val, err = p.Enc.Decode(m.Val); err != nil {
			return err
		}
	}
	return p.Ser.From(m.Val, v)
}
