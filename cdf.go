package reedsolomon

import (
	"math"
)

// ReedsolomonPL calculate the actual-packet-loss rate after RS encoding, when the original-packet-loss rate is pl
func ReedsolomonPL(groupsize, datasize int, pl float64) float64 {
	return receivedPacketLossRate(datasize, groupsize-datasize, pl)
}

/*
  ====================== deepseek generate ============================
*/

func receivedPacketLossRate(datasize, paritysize int, pl float64) float64 {
	n := datasize
	k := paritysize
	if n == 0 {
		return 0.0
	}
	if pl <= 0.0 {
		return 0.0
	}
	if pl >= 1.0 {
		return 1.0
	}

	dataLossProbs := binomialProbabilities(n, pl)

	parityLossProbs := binomialProbabilities(k, pl)
	suffixSum := make([]float64, k+2)

	suffixSum[k+1] = 0.0
	for y := k; y >= 0; y-- {
		suffixSum[y] = suffixSum[y+1] + parityLossProbs[y]
	}

	expectedLoss := 0.0

	currentProb := dataLossProbs[0]
	for x := 0; x <= n; x++ {
		if x > 0 {
			currentProb = dataLossProbs[x]
		}

		if x <= k {
			threshold := k - x + 1
			var probRecoveryFail float64

			if threshold > k {
				probRecoveryFail = 0.0
			} else {
				probRecoveryFail = suffixSum[threshold]
			}

			expectedLoss += float64(x) * currentProb * probRecoveryFail
		} else {
			expectedLoss += float64(x) * currentProb
		}
	}

	lossRate := expectedLoss / float64(n)
	return math.Max(0.0, math.Min(lossRate, 1.0))
}

func binomialProbabilities(n int, p float64) []float64 {
	probs := make([]float64, n+1)
	if n == 0 {
		probs[0] = 1.0
		return probs
	}

	probs[0] = math.Pow(1-p, float64(n))
	if p == 0.0 {
		return probs
	}

	inv1MinusP := 1.0 / (1 - p)
	for x := 1; x <= n; x++ {
		// P(X=x) = P(X=x-1) * (n-x+1)/x * p/(1-p)
		probs[x] = probs[x-1] * float64(n-x+1) / float64(x) * p * inv1MinusP
	}
	return probs
}
