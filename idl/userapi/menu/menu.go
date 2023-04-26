package menu

type Menu struct {
	Tabs []Tab `json:"tabs"` // 菜单tab
}

type Tab struct {
	Category string `json:"category"` // 菜单分类：角色/工具
	Items    []Item `json:"items"`    // 菜单项
}

type Item struct {
	Type string `json:"type"` // 类型：chat/image/audio
	Name string `json:"name"` // 名称
}
