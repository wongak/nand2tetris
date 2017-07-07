// on keypress screen will turn black
// on release screen will turn white
//
// begin
// d = key
// jmp if key != 0 pressed
// if turn == -1
// turn = 0
// goto draw
// pressed:
// turn = -1
// goto draw
// draw:
// d = turn
// for i = scn; i < kbd; i ++
// a = i
// m = d
// d = key
// goto begin

(BEGIN)
	@KBD
	D=M
	@PRESSED
	D;JNE // if KBD != 0 goto PRESSED
	
	@turn
	D=M
	@BEGIN
	D;JEQ // if turn == 0 goto BEGIN

	@turn
	M=0 // turn = 0
	@DRAW
	0;JMP // goto DRAW

(PRESSED)
	@turn
	M=-1

(DRAW)
	@SCREEN
	D=A
	@i
	M=D // i = SCREEN

(LOOP)
	@i
	D=M
	@KBD
	D=D-A
	@BEGIN
	D;JGE // if i - kbd >= 0 goto BEGIN

	@turn
	D=M
	@i
	A=M
	M=D // M[i] = turn

	@i
	M=M+1
	@LOOP
	0;JMP
