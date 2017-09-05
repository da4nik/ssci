package types

// Notification from source controls
type Notification struct {
	Name     string
	CloneURL string
}

// Notificatable interface to convert to Notification
type Notificatable interface {
	Notification() Notification
}
