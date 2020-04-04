//-----------------------------------------------------------------------------
/*

RISC-V Instructions

*/
//-----------------------------------------------------------------------------

package rv

import "github.com/deadsy/rvdbg/util"

//-----------------------------------------------------------------------------

const (
	opcodeLB      = 0x00000003 // lb
	opcodeLH      = 0x00001003 // lh
	opcodeLW      = 0x00002003 // lw
	opcodeLD      = 0x00003003 // ld
	opcodeSB      = 0x00000023 // sb
	opcodeSH      = 0x00001023 // sh
	opcodeSW      = 0x00002023 // sw
	opcodeSD      = 0x00003023 // sd
	opcodeJAL     = 0x0000006f // jal
	opcodeXORI    = 0x00004013 // xori
	opcodeSRLI    = 0x00005013 // srli
	opcodeADDI    = 0x00000013 // addi
	opcodeEBREAK  = 0x00100073 // ebreak
	opcodeCSRRW   = 0x00001073 // csrrw
	opcodeCSRRS   = 0x00002073 // csrrs
	opcodeFMV_X_W = 0xe0000053 // fmv.x.w
	opcodeFMV_W_X = 0xf0000053 // fmv.w.x
	opcodeFMV_D_X = 0xf2000053 // fmv.d.x
	opcodeFMV_X_D = 0xe2000053 // fmv.x.d
	opcodeFLD     = 0x00003007 // fld
	opcodeFSD     = 0x00003027 // fsd
	opcodeFLW     = 0x00002007 // flw
	opcodeFSW     = 0x00002027 // fsw
)

//-----------------------------------------------------------------------------

// InsLD returns "ld rd, ofs(rs1)"
func InsLD(rd, ofs, rs1 uint) uint32 {
	return uint32((util.Bits(ofs, 11, 0) << 20) | (rs1 << 15) | (rd << 7) | opcodeLD)
}

// InsLW returns "lw rd, ofs(rs1)"
func InsLW(rd, ofs, rs1 uint) uint32 {
	return uint32((util.Bits(ofs, 11, 0) << 20) | (rs1 << 15) | (rd << 7) | opcodeLW)
}

// InsLH returns "lh rd, ofs(rs1)"
func InsLH(rd, ofs, rs1 uint) uint32 {
	return uint32((util.Bits(ofs, 11, 0) << 20) | (rs1 << 15) | (rd << 7) | opcodeLH)
}

// InsLB returns "lb rd, ofs(rs1)"
func InsLB(rd, ofs, rs1 uint) uint32 {
	return uint32((util.Bits(ofs, 11, 0) << 20) | (rs1 << 15) | (rd << 7) | opcodeLB)
}

// InsSD returns "sd rs2, ofs(rs1)"
func InsSD(rs2, ofs, rs1 uint) uint32 {
	return uint32((util.Bits(ofs, 11, 5) << 25) | (rs2 << 20) | (rs1 << 15) | (util.Bits(ofs, 4, 0) << 7) | opcodeSD)
}

// InsSW returns "sw rs2, ofs(rs1)"
func InsSW(rs2, ofs, rs1 uint) uint32 {
	return uint32((util.Bits(ofs, 11, 5) << 25) | (rs2 << 20) | (rs1 << 15) | (util.Bits(ofs, 4, 0) << 7) | opcodeSW)
}

// InsSH returns "sh rs2, ofs(rs1)"
func InsSH(rs2, ofs, rs1 uint) uint32 {
	return uint32((util.Bits(ofs, 11, 5) << 25) | (rs2 << 20) | (rs1 << 15) | (util.Bits(ofs, 4, 0) << 7) | opcodeSH)
}

// InsSB returns "sb rs2, ofs(rs1)"
func InsSB(rs2, ofs, rs1 uint) uint32 {
	return uint32((util.Bits(ofs, 11, 5) << 25) | (rs2 << 20) | (rs1 << 15) | (util.Bits(ofs, 4, 0) << 7) | opcodeSB)
}

// InsADDI returns "addi rd, rs1, imm"
func InsADDI(rd, rs1, imm uint) uint32 {
	return uint32((util.Bits(imm, 11, 0) << 20) | (rs1 << 15) | (rd << 7) | opcodeADDI)
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

// InsJAL returns "jal rd, ofs"
func InsJAL(rd, ofs uint) uint32 {
	offset := (util.Bit(ofs, 20) << 19) |
		(util.Bits(ofs, 10, 1) << 9) |
		(util.Bit(ofs, 11) << 8) |
		(util.Bits(ofs, 19, 12) << 0)
	return uint32((offset << 12) | (rd << 7) | opcodeJAL)
}

// InsXORI returns "xori rd, rs1, imm"
func InsXORI(rd, rs1, imm uint) uint32 {
	return uint32((util.Bits(imm, 11, 0) << 20) | (rs1 << 15) | (rd << 7) | opcodeXORI)
}

// InsSRLI returns "srli rd, rs1, shamt"
func InsSRLI(rd, rs1, shamt uint) uint32 {
	return uint32((shamt << 20) | (rs1 << 15) | (rd << 7) | opcodeSRLI)
}

// InsFSD returns "fsd rs2, ofs(rs1)"
func InsFSD(rs2, ofs, rs1 uint) uint32 {
	return uint32((util.Bits(ofs, 11, 5) << 25) | (rs2 << 20) | (rs1 << 15) | (util.Bits(ofs, 4, 0) << 7) | opcodeFSD)
}

// InsFSW returns "fsw rs2, ofs(rs1)"
func InsFSW(rs2, ofs, rs1 uint) uint32 {
	return uint32((util.Bits(ofs, 11, 5) << 25) | (rs2 << 20) | (rs1 << 15) | (util.Bits(ofs, 4, 0) << 7) | opcodeFSW)
}

// InsFLD returns "fld rd, ofs(rs1)"
func InsFLD(rd, ofs, rs1 uint) uint32 {
	return uint32((util.Bits(ofs, 11, 0) << 20) | (rs1 << 15) | (rd << 7) | opcodeFLD)
}

// InsFLW returns "flw rd, ofs(rs1)"
func InsFLW(rd, ofs, rs1 uint) uint32 {
	return uint32((util.Bits(ofs, 11, 0) << 20) | (rs1 << 15) | (rd << 7) | opcodeFLW)
}

//-----------------------------------------------------------------------------
