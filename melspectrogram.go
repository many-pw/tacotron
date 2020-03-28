package main

import "github.com/aclements/go-moremath/vec"

import "fmt"

func melspectrogram(y []float32) {
	getWindow(1100, true)
	//D := stft(y)
	//S = amp_to_db(linear_to_mel(np.abs(D)))
	//return normalize(S)
}

func getWindow(length int, fftbins bool) {
	fac := vec.Linspace(-3.14, 3.14, length)
	fmt.Println(len(fac))
	// 1100
	//#general_hamming(M, 0.5, True)
	//def general_hamming(M, alpha, sym=True):
	//  return general_cosine(M, [alpha, 1. - alpha], sym)
}
