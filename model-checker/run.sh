#!/bin/sh
# Exhaustively analyze the QSC model using the Spin model checker.
spin -a qp.pml || exit 1
gcc -O2 -DSAFETY -DBITSTATE -o pan pan.c || exit 1
./pan -m50000
