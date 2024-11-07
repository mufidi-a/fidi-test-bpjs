package response

type SkillResponse struct {
	Id    int    `json:"id"`
	Skill string `json:"skill"`
	Level string `json:"level"`
}

type SkillList struct {
	Data []*SkillResponse `json:"data"`
}
