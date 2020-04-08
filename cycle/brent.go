package cycle

func Brent(cb func(i int) int, x0 int) (int, int) {
	var λ, µ, power, tortoise, hare int

	// Main phase: search successive powers of two.
	power = 1
	λ = 1
	tortoise = x0
	hare = cb(x0) // f(x0) is the element/node next to x0.

	for tortoise != hare {
		if power == λ { // Time to start a new power of two.
			tortoise = hare
			power *= 2
			λ = 0
		}
		hare = cb(hare)
		λ++
	}

	// Find the position of the first repetition of length λ.
	µ = 0
	tortoise, hare = x0, x0
	for i := 0; i < λ; i++ {
		// produces a list with the values 0,1,...,λ-1.
		hare = cb(hare)
		// The distance between hare and tortoise is now λ.
	}

	// The tortoise and the hare move at the same speed until they agree.
	for tortoise != hare {
		tortoise = cb(tortoise)
		hare = cb(hare)
		µ++
	}

	return λ, µ
}
