package resources

type ResourceRequest struct {
	NotResource bool `json:"not"`
	Resource    string
}
