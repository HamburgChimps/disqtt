package router

import (
	"github.com/eclipse/paho.mqtt.golang/paho"
)

// Router is a trie based MQTT message router which can be used to
// dispatch messages to different handler functions.
type Router struct {
	trie           *Trie
	DefaultHandler HandlerFunc
}

// NewRouter returns a Router instance.
func NewRouter(defaultHandler HandlerFunc, opts ...Options) *Router {
	return &Router{trie: NewTrie(opts...), DefaultHandler: defaultHandler}
}

// RegisterHandler takes a string of the topic, and a MessageHandler
// to be invoked when Publishes are received that match that topic
func (r *Router) RegisterHandler(pattern string, handler HandlerFunc) {
	r.trie.Define(pattern).Handle(handler)
}

// UnregisterHandler takes a string of the topic to remove
// MessageHandlers
func (r *Router) UnregisterHandler(string) {

}

// Route sends messages to registered handlers
func (r *Router) Route(m *paho.Publish) {
	var handler HandlerFunc
	topic := string(m.Topic)
	res := r.trie.Match(topic)

	if res.Node == nil {
		handler = r.DefaultHandler
	} else {
		handler = res.Node.GetHandler()
	}

	handler(m)
}
