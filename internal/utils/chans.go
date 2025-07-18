package utils

import "sync"

func FanIn[T any](cs ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)

	output := func(c <-chan T) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func FanOut[T any](input <-chan T, count int) []<-chan T {
	outputs := make([]chan T, count)
	for i := 0; i < count; i++ {
		outputs[i] = make(chan T)
	}

	go func() {
		defer func() {
			for _, out := range outputs {
				close(out)
			}
		}()

		for item := range input {
			var wg sync.WaitGroup
			wg.Add(len(outputs))

			// Рассылаем сообщение во все выходные каналы
			for _, out := range outputs {
				go func(ch chan<- T) {
					defer wg.Done()
					ch <- item
				}(out)
			}

			wg.Wait()
		}
	}()

	// Конвертируем в read-only каналы
	result := make([]<-chan T, count)
	for i, ch := range outputs {
		result[i] = ch
	}

	return result
}
