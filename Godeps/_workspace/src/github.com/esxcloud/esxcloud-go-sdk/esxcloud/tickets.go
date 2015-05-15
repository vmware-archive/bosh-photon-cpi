package esxcloud

type ResourceTicketCreateSpec struct {
	Name   string          `json:"name"`
	Limits []QuotaLineItem `json:"limits"`
}

type ResourceTicket struct {
	Kind     string          `json:"kind"`
	Usage    []QuotaLineItem `json:"usage"`
	TenantId string          `json:"tenantId"`
	Name     string          `json:"name"`
	ID       string          `json:"id"`
	Limits   []QuotaLineItem `json:"limits"`
	Tags     []string        `json:"tags"`
	SelfLink string          `json:"selfLink"`
}

type ResourceList struct {
	Items []ResourceTicket `json:"items"`
}

type ResourceTicketReservation struct {
	Name   string          `json:"name"`
	Limits []QuotaLineItem `json:"limits"`
}
