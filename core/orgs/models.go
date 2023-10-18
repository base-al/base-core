package orgs

type OrgRequest struct {
	UserID    int    `json:"-"`
	Name      string `json:"name"`
	Size      string `json:"size"`
	TradeName string `json:"tradeName"`
	UIN       string `json:"uin"`
	VAT       string `json:"vat"`
}

type OrgResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	OrgSlug string `json:"orgSlug"`
}

type FindOrgRequest struct {
	UserID int `json:"-"`
}

type OrgWithRole struct {
	OrgID  int    `json:"orgId"`
	RoleID int    `json:"roleId"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
}
