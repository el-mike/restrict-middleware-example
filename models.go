package main

type User struct {
	ID    string
	Roles []string
}

// Subject interface implementation.
func (u *User) GetRoles() []string {
	return u.Roles
}

// Resource interface implementation.
func (u *User) GetResourceName() string {
	return "User"
}

type Conversation struct {
	ID        string
	CreatedBy string
}

// Resource interface implementation.
func (c *Conversation) GetResourceName() string {
	return "Conversation"
}
