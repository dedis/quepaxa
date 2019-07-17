
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


proctype NodeProc(byte i) {
	byte j, r, s, tkt, step, seen, scnt, prsn, best, btkt;
	byte belig, betkt, beseen, k;
	//bool correct = (n < T);

	//printf("Node %d correct %d\n", n, correct);

	for (r : 0 .. ROUNDS-1) {

		atomic {

		// select a "random" (here just arbitrary) ticket
		select (tkt : 1 .. TICKETS);
		node[i].round[r].ticket = tkt;

		// we've already seen our own proposal
		prsn  = 1 << i;

		// finding the "best proposal" starts with our own...
		best = i;
		btkt = tkt;

		} // atomic

		// Run the round to completion
		for (s : 0 .. STEPS-1) {

			// "send" the broadcast for this time-step
			node[i].round[r].sent[s] = 1;

			// collect a threshold of other nodes' broadcasts
			seen = 1 << i;		// we've already seen our own
			scnt = 1;
			do
			::	// Pick another node to try to 'receive' from
				select (j : 1 .. N); j--;
				if
				:: ((seen & (1 << j)) == 0) && 
				    (node[j].round[r].sent[s] != 0) ->

					atomic {

					//printf("%d received from %d\n", n, j);
					seen = seen | (1 << j);
					scnt++;

					// Track the best proposal we've seen
					if
					:: step == 0 ->
						prsn = prsn | (1 << j);
						if
						:: node[j].round[r].ticket > btkt ->
							best = j;
							btkt = node[j].round[r].ticket;
						:: node[j].round[r].ticket == btkt ->
							best = 255;	// tied tickets
						:: else -> skip
						fi

					// Track proposals we've seen indirectly
					:: step > 0 ->
						prsn = prsn | node[j].round[r].prsn[s-1];
						if
						:: node[j].round[r].btkt[s-1] > btkt ->
							best = node[j].round[r].best[s-1];
							btkt  = node[j].round[r].btkt[s-1];
						:: (node[j].round[r].btkt[s-1] == btkt) &&
							(node[j].round[r].best[s-1] != best) ->
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
			node[i].round[r].seen[s] = seen;
			node[i].round[r].prsn[s] = prsn;
			node[i].round[r].best[s] = best;
			node[i].round[r].btkt[s] = btkt;

			printf("%d step %d complete: seen %x best %d ticket %d\n",
				i, s, seen, best, btkt);

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
		for (j : 0 .. N-1) {

			// determine number of nodes that knew of j's proposal
			// by time t+2.
			int jseen = 0;
			for (k : 0 .. N-1) {
				if
				:: ((node[i].round[r].seen[2] & (1 << k)) != 0) &&
					((node[k].round[r].prsn[1] & (1 << j)) != 0) ->
					jseen++;
				:: else ->
					skip
				fi
			}

			if
			:: (jseen >= Fa+1) &&	// j's proposal is eligible
			   (node[j].round[r].ticket > betkt) -> // is better
				belig = j;
				betkt = node[j].round[r].ticket;
				beseen = jseen;
			:: (jseen >= Fa+1) &&	// j's proposal is eligible
			   (node[j].round[r].ticket == betkt) -> // is tied
				belig = 255;
				beseen = 0;
			:: else -> skip
			fi
		}
		printf("%d best eligible proposal %d ticket %d seen by %d\n",
			i, belig, betkt, beseen);

		// we should have found at least one eligible proposal!
		assert(betkt > 0);

		// The round is now complete in terms of picking a proposal.
		node[i].round[r].picked = belig;
		node[i].round[r].done = 1;

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
			printf("%d round %d definitely committed\n", i, r);

			// Verify that what we decided doesn't conflict with
			// the proposal any other node chooses.
			select (j : 1 .. N); j--;
			assert(!node[j].round[r].done ||
				(node[j].round[r].picked == belig));

		:: (belig < 255) && (beseen < T) ->
			printf("%d round %d failed due to threshold\n", i, r);

		:: (belig < 255) && (belig != best) ->
			printf("%d round %d failed due to spoiler\n", i, r);

		:: (belig == 255) ->
			printf("%d round %d failed due to tie\n", i, r);
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

