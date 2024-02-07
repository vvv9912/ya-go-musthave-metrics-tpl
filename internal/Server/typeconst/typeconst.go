package typeconst

type contextKey int64

// для передачи контекста внутри mw/handler
const (
	UserIDContextKey contextKey = 1
)
