package typeconst

type contextKey uint64

// для передачи контекста внутри mw/handler
const (
	UserIDContextKey contextKey = 1
)
