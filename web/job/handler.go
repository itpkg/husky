package job

//Handler job handler
type Handler func(...interface{}) error

var handlers = make(map[string]Handler)

//Register register job worker
func Register(n string, h Handler) {
	handlers[n] = h
}
