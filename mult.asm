// multiplies RAM[0] * RAM[1]
// puts the result in RAM[2]
//
// naive implementation with adding in a loop

	@R0
	D=M
	@i
	M=D // i = R0

	@R2
	M=0 // sum = 0

(LOOP)
	@i
	D=M
	@END
	D;JEQ // if i == 0 end

	@R2
	D=M
	@R1
	D=D+M
	@R2
	M=D // sum = sum + R1
	@i
	D=M
	M=D-1 // i = i -1
	@LOOP
	0;JMP

(END)
	@END
	0;JMP
