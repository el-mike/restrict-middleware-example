package main

type User struct {
	ID   string
	Role string
}

// Subject interface implementation.
func (u *User) GetRole() string {
	return u.Role
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
