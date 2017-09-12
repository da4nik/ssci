package types

// Notification from source controls
type Notification struct {
	Name     string
	FullName string
	CloneURL string
	Version  string
}

// Notificatable interface to convert to Notification
type Notificatable interface {
	Notification() Notification
}
