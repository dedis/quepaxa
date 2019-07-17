#!/bin/sh
# Exhaustively analyze the QSC model using the Spin model checker.
spin -a qsc.pml || exit 1
gcc -O2 -DSAFETY -DBITSTATE -o pan pan.c || exit 1
./pan -m20000
