package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type UnitGroup struct {
	IDs []string
	Key ebiten.Key
}

func NewUnitGroup(IDs []string, key ebiten.Key) *UnitGroup {
	return &UnitGroup{
		IDs: IDs,
		Key: key,
	}
}

type UnitGroupManager struct {
	groups []*UnitGroup
}

func (u *UnitGroupManager) AddGroup(IDs []string, key ebiten.Key) {
	newGroup := NewUnitGroup(IDs, key)
	u.groups = append(u.groups, newGroup)
}
