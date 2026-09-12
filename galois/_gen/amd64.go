package main

import (
	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/ir"
	. "github.com/mmcloughlin/avo/operand"
)

func amd64() {
	amd64_mulVect_16()
	amd64_mulXorVect_16()
	amd64_xorVect_16()
	// amd64_tailZero_16()
}

func amd64_mulVect_16() {
	TEXT("_mulVect_16", NOSPLIT, "func(lohi *[2][16]byte, in, out []byte) ")
	Doc("要求数据长度不小于16")
	lohi := Load(Param("lohi"), GP64())
	in := Load(Param("in").Base(), GP64())
	len := Load(Param("in").Len(), GP64())
	out := Load(Param("out").Base(), GP64())
	end := GP64() // 数据结尾地址
	MOVQ(in, end)
	ADDQ(len, end)

	var (
		lo, hi = XMM(), XMM() // 被查询表
		mask   = XMM()        // 每字节低4位掩码
		l1, h1 = XMM(), XMM() //
	)
	VMOVDQU(Mem{Base: lohi, Disp: 0}, lo)
	VMOVDQU(Mem{Base: lohi, Disp: 16}, hi)
	VMOVDQU(NewDataAddr(Symbol{Name: "·lomask"}, 0), mask)

	Instruction(&ir.Instruction{
		Opcode:   "PCALIGN",
		Operands: []Op{Imm(16)},
	})
	Label("mulVect_start")
	{
		ADDQ(Imm(16), in)
		CMPQ(in, end)
		JA(LabelRef("mulVect_end"))

		VMOVDQU(Mem{Base: in, Disp: -16}, l1)
		VPSRLQ(Imm(4), l1, h1)
		VPAND(mask, l1, l1)
		VPAND(mask, h1, h1)
		VPSHUFB(l1, lo, l1)
		VPSHUFB(h1, hi, h1)
		VPXOR(l1, h1, l1)
		VMOVDQU(l1, Mem{Base: out})

		ADDQ(Imm(16), out)
		JMP(LabelRef("mulVect_start"))
	}
	Label("mulVect_end")
	{ // 尾部处理
		in := Load(Param("in").Base(), GP64())
		out := Load(Param("out").Base(), GP64())
		ADDQ(len, in)
		ADDQ(len, out)

		VMOVDQU(Mem{Base: in, Disp: -16}, l1)
		VPSRLQ(Imm(4), l1, h1)
		VPAND(mask, l1, l1)
		VPAND(mask, h1, h1)
		VPSHUFB(l1, lo, l1)
		VPSHUFB(h1, hi, h1)
		VPXOR(l1, h1, l1)
		VMOVDQU(l1, Mem{Base: out, Disp: -16})
	}
	RET()
}
func amd64_mulXorVect_16() {
	TEXT("_mulXorVect_16", NOSPLIT, "func(lohi *[2][16]byte, in, out []byte) ")
	Doc("要求数据长度不小于16")
	lohi := Load(Param("lohi"), GP64())
	in := Load(Param("in").Base(), GP64())
	len := Load(Param("in").Len(), GP64())
	out := Load(Param("out").Base(), GP64())
	end := GP64() // 数据结尾地址
	MOVQ(in, end)
	ADDQ(len, end)

	var (
		lo, hi = XMM(), XMM() // 被查询表
		mask   = XMM()        // 每字节低4位掩码
		l1, h1 = XMM(), XMM() //
		tail   = XMM()        // 暂存out末尾16B数据包
	)
	VMOVDQU(Mem{Base: lohi, Disp: 0}, lo)
	VMOVDQU(Mem{Base: lohi, Disp: 16}, hi)
	VMOVDQU(NewDataAddr(Symbol{Name: "·lomask"}, 0), mask)
	VMOVDQU(Mem{Base: out, Index: len, Scale: 1, Disp: -16}, tail)

	Instruction(&ir.Instruction{
		Opcode:   "PCALIGN",
		Operands: []Op{Imm(16)},
	})
	Label("mulXorVect_start")
	{
		ADDQ(Imm(16), in)
		CMPQ(in, end)
		JA(LabelRef("mulXorVect_end"))

		VMOVDQU(Mem{Base: in, Disp: -16}, l1)
		VPSRLQ(Imm(4), l1, h1)
		VPAND(mask, l1, l1)
		VPAND(mask, h1, h1)
		VPSHUFB(l1, lo, l1)
		VPSHUFB(h1, hi, h1)
		VPXOR(l1, h1, l1)
		VPXOR(Mem{Base: out}, l1, l1)
		VMOVDQU(l1, Mem{Base: out})

		ADDQ(Imm(16), out)
		JMP(LabelRef("mulXorVect_start"))
	}
	Label("mulXorVect_end")
	{ // 尾部处理
		in := Load(Param("in").Base(), GP64())
		out := Load(Param("out").Base(), GP64())
		ADDQ(len, in)
		ADDQ(len, out)

		VMOVDQU(Mem{Base: in, Disp: -16}, l1)
		VPSRLQ(Imm(4), l1, h1)
		VPAND(mask, l1, l1)
		VPAND(mask, h1, h1)
		VPSHUFB(l1, lo, l1)
		VPSHUFB(h1, hi, h1)
		VPXOR(l1, h1, l1)
		VPXOR(tail, l1, l1)
		VMOVDQU(l1, Mem{Base: out, Disp: -16})
	}
	RET()
}

func amd64_xorVect_16() {
	TEXT("_xorVect_16", NOSPLIT, "func(in, out []byte) ")
	Doc("要求数据长度不小于16")
	in := Load(Param("in").Base(), GP64())
	len := Load(Param("in").Len(), GP64())
	out := Load(Param("out").Base(), GP64())
	end := GP64() // 数据结尾地址
	MOVQ(in, end)
	ADDQ(len, end)

	var (
		x    = XMM()
		tail = XMM()
	)
	VMOVDQU(Mem{Base: out, Index: len, Scale: 1, Disp: -16}, tail)

	Instruction(&ir.Instruction{
		Opcode:   "PCALIGN",
		Operands: []Op{Imm(16)},
	})
	Label("xorVect_start")
	{
		ADDQ(Imm(16), in)
		CMPQ(in, end)
		JA(LabelRef("xorVect_end"))

		VMOVDQU(Mem{Base: in, Disp: -16}, x)
		VPXOR(Mem{Base: out}, x, x)
		VMOVDQU(x, Mem{Base: out})

		ADDQ(Imm(16), out)
		JMP(LabelRef("xorVect_start"))
	}
	Label("xorVect_end")
	{ // 尾部处理
		in := Load(Param("in").Base(), GP64())
		out := Load(Param("out").Base(), GP64())
		ADDQ(len, in)
		ADDQ(len, out)

		VMOVDQU(Mem{Base: in, Disp: -16}, x)
		VPXOR(tail, x, x)
		VMOVDQU(x, Mem{Base: out, Disp: -16})
	}
	RET()
}

func amd64_tailZero_16() {
	TEXT("_tailZero_16", NOSPLIT, "func(b []byte) int")
	Doc("要求数据长度不小于16")
	b := Load(Param("b").Base(), GP64())
	len := Load(Param("b").Len(), GP64())
	end := GP64() // 数据结尾地址
	MOVQ(b, end)
	ADDQ(len, end)

	var (
		x    = XMM()
		zero = XMM()
		mask = GP32()
	)
	VPXOR(zero, zero, zero)

	Instruction(&ir.Instruction{
		Opcode:   "PCALIGN",
		Operands: []Op{Imm(16)},
	})
	Label("tailZero_start")
	{
		SUBQ(Imm(16), end)
		CMPQ(end, b)
		JB(LabelRef("tailZero_end"))

		VMOVDQU(Mem{Base: end}, x)      //
		PCMPEQB(zero, x)                // 按字节比较是否相等, 相等为0xff, 否则为0
		PMOVMSKB(x, mask)               // 将x每字节高位提取, 组成16为整数; 此时0值的mask为1
		CMPW(mask.As16(), Imm(0xffff))  //
		JNE(LabelRef("tailZero_done"))  //
		JMP(LabelRef("tailZero_start")) //
	}
	Label("tailZero_end")
	{ // 处理尾部数据
		MOVQ(b, end)

		VMOVDQU(Mem{Base: end}, x)
		PCMPEQB(zero, x)
		PMOVMSKB(x, mask)
	}

	Label("tailZero_done")
	{
		NOTW(mask.As16())                                          //
		BSRW(mask.As16(), mask.As16())                             // 高位向低位扫描, 第一个1的位置
		JZ(LabelRef("tailZero_notfound"))                          // zero flag 未被设置, 没找到
		m := Mem{Base: end, Index: mask.As64(), Scale: 1, Disp: 1} //
		LEAQ(m, end)                                               // end = end + mask + 1, +1是因为BSRW返回的是索引而非长度
		SUBQ(b, end)                                               //
		Store(end, ReturnIndex(0))                                 //
		RET()

		Label("tailZero_notfound")
		MOVQ(I64(-1), end)
		Store(end, ReturnIndex(0))
		RET()
	}
}
