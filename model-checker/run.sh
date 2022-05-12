#!/bin/sh
# Analyze the consensus model using the Spin model checker.

# Exhaustive verification.
# MEMLIMIT is the memory-usage limit in megabytes.
spin -search -O2 -safety -DMEMLIM=60000 $1

# Set maximum search depth (-m), making it an error to exceed this depth (-b).
#spin -search -O2 -safety -DMEMLIM=60000 -m3870 -b $1

# Exhaustive verification with state vector compression.
#spin -search -O2 -safety -DMEMLIM=60000 -collapse $1
#spin -search -O2 -safety -DMEMLIM=60000 -hc $1

# Bitstate verification - most aggressive state compression.
# -w defines the power of two of the hash table size in bits.
#	examples: -w28: 32MB, -w33: 1GB, -w38: 32GB
#spin -search -O2 -safety -bitstate -w28 $1
#spin -search -O2 -safety -bitstate -w38 $1

