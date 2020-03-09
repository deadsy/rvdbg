//-----------------------------------------------------------------------------
/*

RISC-V Instructions

*/
//-----------------------------------------------------------------------------

package rv

import "github.com/deadsy/rvdbg/util"

//-----------------------------------------------------------------------------

const (
	// 	opcodeLUI         = 0x00000037 // lui
	// 	opcodeAUIPC       = 0x00000017 // auipc
	// 	opcodeJAL         = 0x0000006f // jal
	// 	opcodeJALR        = 0x00000067 // jalr
	// 	opcodeBEQ         = 0x00000063 // beq
	// 	opcodeBNE         = 0x00001063 // bne
	// 	opcodeBLT         = 0x00004063 // blt
	// 	opcodeBGE         = 0x00005063 // bge
	// 	opcodeBLTU        = 0x00006063 // bltu
	// 	opcodeBGEU        = 0x00007063 // bgeu
	// 	opcodeLB          = 0x00000003 // lb
	// 	opcodeLH          = 0x00001003 // lh
	opcodeLW = 0x00002003 // lw
	// 	opcodeLBU         = 0x00004003 // lbu
	// 	opcodeLHU         = 0x00005003 // lhu
	// 	opcodeSB          = 0x00000023 // sb
	// 	opcodeSH          = 0x00001023 // sh
	opcodeSW = 0x00002023 // sw
	// 	opcodeADDI        = 0x00000013 // addi
	// 	opcodeSLTI        = 0x00002013 // slti
	// 	opcodeSLTIU       = 0x00003013 // sltiu
	// 	opcodeXORI        = 0x00004013 // xori
	// 	opcodeORI         = 0x00006013 // ori
	// 	opcodeANDI        = 0x00007013 // andi
	// 	opcodeSLLI        = 0x00001013 // slli
	// 	opcodeSRLI        = 0x00005013 // srli
	// 	opcodeSRAI        = 0x40005013 // srai
	// 	opcodeADD         = 0x00000033 // add
	// 	opcodeSUB         = 0x40000033 // sub
	// 	opcodeSLL         = 0x00001033 // sll
	// 	opcodeSLT         = 0x00002033 // slt
	// 	opcodeSLTU        = 0x00003033 // sltu
	// 	opcodeXOR         = 0x00004033 // xor
	// 	opcodeSRL         = 0x00005033 // srl
	// 	opcodeSRA         = 0x40005033 // sra
	// 	opcodeOR          = 0x00006033 // or
	// 	opcodeAND         = 0x00007033 // and
	// 	opcodeFENCE       = 0x0000000f // fence
	// 	opcodeFENCE_I     = 0x0000100f // fence.i
	// 	opcodeECALL       = 0x00000073 // ecall
	opcodeEBREAK = 0x00100073 // ebreak
	// 	opcodeURET        = 0x00200073 // uret
	// 	opcodeSRET        = 0x10200073 // sret
	// 	opcodeMRET        = 0x30200073 // mret
	// 	opcodeWFI         = 0x10500073 // wfi
	// 	opcodeSFENCE_VMA  = 0x12000073 // sfence.vma
	// 	opcodeHFENCE_BVMA = 0x22000073 // hfence.bvma
	// 	opcodeHFENCE_GVMA = 0xa2000073 // hfence.gvma
	opcodeCSRRW = 0x00001073 // csrrw
	opcodeCSRRS = 0x00002073 // csrrs
// 	opcodeCSRRC       = 0x00003073 // csrrc
// 	opcodeCSRRWI      = 0x00005073 // csrrwi
// 	opcodeCSRRSI      = 0x00006073 // csrrsi
// 	opcodeCSRRCI      = 0x00007073 // csrrci
// 	opcodeMUL         = 0x02000033 // mul
// 	opcodeMULH        = 0x02001033 // mulh
// 	opcodeMULHSU      = 0x02002033 // mulhsu
// 	opcodeMULHU       = 0x02003033 // mulhu
// 	opcodeDIV         = 0x02004033 // div
// 	opcodeDIVU        = 0x02005033 // divu
// 	opcodeREM         = 0x02006033 // rem
// 	opcodeREMU        = 0x02007033 // remu
// 	opcodeLR_W        = 0x1000202f // lr.w
// 	opcodeSC_W        = 0x1800202f // sc.w
// 	opcodeAMOSWAP_W   = 0x0800202f // amoswap.w
// 	opcodeAMOADD_W    = 0x0000202f // amoadd.w
// 	opcodeAMOXOR_W    = 0x2000202f // amoxor.w
// 	opcodeAMOAND_W    = 0x6000202f // amoand.w
// 	opcodeAMOOR_W     = 0x4000202f // amoor.w
// 	opcodeAMOMIN_W    = 0x8000202f // amomin.w
// 	opcodeAMOMAX_W    = 0xa000202f // amomax.w
// 	opcodeAMOMINU_W   = 0xc000202f // amominu.w
// 	opcodeAMOMAXU_W   = 0xe000202f // amomaxu.w
// 	opcodeFLW         = 0x00002007 // flw
// 	opcodeFSW         = 0x00002027 // fsw
// 	opcodeFMADD_S     = 0x00000043 // fmadd.s
// 	opcodeFMSUB_S     = 0x00000047 // fmsub.s
// 	opcodeFNMSUB_S    = 0x0000004b // fnmsub.s
// 	opcodeFNMADD_S    = 0x0000004f // fnmadd.s
// 	opcodeFADD_S      = 0x00000053 // fadd.s
// 	opcodeFSUB_S      = 0x08000053 // fsub.s
// 	opcodeFMUL_S      = 0x10000053 // fmul.s
// 	opcodeFDIV_S      = 0x18000053 // fdiv.s
// 	opcodeFSQRT_S     = 0x58000053 // fsqrt.s
// 	opcodeFSGNJ_S     = 0x20000053 // fsgnj.s
// 	opcodeFSGNJN_S    = 0x20001053 // fsgnjn.s
// 	opcodeFSGNJX_S    = 0x20002053 // fsgnjx.s
// 	opcodeFMIN_S      = 0x28000053 // fmin.s
// 	opcodeFMAX_S      = 0x28001053 // fmax.s
// 	opcodeFCVT_W_S    = 0xc0000053 // fcvt.w.s
// 	opcodeFCVT_WU_S   = 0xc0100053 // fcvt.wu.s
// 	opcodeFMV_X_W     = 0xe0000053 // fmv.x.w
// 	opcodeFEQ_S       = 0xa0002053 // feq.s
// 	opcodeFLT_S       = 0xa0001053 // flt.s
// 	opcodeFLE_S       = 0xa0000053 // fle.s
// 	opcodeFCLASS_S    = 0xe0001053 // fclass.s
// 	opcodeFCVT_S_W    = 0xd0000053 // fcvt.s.w
// 	opcodeFCVT_S_WU   = 0xd0100053 // fcvt.s.wu
// 	opcodeFMV_W_X     = 0xf0000053 // fmv.w.x
// 	opcodeFLD         = 0x00003007 // fld
// 	opcodeFSD         = 0x00003027 // fsd
// 	opcodeFMADD_D     = 0x02000043 // fmadd.d
// 	opcodeFMSUB_D     = 0x02000047 // fmsub.d
// 	opcodeFNMSUB_D    = 0x0200004b // fnmsub.d
// 	opcodeFNMADD_D    = 0x0200004f // fnmadd.d
// 	opcodeFADD_D      = 0x02000053 // fadd.d
// 	opcodeFSUB_D      = 0x0a000053 // fsub.d
// 	opcodeFMUL_D      = 0x12000053 // fmul.d
// 	opcodeFDIV_D      = 0x1a000053 // fdiv.d
// 	opcodeFSQRT_D     = 0x5a000053 // fsqrt.d
// 	opcodeFSGNJ_D     = 0x22000053 // fsgnj.d
// 	opcodeFSGNJN_D    = 0x22001053 // fsgnjn.d
// 	opcodeFSGNJX_D    = 0x22002053 // fsgnjx.d
// 	opcodeFMIN_D      = 0x2a000053 // fmin.d
// 	opcodeFMAX_D      = 0x2a001053 // fmax.d
// 	opcodeFCVT_S_D    = 0x40100053 // fcvt.s.d
// 	opcodeFCVT_D_S    = 0x42000053 // fcvt.d.s
// 	opcodeFEQ_D       = 0xa2002053 // feq.d
// 	opcodeFLT_D       = 0xa2001053 // flt.d
// 	opcodeFLE_D       = 0xa2000053 // fle.d
// 	opcodeFCLASS_D    = 0xe2001053 // fclass.d
// 	opcodeFCVT_W_D    = 0xc2000053 // fcvt.w.d
// 	opcodeFCVT_WU_D   = 0xc2100053 // fcvt.wu.d
// 	opcodeFCVT_D_W    = 0xd2000053 // fcvt.d.w
// 	opcodeFCVT_D_WU   = 0xd2100053 // fcvt.d.wu
)

//-----------------------------------------------------------------------------

// InsLW returns "lw rd, ofs(base)"
func InsLW(rd, base, ofs uint) uint32 {
	return uint32((util.Bits(ofs, 11, 0) << 20) | (base << 15) | (rd << 7) | opcodeLW)
}

// InsSW returns "sw rs, ofs(base)"
func InsSW(rs, base, ofs uint) uint32 {
	return uint32((util.Bits(ofs, 11, 5) << 25) | (rs << 20) | (base << 15) | (util.Bits(ofs, 4, 0) << 7) | opcodeSW)
}

// InsEBREAK returns "ebreak"
func InsEBREAK() uint32 {
	return uint32(opcodeEBREAK)
}

// InsCSRR returns "csrr rd, csr"
func InsCSRR(rd, csr uint) uint32 {
	// csrrs rd, csr, x0
	return uint32((csr << 20) | (RegZero << 15) | (rd << 7) | opcodeCSRRS)
}

// InsCSRW returns "csrw csr, rs1"
func InsCSRW(csr, rs1 uint) uint32 {
	// csrrw x0, csr, rs1
	return uint32((csr << 20) | (rs1 << 15) | (RegZero << 7) | opcodeCSRRW)
}

//-----------------------------------------------------------------------------
