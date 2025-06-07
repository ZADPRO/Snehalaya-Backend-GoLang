package resolver

type AdminResponse struct {
	Name string `json:"name"`
}

func ResolveAdmin(name string) AdminResponse {
	return AdminResponse{Name: name}
}
