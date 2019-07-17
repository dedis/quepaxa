
#define N	4		// total number of nodes
#define Fa	1		// max number of availability failures
#define Fc	1		// max number of unknown correctness failures
#define T	(Fa+Fc+1)	// consensus threshold required

#define STEPS	3		// TLC time-steps per consensus round
#define ROUNDS	2		// number of consensus rounds to run
#define TICKETS	3		// proposal lottery ticket space


typedef Round {
	bit sent[STEPS];	// whether we've sent yet each time step
	byte ticket;		// lottery ticket assigned to proposal at t+0
	byte seen[STEPS];	// bitmask of msgs we've seen from each step
	byte prsn[STEPS];	// bitmaks of proposals we've seen after each
	byte best[STEPS];
	byte btkt[STEPS];
	byte picked;		// which proposal this node picked this round
	bit done;		// set to true when round complete
}

typedef Node {
	Round round[ROUNDS];	// each node's per-consensus-round information
}

Node node[N];			// all state of each node


proctype NodeProc(byte n) {
	byte rnd, tkt, step, seen, scnt, prsn, best, btkt, nn;
	byte belig, betkt, beseen, k;
	//bool correct = (n < T);

	//printf("Node %d correct %d\n", n, correct);

	for (rnd : 0 .. ROUNDS-1) {

		atomic {

		// select a "random" (here just arbitrary) ticket
		select (tkt : 1 .. TICKETS);
		node[n].round[rnd].ticket = tkt;

		// we've already seen our own proposal
		prsn  = 1 << n;

		// finding the "best proposal" starts with our own...
		best = n;
		btkt = tkt;

		} // atomic

		// Run the round to completion
		for (step : 0 .. STEPS-1) {

			// "send" the broadcast for this time-step
			node[n].round[rnd].sent[step] = 1;

			// collect a threshold of other nodes' broadcasts
			seen = 1 << n;		// we've already seen our own
			scnt = 1;
			do
			::	// Pick another node to try to 'receive' from
				select (nn : 1 .. N); nn--;
				if
				:: ((seen & (1 << nn)) == 0) && 
				    (node[nn].round[rnd].sent[step] != 0) ->

					atomic {

					//printf("%d received from %d\n", n, nn);
					seen = seen | (1 << nn);
					scnt++;

					// Track the best proposal we've seen
					if
					:: step == 0 ->
						prsn = prsn | (1 << nn);
						if
						:: node[nn].round[rnd].ticket > btkt ->
							best = nn;
							btkt = node[nn].round[rnd].ticket;
						:: node[nn].round[rnd].ticket == btkt ->
							best = 255;	// tied tickets
						:: else -> skip
						fi

					// Track proposals we've seen indirectly
					:: step > 0 ->
						prsn = prsn | node[nn].round[rnd].prsn[step-1];
						if
						:: node[nn].round[rnd].btkt[step-1] > btkt ->
							best = node[nn].round[rnd].best[step-1];
							btkt  = node[nn].round[rnd].btkt[step-1];
						:: (node[nn].round[rnd].btkt[step-1] == btkt) &&
							(node[nn].round[rnd].best[step-1] != best) ->
							best = 255;	// tied tickets
						:: else -> skip
						fi
					fi

					} // atomic

				:: else -> skip
				fi

				// Threshold test: have we seen enough?
				if
				:: scnt >= T -> break;
				:: else -> skip;
				fi
			od

			atomic {

			// Record what we've seen for the benefit of others
			node[n].round[rnd].seen[step] = seen;
			node[n].round[rnd].prsn[step] = prsn;
			node[n].round[rnd].best[step] = best;
			node[n].round[rnd].btkt[step] = btkt;

			printf("%d step %d complete: seen %x best %d ticket %d\n",
				n, step, seen, best, btkt);

			} // atomic
		}

		atomic {

		// Find the best propposal we can determine to be eligible.
		// We deem a proposal to be eligible if we can see that
		// it was seen by at least f+1 nodes by time t+1.
		// This ensures that ALL nodes at least know of its existence
		// (though not necessarily its eligibility) by t+2.
		belig = 255;	// start with a fake 'tie' state
		betkt = 0;		// worst possible ticket value
		beseen = 0;
		for (nn : 0 .. N-1) {

			// determine number of nodes that knew of nn's proposal
			// by time t+2.
			int nnseen = 0;
			for (k : 0 .. N-1) {
				if
				:: ((node[n].round[rnd].seen[2] & (1 << k)) != 0) &&
					((node[k].round[rnd].prsn[1] & (1 << nn)) != 0) ->
					nnseen++;
				:: else ->
					skip
				fi
			}

			if
			:: (nnseen >= Fa+1) &&	// nn's proposal is eligible
			   (node[nn].round[rnd].ticket > betkt) -> // is better
				belig = nn;
				betkt = node[nn].round[rnd].ticket;
				beseen = nnseen;
			:: (nnseen >= Fa+1) &&	// nn's proposal is eligible
			   (node[nn].round[rnd].ticket == betkt) -> // is tied
				belig = 255;
				beseen = 0;
			:: else -> skip
			fi
		}
		printf("%d best eligible proposal %d ticket %d seen by %d\n",
			n, belig, betkt, beseen);

		// we should have found at least one eligible proposal!
		assert(betkt > 0);

		// The round is now complete in terms of picking a proposal.
		node[n].round[rnd].picked = belig;
		node[n].round[rnd].done = 1;

		// Can we determine a proposal to be definitely committed?
		// To do so, we must be able to see that:
		//
		// 1. it was seen by t+2 by ALL nodes we have info from.
		// 2. we know of no other proposal competitive with it.
		// 
		// #1 ensures ALL nodes will judge this proposal as eligible;
		// #2 ensures no node could judge another proposal as eligible.
		if
		:: (belig < 255) && (beseen >= T) && (belig == best) ->
			printf("%d round %d definitely committed\n", n, rnd);

			// Verify that what we decided doesn't conflict with
			// the proposal any other node chooses.
			select (nn : 1 .. N); nn--;
			assert(!node[nn].round[rnd].done ||
				(node[nn].round[rnd].picked == belig));

		:: (belig < 255) && (beseen < T) ->
			printf("%d round %d failed due to threshold\n", n, rnd);

		:: (belig < 255) && (belig != best) ->
			printf("%d round %d failed due to spoiler\n", n, rnd);

		:: (belig == 255) ->
			printf("%d round %d failed due to tie\n", n, rnd);
		fi

		} // atomic
	}
}


init {
	atomic {
		int i;
		for (i : 0 .. N-1) {
			run NodeProc(i)
		}
	}
}

