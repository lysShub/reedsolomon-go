package reedsolomon

import (
	"math"
)

// ReedsolomonPL 计算在原始丢包率为pl情况下, rs编码后的实际丢包率
func ReedsolomonPL(groupsize, datasize int, pl float64) float64 {
	return receivedPacketLossRate(datasize, groupsize-datasize, pl)
}

/*
  ====================== 以下代码由deepseek生成 ============================
*/

// 计算接收方数据包丢包率
func receivedPacketLossRate(datasize, paritysize int, pl float64) float64 {
	n := datasize
	k := paritysize

	// 特殊情况处理
	if n == 0 {
		return 0.0
	}
	if pl <= 0.0 {
		return 0.0
	}
	if pl >= 1.0 {
		return 1.0
	}

	// 计算数据包丢失的概率分布
	dataLossProbs := binomialProbabilities(n, pl)

	// 计算校验包丢失的概率分布和后缀和
	parityLossProbs := binomialProbabilities(k, pl)
	suffixSum := make([]float64, k+2) // 后缀和数组，索引0..k+1

	// 初始化后缀和：suffixSum[i] = P(Y >= i)
	suffixSum[k+1] = 0.0 // Y >= k+1 的概率为0
	for y := k; y >= 0; y-- {
		suffixSum[y] = suffixSum[y+1] + parityLossProbs[y]
	}

	expectedLoss := 0.0 // 期望丢失的数据包数量

	// 计算数据包丢失数x=0到n的情况
	currentProb := dataLossProbs[0]
	for x := 0; x <= n; x++ {
		if x > 0 {
			// 使用递推关系更新概率
			currentProb = dataLossProbs[x]
		}

		if x <= k {
			// 需要满足 x + y > k 即 y > k - x
			threshold := k - x + 1
			var probRecoveryFail float64

			if threshold > k {
				probRecoveryFail = 0.0 // 不可能满足条件
			} else {
				probRecoveryFail = suffixSum[threshold] // P(Y >= threshold)
			}

			// 当恢复失败时，丢失x个数据包
			expectedLoss += float64(x) * currentProb * probRecoveryFail
		} else {
			// x > k 时恢复一定失败
			expectedLoss += float64(x) * currentProb
		}
	}

	// 计算丢包率并确保在[0,1]范围内
	lossRate := expectedLoss / float64(n)
	return math.Max(0.0, math.Min(lossRate, 1.0))
}

// 计算二项分布概率数组
func binomialProbabilities(n int, p float64) []float64 {
	probs := make([]float64, n+1)
	if n == 0 {
		probs[0] = 1.0
		return probs
	}

	// 初始化x=0的概率
	probs[0] = math.Pow(1-p, float64(n))
	if p == 0.0 {
		return probs
	}

	// 递推计算其他概率
	inv1MinusP := 1.0 / (1 - p)
	for x := 1; x <= n; x++ {
		// P(X=x) = P(X=x-1) * (n-x+1)/x * p/(1-p)
		probs[x] = probs[x-1] * float64(n-x+1) / float64(x) * p * inv1MinusP
	}
	return probs
}
