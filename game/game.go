package game

import (
	"fmt"
	"math/rand"
	"time"
)

const defaultStartingHand = 5

type ErrorGameStarted error
type ErrorAlreadyPlaying error
type ErrorBidOutOfTurn error

type Face uint8

const (
	One   Face = iota
	Two   Face = iota
	Three Face = iota
	Four  Face = iota
	Five  Face = iota
	Six   Face = iota
)

type State uint8

const (
	Open State = iota
	Play State = iota
	Done State = iota
)

type person string
type place uint

// bid is encoded as (quantity * 6) + face.
// This gives us a single number that is monotonically increasing in both dimensions
type bid uint

// A quantity of dice
type dice uint

func successor(o order, bound int) order {
	return order((int(o) + 1) % bound)
}

//Table reflects a single table where a game is in progress
type Table struct {
	// Order of people at the table
	people []person
	// Map of people to dice left
	dice map[person]dice
	// The starting hand size
	startingHand dice
	// Set of faces which are wild at this table
	wilds map[Face]bool
	// The last loser
	loser person
	// The current round
	*round
	// Random number source
	*rand.Rand
}

//Round represents a single round from rolls to bidding to calling liar
type round struct {
	turn    place
	latest  bid
	players []person
	hands   map[person][6]uint
}

//New creates a new table running liar's dice
func New() *Table {
	return &Table{
		people: make([]person, 0),
		dice:   make(map[person]dice),
		startingHand: defaultStartingHand,
		wilds:  make(map[Face]bool),
		Rand:   rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (t *Table) newRound(offset place) {
	//Find the position of the last loser
	var s place
	for l := place(len(t.people)); s < l && (t.people[s] != t.loser); s++ {
	}

	//Build a subset of players with dice left
	// we start with the person who was the last loser
	var players []person
	for i, l := s, place(len(t.people)); i < l+s; i++ {
		p := t.people[i%l]
		if t.dice[p] > 0 {
			players = append(players, p)
		}
	}

	//Roll starting hands
	hands := make(map[person][6]uint)
	for _, p := range players {
		var d [6]uint
		for h := t.dice[p]; h > 0; h-- {
			d[t.Intn(5)]++
		}
		hands[p] = d
	}

	r := round{
		turn:    place(0),
		latest:  bid(0),
		players: players,
		hands:   hands,
	}
}

//Play adds the person to the table at the end of the current order.
// throws an error if they were already playing
// converts them to playing if they were a spectator
func (t *Table) Play(p person) error {
	if 

	for _, r := range t.players {
		if p == r {
			if t.dice[p] > 0 {
				return fmt.Errorf("%s already playing", p)
			}
			t.dice[p] = t.startingHand
			return nil
		}
	}
	t.players = append(t.players, p)
	t.dice[p] = t.startingHand
	return nil
}

//Watch adds the player to the table as a spectator.
//Throws an error if they were already watching.
function (t *Table) Watch(p person) error {
	for _, r := range t.players {
		if p == r {
			return fmt.Errorf("%s already watching or playing", p)
		}
	}
	t.players = append(t.players, p)
	t.dice[p] = t.startingHand
	return nil
}

/*
* Bid changes the last bid to the new bid if and only if
* the turn matches the current player and
* it is strictly greater than the last bid. It also
* advances the turn by one
*
* Returns an error or nil if sucessful
 */
func (t *Table) Bid(p player, quantity uint, f Face) error {
	if p != t.players[t.turn] {
		return ErrorOutOfTurn(fmt.Errorf("player out of turn was %s", p))
	}

	b := bid((quantity * 6) + uint(f))
	if b <= t.last {
		return ErrorBidTooLow(fmt.Errorf("bid was %d", b))
	}

	t.last = b
	t.turn = successor(t.turn, len(t.players))
	return nil
}

/*
* Liar returns true if the last bid is inconsistent
* with the game state.
*
* Consistency means that the number of wilds + Faces bid
* is greater than the quantity
 */
func (t *Table) Liar() bool {
	f := Face(t.last % 6)

	//Sum all hands for a total
	var totals [6]uint
	for _, h := range t.hands {
		for i, v := range h {
			totals[i] += v
		}
	}

	//Make sure we don't double count if someone bet wilds
	toCount := map[Face]bool{f: true}
	for k, v := range t.wilds {
		toCount[k] = v
	}

	//Now we count up everything
	var count uint
	for k, v := range toCount {
		if v {
			count += totals[k]
		}
	}

	quantity := uint(t.last / 6)
	return count >= quantity
}
