package telegram

import "context"

type Middleware func(Handler) Handler

func chain(
	handler Handler,
	middleware ...Middleware,
) Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}

	return handler
}

// WithContext позволяет положить значение
// в обычный context.Context.
//
// Например:
//
//	handler := func(t telegram.Context) error {
//		userID := t.Context.Value(userIDKey{})
//		...
//	}
func WithContext(value any) Middleware {
	return func(next Handler) Handler {
		return func(t Context) error {
			t.Context = context.WithValue(
				t.Context,
				contextKey{},
				value,
			)

			return next(t)
		}
	}
}

type contextKey struct{}
