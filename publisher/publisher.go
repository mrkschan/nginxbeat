package publisher

// Publisher publishes Nginx status via libbeat.
type Publisher interface {
	Publish(s map[string]interface{})
}
