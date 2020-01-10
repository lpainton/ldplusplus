package table

import (
	"fmt"
)

type Face uint8
const (
	ONE Face = iota
	TWO Face = iota
	THREE Face = iota
	FOUR Face = iota
	FIVE Face = iota
	SIX Face = iota
)

type ErrorOutOfTurn error
type ErrorBidTooLow error

type player string
type order uint

// bid is encoded as (quantity * 6) + face.
// This gives us a single number that is monotonically increasing in both dimensions
type bid uint

//Table reflects a single table where a game is in progress
type Table struct {
	// Current person whose turn it is to bid
	turn   order
	last bid
	players []player
	wilds map[Face]bool
	hands   map[player][6]uint
	dice    map[player]uint
}

//New creates a new table running liar's dice
func New() *Table {
	return &Table{
		turn: order(0),
		last: bid(0),
		players: make([]player, 0),
		wilds:   make(map[Face]bool),
		hands:   make(map[player][6]uint),
		dice:    make(map[player]uint),
	}
}

/* 
* Play adds the player to the table at the end of the current order
* if an only if they aren't already present.
*
* Returns true if a change was made
*/ 
func (t *Table) Play(p player) bool {
	for _, v := range t.players {
		if p == v {
			return false
		}
	}
	t.players = append(t.players, p)
	return true
}

/* 
* Bid changes the to the new bid if and only if
* the turn matches the current player and
* it is strictly greater than the last bid. It also
* advances the turn by one
*
* Returns an error or nil if sucessful
*/
func (t *Table) Bid(p player, quantity uint, f Face) error {
	if p != t.players[t.turn] {
		return ErrorOutOfTurn(fmt.Errorf("player out of turn was %s",p))
	}

	b := bid((quantity * 6) + uint(f))
	if b <= t.last {
		return ErrorBidTooLow(fmt.Errorf("bid was %d",b))
	}

	t.last = b
	t.turn = order((int(t.turn) + 1)%len(t.players))
	return nil
}

/*
* Lair returns true if the last bid is inconsistent
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
		for i,v := range h {
			totals[i] += v
		}
	}

	//Make sure we don't double count if someone bet wilds
	toCount := map[Face]bool{f: true}
	for k,v := range t.wilds {
		toCount[k] = v
	}

	//Now we count up everything
	var count uint
	for k,v := range toCount {
		if v {
			count += totals[k]
		}
	}

	quantity := uint(t.last / 6)
	return count >= quantity
}