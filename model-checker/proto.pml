
#define N 3			// number of nodes
#define F 1			// number of faulty nodes
#define T (N-F)			// consensus threshold

#define TIME	5		// number of time steps to run
#define CMAX	((2+N)*N*TIME)	// max size of channels

mtype = {MSG, ACK, SUC};

typedef Recv {
	int prop[N];		// value proposed by each node, 0 if none
	bit win[N];		// whether value proposed is winning
	byte acks[N];		// # acks we've received for each node's prop
	byte sucs[N];		// # success confs we've gotten for proposal
}

typedef Node {
	chan in = [CMAX] of { byte, mtype, byte, byte, byte, bit } // node's inbox
}

Node node[N];			// each node's channel and receive state
byte sucp[TIME];			// successful proposal at each time, if any

active[N] proctype NodeProc() {
	int ltime = 0;		// local clock
	Recv recv[TIME];	// everything we've received for each timestep
	int from, stime, sndr, sprop;
	bit swin;
	do
	::	// pick a nondeterministic local proposal
		int lprop;
		if
		:: skip -> lprop = 111;
		:: skip -> lprop = 222;
		fi

		// "flip a coin" to see if it's a winning proposal
		// (should be prob 1/N but only matters for performance)
		bit lwin;
		if
		:: skip -> lwin = 0;
		:: skip -> lwin = 1;
		fi

		// broadcast our proposal
		int i;
		for (i : 0 .. N-1) {
			node[i].in ! _pid, MSG, ltime, _pid, lprop, lwin;
		}

		// process messages until we receive enough acks for this time
		do
		:: node[_pid].in ? from, MSG, stime, sndr, sprop, swin ->
			assert(sprop > 0);
			assert(recv[stime].prop[sndr] == 0 ||
				recv[stime].prop[sndr] == sprop);
			recv[stime].prop[sndr] = sprop;
			recv[stime].win[sndr] = swin;

			// Acknowledge it.
			for (i : 0 .. N-1) {
				node[i].in ! _pid, ACK, stime, sndr, sprop, swin;
			}

		:: node[_pid].in ? from, ACK, stime, sndr, sprop, swin ->
			recv[stime].acks[sndr]++;
			assert(sprop > 0);

			// We might receive an ACK for sndr's message
			// before the MSG itself; that's OK.
			assert(recv[stime].prop[sndr] == 0 ||
				recv[stime].prop[sndr] == sprop);
			recv[stime].prop[sndr] = sprop;
			recv[stime].win[sndr] = swin;

			// Increment the local clock as soon as we get
			// a threshold T of acks on a message from anyone
			assert(recv[stime].acks[sndr] <= N);
			if
			:: (stime >= ltime &&
			    recv[stime].acks[sndr] == T) ->
				break;
			:: else
			fi

		:: node[_pid].in ? from, SUC, stime, sndr, sprop, swin ->
			recv[stime].sucs[sndr]++;
			assert(recv[stime].sucs[sndr] <= N);

			// Do we have a threshold of success judgments?
			if
			:: (recv[stime].sucs[sndr] >= T) ->
				printf("winner at time %d\n", stime);
				assert(sucp[stime] == 0 ||
					sucp[stime] == recv[stime].prop[sndr]);
				sucp[stime] = recv[stime].prop[sndr];
			:: else
			fi
		od

		// Increment our local time, catching up with others as needed
		do
		:: (ltime <= stime) ->		// catch up

			// How many winning proposals did we see at ltime?
			int nwin = 0, winr;
			for (i : 0 .. N-1) {
				if
				:: (recv[ltime].win[i]) ->
					nwin++;
					winr = i;
				:: else
				fi
			}
			if
			:: (nwin == 1) ->	// might be a success round
				for (i : 0 .. N-1) {
					node[i].in ! _pid, SUC, ltime, winr,
						recv[ltime].prop[winr], 1;
				}
			:: else
			fi

			// Increment the current time
			ltime++;

		:: (ltime == stime+1) ->	// we're caught up
			break;
i
		:: else ->
			assert(0);
		od

		if
		:: (ltime >= TIME) -> break
		:: else
		fi;
	od;
}

