package request

type CreateSkillRequest struct {
	ProfileCode int    `param:"profileCode" validate:"required"`
	Skill       string `json:"skill"`
	Level       string `json:"level"`
}
