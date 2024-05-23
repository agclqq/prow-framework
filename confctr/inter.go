package confctr

type Val struct {
	Value string
}
type WatchType string

const (
	WatchTypeKey       WatchType = "key"
	WatchTypeKeyPrefix WatchType = "keyprefix"
	WatchTypeService   WatchType = "service"
	WatchTypeNode      WatchType = "node"
	WatchTypeEvent     WatchType = "event"
)

type WatchOption struct {
	// used for consul
	Datacenter string
	Token      string
	Type       WatchType
	Args       []string
	Tag        []string
}
type CallBack func(eventName string, key string, val string)

// CC config center interface
type CC interface {
	// Origin get the underlying config center implementation.
	Origin() interface{}
	// Get Retrieve an item from the config center by key.
	Get(key string) ([]Val, error)
	// Create a value to config center
	Create(key, value string) error
	// Update a value to config center
	Update(key, value string) error
	// Delete a value from config center
	Delete(key string) error
	// Watch a value from config center
	Watch(key string, f CallBack, option ...*WatchOption) error
}
