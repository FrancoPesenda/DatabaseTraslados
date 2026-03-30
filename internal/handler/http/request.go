package http

type createAdminRequest struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}
