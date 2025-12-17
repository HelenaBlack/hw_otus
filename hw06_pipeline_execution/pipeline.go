package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// orDone возвращает канал, который пересылает значения из in,
// но закрывается сразу при закрытии done.
func orDone(done Bi, in In) In {
	if done == nil {
		return in
	}
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-done:
					return
				}
			}
		}
	}()
	return out
}

// ExecutePipeline передаёт данные через цепочку стадий.
// Между стадиями запускается форвардер, который при закрытии done
// быстро закрывает следующий канал и в фоне дренирует исходный,
// чтобы дать завершиться горутинам стадий без задержки возвращаемого канала.
func ExecutePipeline(in In, done Bi, stages ...Stage) Out {
	if len(stages) == 0 {
		return orDone(done, in)
	}

	cur := in
	for i, s := range stages {
		// не даём стадиям читать после закрытия done
		cur = orDone(done, cur)

		origOut := s(cur)

		// если последняя стадия — возвращаем её выход через форвардер
		next := make(Bi)
		go func(orig Out, out Bi, done Bi) {
			defer close(out)
			if done == nil {
				for v := range orig {
					out <- v
				}
				return
			}
			for {
				select {
				case <-done:
					// при done — сразу запускаем дрен и выходим
					go drain(orig)
					return
				case v, ok := <-orig:
					if !ok {
						return
					}
					select {
					case out <- v:
						// успешно переслали
					case <-done:
						// при done — запустить дрен и выйти, не блокируя
						go drain(orig)
						return
					}
				}
			}
		}(origOut, next, done)

		// если это последняя стадия — вернуть канал next
		if i == len(stages)-1 {
			return next
		}
		cur = next
	}

	// неиспользуемый путь, но нужен для компилятора
	return orDone(done, cur)
}

func drain(ch <-chan interface{}) {
	for v := range ch {
		_ = v
	}
}
