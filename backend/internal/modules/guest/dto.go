package guest

type CreateIdentityTypeRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateIdentityTypeRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateGuestRequest struct {
	Title          string  `json:"title"`
	FirstName      string  `json:"first_name" binding:"required"`
	LastName       string  `json:"last_name"`
	IdentityTypeId *string `json:"identity_type_id"`
	IdentityNumber string  `json:"identity_number"`
	IdentityImage  string  `json:"identity_image"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	Nationality    string  `json:"nationality"`
	Address        string  `json:"address"`
	GuestType      string  `json:"guest_type"`
}

type UpdateGuestRequest struct {
	Title          string  `json:"title"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	IdentityTypeId *string `json:"identity_type_id"`
	IdentityNumber string  `json:"identity_number"`
	IdentityImage  string  `json:"identity_image"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	Nationality    string  `json:"nationality"`
	Address        string  `json:"address"`
	GuestType      string  `json:"guest_type"`
}
