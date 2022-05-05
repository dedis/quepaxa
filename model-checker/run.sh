#!/bin/sh
# Analyze the consensus model using the Spin model checker.

# Exhaustive verification.
# MEMLIMIT is the memory-usage limit in megabytes.
spin -search -O2 -safety -DMEMLIM=60000 qp.pml

# Bitstate verification - most aggressive state compression.
# -w defines the power of two of the hash table size in bits.
#spin -search -O2 -safety -bitstate -w28 qp.pml

