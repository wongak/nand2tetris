package main

const (
	stackPointerAddress = 0 // SP
	stackBaseAddress    = 256

	lclBasePointerAddress  = 1 // LCL
	argBasePointerAddress  = 2 // ARG
	thisBasePointerAddress = 3 // THIS
	thatBasePointerAddress = 4 // THAT
)

// pop local 2
// example demonstrating hack assembly for working with stacks
//
// implementation: addr=LCL+2, SP--, *addr=*SP
//
// hack assembly:
//
//     @2
//     D=A
//     @1 // LCL
//     D=D+M // LCL+2
//     @addr
//     M=D // addr = LCL+2
//
//     @0 // SP
//     M=M-1 // SP--
//     A=M
//     D=M // D=*SP
//     @addr
//     M=D // *addr=*SP
//
// can be generalized with 2 = i
//
// push local i
//
// implementation: addr=LCL+i, *SP=*addr, SP++
//
// hack assembly:
//
//     @i // i
//     D=A
//     @1 // LCL
//     D=D+M // LCL+i
//     A=D
//     D=M // *(LCL+i)
//
//     @0 // SP
//     A=M
//     M=D // *SP=*addr
//     @0
//     M=M+1
//
// push constant i
//
// implementation: *SP=i, SP++
//
// hack assembly:
//
//     @i // i
//     D=A
//     @0 // SP
//     A=M // *SP
//     M=D // =i
//     @0
//     M=M+1
