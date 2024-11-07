package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Skill struct {
	bun.BaseModel `bun:"table:skill"`

	ProfileCode int       `bun:"profile_code"`
	Id          int       `bun:"id,pk,type:int,autoincrement"`
	Skill       string    `bun:"skill"`
	Level       string    `bun:"level"`
	CreatedAt   time.Time `bun:"created_at,default:current_timestamp"`
}

type SkillDTO struct {
	ProfileCode int       `json:"profileCode"`
	Id          int       `json:"id"`
	Skill       string    `json:"skill"`
	Level       string    `json:"level"`
	CreatedAt   time.Time `json:"createdAt"`
}
