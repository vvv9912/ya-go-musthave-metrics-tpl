// Пакет typeconst c константами.
package typeconst

// Новый тип данных для передачи по контексту.
type contextKey int64

// Значение для передачи контекста внутри mw/handler.
const (
	UserIDContextKey contextKey = 1
)
