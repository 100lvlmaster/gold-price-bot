package models

import "gorm.io/datatypes"

type Setting struct {
	ID    uint           `gorm:"primaryKey" json:"id"`
	Name  string         `json:"name"`
	Value datatypes.JSON `json:"value"`
}
