package worker

import (
	"context"
)

func _(ctx context.Context, inputChannel chan any, receiver func(any)) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-inputChannel:
			receiver(message)
		}
	}
}
