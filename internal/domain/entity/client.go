package entity

// Client adalah interface kontrak koneksi client yang aktif di server
type Client interface {
	Username() string
	SendChannel() chan Message
	Close()
}
