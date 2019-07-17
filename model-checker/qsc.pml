
#define N		3		// total number of nodes
#define Fa	1		// max number of availability failures
#define Fc	0		// max number of correctness failures
#define T	(Fa+Fc+1)	// consensus threshold required

#define STEPS	3		// TLC time-steps per consensus round
#define ROUNDS	2		// number of consensus rounds to run
#define TICKETS	3		// proposal lottery ticket space

// TLC state for each logical time-step
typedef Step {
	bit sent;			// true if we've sent our raw proposal
	bit seen[1+N];		// nodes whose raw proposals we've received
	bit ackd[1+N];		// nodes who have acknowledged our raw proposal
	bit witd;			// true if our proposal is threshold witnessed
	bit witn[1+N];		// nodes we've gotten threshold witnessed msgs from
}

// QSC summary information for a "best" proposal seen so far
typedef Best {
	byte from;			// node number the proposal is from, 0 if tied spoiler
	byte tkt;			// proposal's genetic fitness ticket value
}

// TLC and QSC state per round
typedef Round {
	Step step[STEPS];	// TLC state for each logical time-step

	byte ticket;		// QSC lottery ticket assigned to proposal at t+0
	Best spoil;			// best potential spoiler(s) we've found so far
	Best conf;			// best confirmed proposal we've seen so far
	Best reconf;		// best reconfirmed proposal we've seen so far
	byte picked;		// which proposal this node picked this round, 0 if not yet
}

// Per-node state
typedef Node {
	Round rnd[ROUNDS];	// each node's per-consensus-round information
}

Node node[1+N];			// all state of each node 1..N


// Implement a given node i.
proctype NodeProc(byte i) {
	byte j, r, s, tkt, step, acks, wits;

	for (r : 0 .. ROUNDS-1) {

		atomic {
			// select a "random" (here just arbitrary) ticket
			select (tkt : 1 .. TICKETS);
			node[i].rnd[r].ticket = tkt;

			// start with our own proposal as best potential spoiler
			node[i].rnd[r].spoil.from = i;
			node[i].rnd[r].spoil.tkt = tkt;
		} // atomic

		// Run the round to completion
		for (s : 0 .. STEPS-1) {

			// "send" the broadcast for this time-step
			node[i].rnd[r].step[s].sent = 1;

			// collect a threshold of other nodes' broadcasts
			acks = 0;
			wits = 0;
			do
			::	// Pick another node to "receive" a message from
				select (j : 1 .. N);
				atomic {

					// Track the best potential spoiler we encounter
					if
					// Node j knows about a strictly better potential spoiler
					:: node[j].rnd[r].spoil.tkt > node[i].rnd[r].spoil.tkt ->
						node[i].rnd[r].spoil.from = node[j].rnd[r].spoil.from;
						node[i].rnd[r].spoil.tkt = node[j].rnd[r].spoil.tkt;

					// Node j knows about a spoiler that's tied with our best
					:: node[j].rnd[r].spoil.tkt == node[i].rnd[r].spoil.tkt &&
						node[j].rnd[r].spoil.from != node[i].rnd[r].spoil.from ->
						node[i].rnd[r].spoil.from = 0; // tied, so mark invalid

					:: else -> skip
					fi

					// Track the best confirmed proposal we encounter
					if
					:: node[j].rnd[r].conf.tkt > node[i].rnd[r].conf.tkt ->
						node[i].rnd[r].conf.from = node[j].rnd[r].conf.from;
						node[i].rnd[r].conf.tkt = node[j].rnd[r].conf.tkt;
					:: else -> skip
					fi

					// Track the best reconfirmed proposal we encounter
					if
					:: node[j].rnd[r].reconf.tkt > node[i].rnd[r].reconf.tkt ->
						node[i].rnd[r].reconf.from = node[j].rnd[r].reconf.from;
						node[i].rnd[r].reconf.tkt = node[j].rnd[r].reconf.tkt;
					:: else -> skip
					fi

					// Now handle specific types of messages: Raw, Ack, or Wit.
					if

					// We "receive" a raw unwitnessed message from node j
					:: node[j].rnd[r].step[s].sent && !node[i].rnd[r].step[s].seen[j] ->

						node[i].rnd[r].step[s].seen[j] = 1;

					// We "receive" an acknowledgment of our message from node j
					:: node[j].rnd[r].step[s].seen[i] && !node[i].rnd[r].step[s].ackd[j] ->

						node[i].rnd[r].step[s].ackd[j] = 1;
						acks++;
						if
						:: acks >= T ->
							// Our proposal is now fully threshold witnessed
							node[i].rnd[r].step[s].witd = 1

							// See if our proposal is now the best confirmed proposal
							if
							:: s == 0 &&
								node[i].rnd[r].ticket > node[i].rnd[r].conf.tkt ->
								node[i].rnd[r].conf.from = i;
								node[i].rnd[r].conf.tkt = node[i].rnd[r].ticket;
							:: else -> skip
							fi

							// See if we're reconfirming a best confirmed proposal
							if
							:: s == 1 &&
								node[i].rnd[r].conf.tkt > node[i].rnd[r].reconf.tkt ->
								node[i].rnd[r].reconf.from = node[i].rnd[r].conf.from;
								node[i].rnd[r].reconf.tkt = node[i].rnd[r].conf.tkt;
							:: else -> skip
							fi

						:: else -> skip
						fi

					// We "receive" a fully threshold witnessed message from node j
					:: node[j].rnd[r].step[s].witd && !node[i].rnd[r].step[s].witn[j] ->

						node[i].rnd[r].step[s].witn[j] = 1
						wits++;

					// End this step if we've seen enough witnessed proposals
					:: wits >= T -> break;

					:: else -> skip
					fi
				} // atomic
			od
		}

		atomic {
			printf("%d best spoiler %d ticket %d\n",
						i, node[i].rnd[r].spoil.from, node[i].rnd[r].spoil.tkt);
			printf("%d best confirmed %d ticket %d\n",
						i, node[i].rnd[r].conf.from, node[i].rnd[r].conf.tkt);
			printf("%d best reconfirmed %d ticket %d\n",
						i, node[i].rnd[r].reconf.from, node[i].rnd[r].reconf.tkt);

			// The round is now complete in terms of picking a proposal.
			node[i].rnd[r].picked = node[i].rnd[r].conf.from;

			// We can be sure everyone has converged on this proposal
			// if it is also the best spoiler and best reconfirmed proposal.
			if
			:: node[i].rnd[r].spoil.from == node[i].rnd[r].picked &&
				node[i].rnd[r].reconf.from == node[i].rnd[r].picked ->
				printf("%d round %d definitely COMMITTED\n", i, r);

				// Verify that what we decided doesn't conflict with
				// the proposal any other node chooses.
				select (j : 1 .. N);
				assert(!node[j].rnd[r].picked ||
					(node[j].rnd[r].picked == node[i].rnd[r].picked));

			:: node[i].rnd[r].reconf.from != node[i].rnd[r].picked ->
				printf("%d round %d FAILED to be reconfirmed\n", i, r);

			:: node[i].rnd[r].spoil.from != node[i].rnd[r].picked ->
				printf("%d round %d FAILED due to spoiler\n", i, r);

			:: node[i].rnd[r].spoil.from == 0 ->
				printf("%d round %d FAILED due to tie\n", i, r);

			:: else ->
			fi
		} // atomic
	}
}

init {
	atomic {
		int i;
		for (i : 1 .. N) {
			run NodeProc(i)
		}
	}
}

