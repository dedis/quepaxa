#!/bin/sh
# Exhaustively analyze the QSC model using the Spin model checker.
spin -a qsc.pml
gcc -O2 -DSAFETY -DBITSTATE -o pan pan.c 
./pan -m20000
