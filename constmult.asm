// multiplies R[0] and R[1]
// result in R[2]

	@i
	M=0 // i = 0

	@mask
	M=1 // mask = 1

	@R2
	M=0 // R2 = 0

(LOOP)
	@mask
	D=M
	@END
	D;JLT // if mask < 0 jmp END

	@R0
	D=D&M // D = mask already
	@CONTINUE
	D;JEQ // if sub == 0 goto CONTINUE

	@i
	D=M
	@tens
	M=D // tens = i

	@shifts
	M=0 // shifts = 0

(COUNTSHIFTS)
	@tens
	D=M
	@SHIFTS
	D;JLE // if tens <= 0 goto SHIFTS

	@shifts
	M=M+1 // shifts ++

	@tens
	M=M-1 // tens--

	@COUNTSHIFTS
	0;JMP // goto COUNTSHIFTS

(SHIFTS)
	@R1
	D=M
	@interim
	M=D // interim = R1

(SHIFTLOOP)
	@shifts
	D=M
	@SUM
	D;JLE // if shifts <= 0 goto SUM

	@interim
	D=M
	M=D+M // interim = interim + interim (*2 or << 1)

	@shifts
	M=M-1 // shifts--
	@SHIFTLOOP
	0;JMP

(SUM)
	@interim
	D=M
	@R2
	M=D+M // R2 = R2 + interim

(CONTINUE)
	@i
	M=M+1 // i++
	@mask
	D=M
	M=D+M // mask = mask + mask (*2)
	@LOOP
	0;JMP

(END)
	@END
	0;JMP
