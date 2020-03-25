package main

import "math"

/*
def mel(sr, n_fft, n_mels=128, fmin=0.0, fmax=None, htk=False,
        norm='slaney', dtype=np.float32):
    """Create a Filterbank matrix to combine FFT bins into Mel-frequency bins
*/

type Thing struct {
	sr float64
	nfft float64
	fmax float64
	nmels int
}

func (t *Thing) mel() {
	if t.fmax == 0.0 {
		t.fmax = t.sr / 2
	}
//	    weights = np.zeros((n_mels, int(1 + n_fft // 2)), dtype=dtype)
  math.Floor(1.0 + t.nfft / 2.0)
}
