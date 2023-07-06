package hardware

// All Instructions
var ALL_INS = [16]Instruction{
	{
		Opcode: "BRK",
		Binary: 0b0000,
		Nargs:  0,
	},
	{
		Opcode: "LDA imm.",
		Binary: 0b0001,
		Nargs:  1,
	},
	{
		Opcode: "LDA mem.",
		Binary: 0b0010,
		Nargs:  2,
	},
	{
		Opcode: "STA",
		Binary: 0b0011,
		Nargs:  2,
	},
	{
		Opcode: "IDC",
		Binary: 0b0100,
		Nargs:  2,
	},
	{
		Opcode: "ADD",
		Binary: 0b0101,
		Nargs:  2,
	},
	{
		Opcode: "---",
		Binary: 0b0110,
		Nargs:  0,
	},
	{
		Opcode: "---",
		Binary: 0b0111,
		Nargs:  0,
	},
	{
		Opcode: "NOT",
		Binary: 0b1000,
		Nargs:  0,
	},
	{
		Opcode: "ORA",
		Binary: 0b1001,
		Nargs:  2,
	},
	{
		Opcode: "AND",
		Binary: 0b1010,
		Nargs:  2,
	},
	{
		Opcode: "SHF",
		Binary: 0b1011,
		Nargs:  1,
	},
	{
		Opcode: "SLP",
		Binary: 0b1100,
		Nargs:  2,
	},
	{
		Opcode: "BNE",
		Binary: 0b1101,
		Nargs:  2,
	},
	{
		Opcode: "JMP imm.",
		Binary: 0b1110,
		Nargs:  2,
	},
	{
		Opcode: "JMP mem.",
		Binary: 0b1111,
		Nargs:  2,
	},
}
