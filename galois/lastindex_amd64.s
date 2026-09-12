// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT	·_LastIndexByte(SB), NOSPLIT, $0-40
	MOVQ b_base+0(FP), SI
	MOVQ b_len+8(FP), CX
	MOVBLZX c+24(FP), AX
	LEAQ ret+32(FP), R8
	JMP  lastindexbytebody<>(SB)

// TEXT	·_LastIndexByteString(SB), NOSPLIT, $0-32
// 	MOVQ s_base+0(FP), SI
// 	MOVQ s_len+8(FP), CX
// 	MOVBLZX c+16(FP), AX
// 	LEAQ ret+24(FP), R8
// 	JMP  lastindexbytebody<>(SB)

// input:
//   SI: data
//   CX: data len
//   AL: byte sought
//   R8: address to put result
TEXT	lastindexbytebody<>(SB), NOSPLIT, $0
	// Shuffle X0 around so that each byte contains
	// the character we're looking for.
	MOVD AX, X0
	PUNPCKLBW X0, X0
	PUNPCKLBW X0, X0
	PSHUFL $0, X0, X0

	CMPQ CX, $16
	JLT small

	CMPQ CX, $32
	JA avx2
sse:
	LEAQ -16(SI)(CX*1), DI
	MOVQ SI, AX	// AX = address of first 16 bytes
	JMP	sseloopentry

sseloop:
	// Move the next 16-byte chunk of the data into X1.
	MOVOU	(DI), X1
	// Compare bytes in X0 to X1.
	PCMPEQB	X0, X1
	// Take the top bit of each byte in X1 and put the result in DX.
	PMOVMSKB X1, DX
	// Find last set bit, if any.
	BSRL	DX, DX
	JNZ	ssesuccess
	// Advance to previous block.
	SUBQ	$16, DI
sseloopentry:
	CMPQ	DI, AX
	JA	sseloop

	// Search the first 16-byte chunk. This chunk may overlap with the
	// chunks we've already searched, but that's ok.
	MOVQ	AX, DI
	MOVOU	(AX), X1
	PCMPEQB	X0, X1
	PMOVMSKB X1, DX
	BSRL	DX, DX
	JNZ	ssesuccess

failure:
	MOVQ $-1, (R8)
	RET

// We've found a chunk containing the byte.
// The chunk was loaded from DI.
// The index of the matching byte in the chunk is DX.
// The start of the data is SI.
ssesuccess:
	SUBQ SI, DI	// Compute offset of chunk within data.
	ADDQ DX, DI	// Add offset of byte within chunk.
	MOVQ DI, (R8)
	RET

// handle for lengths < 16
small:
	TESTL	CX, CX
	JEQ	failure

	// Check if we'll load across a page boundary.
	LEAQ	16(SI), AX
	TESTW	$0xff0, AX // Test if there is any bit higher than the first four set in address SI+16.
	JZ	endofpage      // If not then reading 16 bytes starting at SI crosses a 4kbyte page boundary. 

	MOVOU	 (SI), X1 // Load data into the high end of X1.
	PCMPEQB	 X0, X1	  // Compare target byte with each byte in data.
	PMOVMSKB X1, DX	  // Move result bits to integer register.
	NEGL	CX
	ADDL	$32, CX
	SHLL	CX, DX    // Clear bits that represent matches past end of data.
	SHRL	CX, DX	  // Shift desired bits back down to bottom of register.
	BSRL	DX, DX	  // Find last set bit.
	JZ	failure	      // No set bit, failure.
	MOVQ	DX, (R8)
	RET

endofpage:
	MOVOU	-16(SI)(CX*1), X1 // Load data
	PCMPEQB	X0, X1	    // Compare target byte with each byte in data.
	PMOVMSKB X1, DX	    // Move result bits to integer register.
	SHLL	CX, DX
	SHRL	$16, DX	    // Shift desired bits down to bottom of register.
	BSRL	DX, DX	    // Find last set bit.
	JZ	failure	        // No set bit, failure.
	CMPL	DX, CX
	JAE	failure	        // Match is past end of data.
	MOVQ	DX, (R8)
	RET

avx2:
	CMPB ·hasAVX2(SB), $1
	JNE sse

	VPBROADCASTB  X0, Y1
	LEAQ -32(SI)(CX*1), DI
	MOVQ SI, AX	// AX = address of first 32 bytes
avx2loop:
	// Move the next 32-byte chunk of the data into Y2.
	VMOVDQU (DI), Y2
	// Compare bytes in Y1 to Y2 an store in result in Y3.
	VPCMPEQB Y1, Y2, Y3
	// Test if any bit in Y3 is set.
	VPTEST Y3, Y3
	JNZ avx2success
	// Advance to previous block.
	SUBQ $32, DI
avx2loopentry:
	CMPQ DI, AX
	JA avx2loop

	// Search the first 16-byte chunk. This chunk may overlap with the
	// chunks we've already searched, but that's ok.
	MOVQ AX, DI
	VMOVDQU (DI), Y2
	VPCMPEQB Y1, Y2, Y3
	VPTEST Y3, Y3
	JNZ avx2success
	VZEROUPPER
	MOVQ $-1, (R8)
	RET

avx2success:
	VPMOVMSKB Y3, DX // Move result bits to integer register.
	BSRL DX, DX      // Find last set bit, if any.
	SUBQ SI, DI      // Compute offset of chunk within data.
	ADDQ DI, DX      // Add offset of byte within chunk.
	MOVQ DX, (R8)
	VZEROUPPER
	RET
