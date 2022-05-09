// Simple model of Que Paxa consensus.
// Recorder logic runs atomically in-line within the proposer code.

#define N		3	// total number of recorder (state) nodes
#define F		1	// number of failures tolerated
#define T		(N-F)	// consensus threshold required

#define M		2	// number of proposers (clients)

#define ROUNDS		2	// number of consensus rounds to run
#define RAND		2	// random part of fitness space is 1..RAND
#define HI		(RAND+1) // top priority for proposals by leader
#define VALS		2	// space of preferred values is 1..VALS

// A proposal is an integer divided into three bit-fields: FIT, CLI, VAL.
#define	VALBITS		4
#define FITBITS		4
#define VALSHIFT	(0)
#define FITSHIFT	(VALBITS)
#define PROP(f,v)	(((f)<<FITSHIFT) | ((v)<<VALSHIFT))
#define VAL(p)		(((p) >> VALSHIFT) & ((1 << VALBITS)-1))
#define FIT(p)		(((p) >> FITSHIFT) & ((1 << FITBITS)-1))

#define MAX(a, b)	((a) > (b) -> (a) : (b))

// Recorder state: implements an epoch summary primitive (ESP),
// which returns the first value submitted in this epoch
// and the maximum of all values submitted in the prior epoch.
typedef Rec {
	byte s;			// step/epoch number
	byte f;			// first value submitted in this epoch
	byte a;			// maximum value seen so far in this epoch
	byte m;			// maximum value seen in prior epoch (s-1)
}

Rec rec[1+N];			// state of recorder nodes 1..N
byte decided;			// proposed value that we've decided on
byte leader;			// which proposer is the well-known leader

#define DECIDE(j, s, p)	atomic {					\
	printf("%d step %d decided <%d,%d>", j, s, FIT(p), VAL(p));	\
	assert(decided == 0 || decided == VAL(p));			\
	decided = VAL(p);						\
}


// We model one process per proposer.
proctype Proposer(byte j) {			// We're proposer j in 1..M
	byte s, t;
	byte p, g;
	byte i, recs, mask;	// recorders we've interacted with
	bit done;		// for detecting early-decision opportunities

	// Choose the arbitrary initial "preferred value" of this proposer
	s = 4;
	select (t : 1 .. VALS);	// select a "random" value into temporary
	p = PROP(HI, t);
	printf("%d proposing %d\n", j, t);

	do			// iterate over time-steps
	:: s <= ROUNDS*4 ->
		printf("%d step %d\n", j, s);

		// Send <s,p> and get reply from threshold of recorders
		recs = 0;	// number of recorders we've heard from
		mask = 0;	// bit mask of those recorders
		g = 0;		// gather best response proposer saw so far
		done = true;
		select (i : 1 .. N);	// first recorder to interact with
		do		// interact with the recorders in any order
		:: recs < T && (mask & (1 << i)) == 0 ->

			atomic {
				// Randomize fitnesses if we're not the leader
				if
				:: (s & 3) == 0 && j != leader ->
					select(t : 1 .. RAND);
					p = PROP(t, VAL(p));
				:: else -> skip
				fi
				assert(FIT(p) > 0 && VAL(p) > 0);

				// enter the recorder/ESP role (via "RPC").
				printf("%d step %d ESP <%d,%d> to %d\n",
					j, s, FIT(p), VAL(p), i);

				// first catch up the recorder if appropriate
				if
				:: s > rec[i].s ->
					rec[i].m = ((s == (rec[i].s+1)) ->
							rec[i].a : 0);
					rec[i].s = s;
					rec[i].f = p;
					rec[i].a = p;

				:: s == rec[i].s ->
					rec[i].a = MAX(rec[i].a, p);

				:: else -> skip
				fi

				// we're back to the proposer's logic now,
				// incorporating the recorder's "response".
				assert(s <= rec[i].s);
				if
				:: s == rec[i].s && (s & 3) == 0 ->
					g = MAX(g, rec[i].f);	// gather props
					done = done && (FIT(rec[i].f) == HI);

				:: s == rec[i].s && (s & 3) == 1 -> skip

				:: s == rec[i].s && (s & 3) >= 2 ->
					printf("%d step %d got <%d,%d> from %d\n", j, s, FIT(rec[i].m), VAL(rec[i].m), i);
					g = MAX(g, rec[i].m);	// gather E/C

				:: s < rec[i].s ->	// catch up proposer
					s = rec[i].s;
					p = rec[i].f;
					break;
				fi
				assert(s == rec[i].s);

				recs++;	 // this recorder has now "replied"
				mask = mask | (1 << i);

				select (i : 1 .. N);	// choose next recorder

			} // atomic

		:: recs < T && (mask & (1 << i)) != 0 ->
			// we've already gotten a reply from this recorder,
			// so just pick a different one.
			select (i : 1 .. N);

		:: recs == T ->	// we've heard from a threshold of recorders

			if
			:: (s & 3) == 0 ->	// propose phase
				assert(FIT(g) > 0 && VAL(g) > 0);
				p = g;		// pick best of some E set

				// Decide early if all proposals were HI fit
				if
				:: done ->
					DECIDE(j, s, p);
				:: else -> skip
				fi

			:: (s & 3) == 1 -> skip	// spreadE phase

			:: (s & 3) == 2 ->	// gatherEspreadC phase
				// p is now the best of a U set;
				// g is the best of all gathered E sets
				assert(FIT(g) > 0 && VAL(g) > 0);
				if
				:: p == g ->
					DECIDE(j, s, p);
				:: else -> skip
				fi

			:: (s & 3) == 3 ->	// gatherC phase
				// g is the best of all gathered C sets.
				// this is our proposal for the next round.
				assert(FIT(g) > 0 && VAL(g) > 0);
				p = g;
			fi
			s = s + 1;
			break;
		od

	:: s > ROUNDS*4 ->	// we've simulated enough time-steps
			break;
	od
}

init {
	assert(HI < 1 << FITBITS);
	assert(VALS < 1 << VALBITS);

	decided = 0;				// we haven't decided yet

	// first choose the "well-known" leader, or 0 for no leader
	//leader = 0;				// no leader
	leader = 1;				// fixed leader
	//select (leader : 0 .. M);		// any (or no) leader

	atomic {
		int i;
		for (i : 1 .. M) {		// Launch M proposers
			run Proposer(i)
		}
	}
}

