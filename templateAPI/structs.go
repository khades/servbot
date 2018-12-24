package templateAPI


type templatesResponse struct {
	Templates []TemplateInfo `json:"list"`
	IsMod     bool                  `json:"isMod"`
}

type templatePushRequest struct {
	Template string `json:"get"`
}

type aliasToRequest struct {
	AliasTo string `json:"aliasTo"`
}