package tdo

type Profile struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Avatar    string `json:"avatar"`
	Role      string `json:"role"`
	Provider  string `json:"provider"`
	CreatedAt string `json:"created_at"`
	DeletedAt string `json:"deleted_at"`
}

func NewProfile(id, email, fullName, avatar, role, provider, createdAt, deletedAt string) Profile {
	return Profile{
		ID:        id,
		Email:     email,
		FullName:  fullName,
		Avatar:    avatar,
		Role:      role,
		Provider:  provider,
		CreatedAt: createdAt,
		DeletedAt: deletedAt,
	}
}
