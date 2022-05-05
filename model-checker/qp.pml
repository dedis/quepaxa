// Simple model of Que Paxa consensus

#define N		3	// total number of recorder (state) nodes
#define F		1	// number of failures tolerated
#define T		(N-F)	// consensus threshold required

#define M		2	// number of proposers (clients)

#define ROUNDS		2	// number of consensus rounds to run
#define RAND		2	// random part of fitness space is 1..RAND
#define HI		(RAND*N+N+1) // top priority for proposals by leader
#define VALS		2	// space of preferred values is 1..VALS

// A proposal is a <fitness, value> pair.
typedef Prop {
	byte fit;		// proposal priority: random*N+node#, or HI
	byte val;		// proposed value
}

#define ZERO(d)		d.fit = 0; d.val = 0;
#define COPY(d, s)	d.fit = s.fit; d.val = s.val;
#define BEST(d, s)	if						\
			:: s.fit > d.fit -> 				\
				d.fit = s.fit;				\
				d.val = s.val;				\
			:: s.fit == d.fit ->				\
				assert(s.val == d.val);			\
			:: else -> skip					\
			fi

// Recorder state
typedef Rec {
	byte s;			// step number
	Prop p;			// (best) proposal from proposer this round
	Prop e;			// best existent proposal seen this round
	Prop c;			// best common proposal seen this round
}

Rec rec[1+N];			// state of recorder nodes 1..N
byte decided;			// proposed value that we've decided on
byte leader;			// which proposer is the well-known leader

#define DECIDE(j, s, p)	atomic {					\
				printf("%d step %d decided <%d,%d>", j, s, p.fit, p.val); \
				assert(decided == 0 || decided == p.val); \
				decided = p.val;			\
			}


// We model one process per proposer.
proctype Proposer(byte j) {			// We're proposer j in 1..M
	byte s, t;
	Prop p, e, c, g;
	byte i, recs, mask;	// recorders we've interacted with
	bit done;		// for detecting early-decision opportunities

	// Choose the arbitrary initial "preferred value" of this proposer
	s = 4;
	select (t : 1 .. VALS);	// select a "random" value into temporary
	p.val = t;
	printf("%d proposing %d\n", j, t);
	ZERO(e);		// best of no proposals so far
	ZERO(c);		// best of no proposals so far

	do			// iterate over time-steps
	:: s <= ROUNDS*4 ->
		printf("%d step %d\n", j, s);

		// Send <s,p,e,c> and get reply from threshold of recorders
		recs = 0;	// number of recorders we've heard from
		mask = 0;	// bit mask of those recorders
		ZERO(g);	// gather best response proposer saw so far
		done = true;
		select (i : 1 .. N);	// first recorder to interact with
		do		// interact with the recorders in any order
		:: recs < T && (mask & (1 << i)) == 0 ->

			atomic {
				// enter the recorder's role (via "RPC").

				// first catch up the recorder if appropriate
				if
				:: s > rec[i].s ->

					// determine proposal to adopt
					COPY(rec[i].p, p);
					if
					:: (s & 3) == 0 ->
						// Choose a fitness/priority
						if
						:: j == leader ->
							rec[i].p.fit = HI;
						:: else ->
							select(t : 1 .. RAND);
							rec[i].p.fit = t*N + i;
						fi
					:: else -> skip
					fi

					// determine E state to adopt
					if
					:: (s & 3) == 2 && s == rec[i].s+1 ->
						BEST(rec[i].e, e);
					:: else ->
						COPY(rec[i].e, e);
					fi

					// determine C state to adopt
					if
					:: (s & 3) == 3 && s == rec[i].s+1 ->
						BEST(rec[i].c, c);
					:: else ->
						COPY(rec[i].c, c);
					fi

					rec[i].s = s;

				:: else -> skip
				fi
				assert(rec[i].p.val > 0 && rec[i].p.fit > 0);

				// accumulate best proposals in recorder
				if
				:: s == rec[i].s && (s & 3) == 1 ->
					assert(p.fit > 0 && p.val > 0);
					BEST(rec[i].e, p);	// spread E
					assert(rec[i].e.fit > 0 && rec[i].e.val > 0);
				:: s == rec[i].s && (s & 3) == 2 ->
					assert(p.fit > 0 && p.val > 0);
					BEST(rec[i].c, p);	// spread C
				:: else -> skip
				fi

				// we're back to the proposer's logic now,
				// incorporating the recorder's "response".
				assert(s <= rec[i].s);
				if
				:: s == rec[i].s && (s & 3) == 0 ->
					BEST(g, rec[i].p);	// gather P
					assert(g.fit > 0 && g.val > 0);
					done = done && (rec[i].p.fit == HI);

				:: s == rec[i].s && (s & 3) == 1 ->
					skip

				:: s == rec[i].s && (s & 3) == 2 ->
					BEST(g, rec[i].e);	// gather E
					// XXX detect early decide condition

				:: s == rec[i].s && (s & 3) == 3 ->
					BEST(g, rec[i].c);	// gather C

				:: s < rec[i].s ->	// catch up proposer
					s = rec[i].s;
					COPY(p, rec[i].p);
					COPY(e, rec[i].e);
					COPY(c, rec[i].c);
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
				assert(g.fit > 0 && g.val > 0);
				COPY(p, g);	// best of some E set

				// Decide early if all proposals were HI fit
				if
				:: done ->
					DECIDE(j, s, p);
				:: else -> skip
				fi

			:: (s & 3) == 1 ->
				skip

			:: (s & 3) == 2 ->	// gatherEspreadC phase
				// p is now the best of a U set;
				// g is the best of all gathered E sets
				assert(g.fit > 0 && g.val > 0);
				if
				:: p.fit == g.fit ->
					assert(p.val == g.val);
					DECIDE(j, s, p);
				:: else -> skip
				fi

			:: (s & 3) == 3 ->	// gatherC phase
				// g is the best of all gathered C sets.
				// this is our proposal for the next round.
				assert(g.fit > 0 && g.val > 0);
				COPY(p, g);
				ZERO(e);	// start fresh E sets
				ZERO(c);	// start fresh C sets
			fi
			s = s + 1;
			break;
		od

	:: s > ROUNDS*4 ->	// we've simulated enough time-steps
			break;
	od
}

init {
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

