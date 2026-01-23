package di

import "context"

func Get(ctx context.Context, key string) any {
	ctn, ok := ctx.Value(containerKey).(*container)
	if !ok {
		return nil
	}

	return ctn.Get(key)
}
