// Simple model of QuePaxa consensus.
// Uses explicit message-based communication with recorders.

#define N		3	// total number of recorder (state) nodes
#define F		1	// number of failures tolerated
#define T		(N-F)	// consensus threshold required

#define M		2	// number of proposers (clients)

#define STEPHI		11	// highest step number to simulate
#define RAND		2	// random part of fitness space is 1..RAND
#define HI		(RAND+1) // top priority for proposals by leader
#define VALS		2	// space of preferred values is 1..VALS

// A proposal is an integer divided into two bit-fields: fitness and value.
#define	VALBITS		4
#define FITBITS		4
#define VALSHIFT	(0)
#define FITSHIFT	(VALBITS)
#define PROP(f,v)	(((f)<<FITSHIFT) | ((v)<<VALSHIFT))
#define VAL(p)		(((p) >> VALSHIFT) & ((1 << VALBITS)-1))
#define FIT(p)		(((p) >> FITSHIFT) & ((1 << FITBITS)-1))

#define MAX(a, b)	((a) > (b) -> (a) : (b))

byte leader;			// which proposer is the well-known leader
byte decided;			// proposed value that we've decided on
byte propsdone;			// number of proposers that have finished

// Channels for recorder/proposer communication.
chan creq[1+N] = [0] of { byte, byte, byte }		// <j, s, v>
chan crep[1+M] = [0] of { byte, byte, byte, byte, byte}	// <i, s, s', f', m'>

#define DECIDE(j, s, p)	atomic {					\
	printf("%d step %d decided <%d,%d>", j, s, FIT(p), VAL(p));	\
	assert(decided == 0 || decided == VAL(p));			\
	decided = VAL(p);						\
}

// Each proposer is a process.
proctype Proposer(byte j) {			// We're proposer j in 1..M
	byte s;
	byte p, g;
	byte ri, rs, rsn, rfn, rmn;	// responses we get from recorders
	byte i, sent, recs;	// request send and reply receiving state
	bit done;		// for detecting early-decision opportunities

	// Choose the arbitrary initial "preferred value" of this proposer
	s = 4;
	select (p : 1 .. VALS);	// select a "random" value into temporary
	printf("%d proposing %d\n", j, p);
	p = PROP(HI, p);

	// Initialize per-step state for the first step of the first round.
	printf("%d step %d\n", j, s);
	sent = 0;		// bit mask of recorders we've sent to
	recs = 0;		// number of recorders we've heard from
	g = 0;			// gather best response proposer saw so far
	done = true;

	i = 0;			// first, send to a channel no one listens on
	do
	:: creq[i] ! j, s, p ->	// send a request for this step to recorder i
		printf("%d step %d sent <%d,%d> to %d\n",
			j, s, FIT(p), VAL(p), i);
		sent = sent | (1 << i);		// successfully sent
		i = 0;				// now we have no target again

	:: s <= STEPHI && recs < T ->		// choose a recorder to send to

		// randomize fitness in phase 0 if we're not the leader
		if
		:: (s & 3) == 0 && j != leader ->
			byte r;
			select(r : 1 .. RAND);
			p = PROP(r, VAL(p));
		:: else -> skip
		fi
		assert(FIT(p) > 0 && VAL(p) > 0);

		// choose a recorder that we haven't already sent a request to
		// revert to i=0 if we've already sent to selected recorder
		select (i : 1 .. N);
		i = ((sent & (1 << i)) == 0 -> i : 0);

	:: crep[j] ? ri, rs, rsn, rfn, rmn ->	// get response from a recorder
		printf("%d step %d recv %d %d <%d,%d>,<%d,%d> from %d\n",
			j, s, rs, rsn, FIT(rfn), VAL(rfn),
			FIT(rmn), VAL(rmn),  i);
		assert(rs <= s);	// should get replies only to requests
		if
		:: rs < s -> skip	// discard old unneeded replies

		:: rs == s && rsn > s -> // catch up to new recorder state
			s = rsn;	// adopt recorder's round start state
			p = rfn;

			// initialize per-step state for the new time-step
			printf("%d step %d\n", j, s);
			sent = 0;	// bit mask of recorders we've sent to
			recs = 0;	// number of recorders we've heard from
			g = 0;		// best response proposer saw so far
			done = true;

		:: rs == s && rsn == s && (s & 3) == 0 -> // propose phase
			g = MAX(g, rfn); // gather best of all first proposals
			done = done && (FIT(rfn) == HI);
			recs++;	 	// this recorder has now replied

		:: rs == s && rsn == s && (s & 3) == 1 -> // spread E phase
			recs++;	 	// this recorder has now replied

		:: rs == s && rsn == s && (s & 3) >= 2 -> // gather E spread C
			g = MAX(g, rmn); // gather best of E or C sets
			recs++;	 	// this recorder has now replied
		fi
		assert(recs <= N);	// shouldn't get any extra replies

		ri = 0;			// clear temporaries
		rs = 0;
		rsn = 0;
		rfn = 0;
		rmn = 0;

	:: s <= STEPHI && recs >= T ->	// got a quorum of replies

		// handle the proposer's completion of this round
		if
		:: (s & 3) == 0 ->		// propose phase
			assert(FIT(g) > 0 && VAL(g) > 0);
			p = g;		// pick best of some E set

			// Decide early if all proposals were HI fit
			if
			:: done ->
				DECIDE(j, s, p);
			:: else -> skip
			fi

		:: (s & 3) == 1 -> skip		// spread E phase: nothing to do

		:: (s & 3) == 2 ->		// gather E spread C phase
			// p is now the best of some universal (U) set;
			// g is the best of all the E sets we gathered.
			assert(FIT(g) > 0 && VAL(g) > 0);
			if
			:: p == g ->
				DECIDE(j, s, p);
			:: else -> skip
			fi

		:: (s & 3) == 3 ->		// gather C phase
			// g is the best of all common (C) sets we gathered;
			// this becomes our proposal for the next round.
			assert(FIT(g) > 0 && VAL(g) > 0);
			p = g;
		fi

		// proceed to next logical time-step
		s = s + 1;

		// initialize per-step state for the new time-step
		printf("%d step %d\n", j, s);
		sent = 0;	// bit mask of recorders we've sent to
		recs = 0;	// number of recorders we've heard from
		g = 0;		// best response proposer saw so far
		done = true;

	:: s > STEPHI ->		// we've simulated enough time-steps
		break;
	od

	// count terminated proposers so recorders can terminate too
	atomic {
		propsdone++;
	}
}

// Each recorder is a process implementing an interval summary register (ISR).
proctype Recorder(byte i) {			// We're proposer j in 1..M
	byte s, f, a, m;
	byte rj, rs, rv;

	do
	:: creq[i] ? rj, rs, rv ->		// got request from proposer rj
		if
		:: rs == s ->
			a = MAX(a, rv);		// accumulate max of all values

		:: rs > s ->			// forward to a later step
			m = (rs == s+1 -> a : 0);
			s = rs;
			f = rv;
			a = rv;

		:: else -> skip
		fi

		// send reply to the proposer --
		// but don't block forever if all proposers terminate.
		if
		:: crep[rj] ! i, rs, s, f, m 	// reply succeeded
		:: propsdone == M -> break	// done while trying to send
		fi

		rj = 0;				// clear temporaries
		rs = 0;
		rv = 0;

	:: propsdone == M ->			// all proposers terminated?
		break;				// terminate recorder thread
	od
}

// The initialization process just gets things launched.
init {
	assert(HI < 1 << FITBITS);
	assert(VALS < 1 << VALBITS);

	decided = 0;				// we haven't decided yet

	// first choose the "well-known" leader, or 0 for no leader
	//leader = 0;				// no leader
	leader = 1;				// fixed leader
	//select (leader : 0 .. M);		// any (or no) leader

	atomic {
		int i, j;

		for (i : 1 .. N) {		// Launch N recorders
			run Recorder(i)
		}
		for (j : 1 .. M) {		// Launch M proposers
			run Proposer(j)
		}
	}
}

