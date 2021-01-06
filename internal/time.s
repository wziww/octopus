#include "textflag.h"

TEXT ·walltime(SB),NOSPLIT,$0-12
	JMP	runtime·walltime1(SB)
