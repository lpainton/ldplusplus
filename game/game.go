/*
* Package game provides a complete model for a single game of liars dice.
* For our purposes a game is a discrete play session starting from a fixed set of players who all have dice and ending when only one player has dice remaining.
* Players
 */
package game

import (
	"math/rand"
	"time"
)

type prng interface {
	Intn(int) int
}

//Monotonic maps two-component vector bid to monotonic space
func Monotonic(quantity uint, face uint) uint {
	return (quantity * 6) + face
}

//Factorize maps a monotinc bid to two-component space
func Factorize(bid uint) (quantity uint, face uint) {
	quantity = bid / 6
	face = bid % 6
	return
}

//Player is a player object with associated reporting fields
type Player struct {
	ID   string
	Dice uint
	Hand [6]uint
}

//Rules of the game
type Rules struct {
	Dice  uint    //Number of dice players start with
	Wilds [6]bool //Array of faces which are wild
}

type hand = [6]uint
type playerID = string

//Game is a full game from start to finish
type Game struct {
	players []playerID
	cups    map[playerID]uint
	hands   map[playerID]hand
	bidder  int  //The current bidder
	bid     uint //The previous bid
	prev    int  //The previous bidder
	prng
	Rules
}

//Pending is a completely new game that hasn't started play yet
type Pending struct{ Game }

//RoundStart is a game which is at the beginning of the round. Dice have been rolled but no bids made
type RoundStart struct{ Game }

//Bidding is a game where dice have been rolled, bids have been made but nobody has called liar
type Bidding struct{ Game }

//LiarCalled is one where someone has called liar and either won or lost
type LiarCalled struct{ Game }

//Over is a game where only one player has dice remaining
type Over struct{ Game }

//New creates a pending game from the provided set of rules
func New(r Rules) *Pending {
	return &Pending{Game{
		prng:  rand.New(rand.NewSource(time.Now().Unix())),
		Rules: r,
	}}
}

//decides if the player with id exists in the game
func (g Game) exists(id playerID) bool {
	for _, p := range g.players {
		if p == id {
			return true
		}
	}
	return false
}

//StartGame starts a pending game
func (p Pending) StartGame(bidder int) (*RoundStart, error) {
	//Ensure there are enough players to play a round
	if len(p.players) < 2 {
		return nil, errNotEnoughPlayers()
	}

	//Add dice to cups
	var cups = make(map[playerID]uint)
	for _, id := range p.players {
		cups[id] = p.Dice
	}

	//Roll starting hands
	h := rollAll(cups, p)

	return &RoundStart{Game{
		players: p.players,
		cups:    cups,
		hands:   h,
		prng:    p.prng,
		Rules:   p.Rules,
	}}, nil
}

//NewRound starts a new round post liar being called
func (l LiarCalled) NewRound(bidder int) (*RoundStart, error) {
	//Ensure there are enough players to play a round
	var count uint
	for _, c := range l.cups {
		if c > 0 {
			count++
		}
	}
	if count > 0 {
		return nil, errNotEnoughPlayers()
	}

	return &RoundStart{l.Game}, nil
}

//rolls all dice in the provided cups
func rollAll(cups map[playerID]uint, rng prng) map[playerID]hand {
	m := make(map[playerID]hand)
	for p, d := range cups {
		m[p] = roll(d, rng)
	}
	return m
}

//rolls a new hand of dice of the given size given the provided rng
func roll(dice uint, rng prng) hand {
	var h [6]uint
	for i := uint(0); i < dice; i++ {
		h[rng.Intn(5)]++
	}
	return h
}

//returns the index of the next bidder or an error if there is no valid bidder left
func (g *Game) next() (int, error) {
	for i := 1; i < len(g.players); i++ {
		b := (g.bidder + i) % len(g.players)
		if g.cups[g.players[b]] > 0 {
			return b, nil
		}
	}
	return g.bidder, errNoBidder()
}

//returns the number of dice in the game
func (g *Game) dice() uint {
	var sum uint
	for _, c := range g.cups {
		sum += c
	}
	return sum
}

//counts the number of dice in the game showing a particular face
func (g *Game) count(face int) uint {
	var sum uint
	for _, p := range g.players {
		sum += g.hands[p][face]
	}
	return sum
}

//Player finds a player in the game by id, returns error if not found
func (g *Game) Player(id string) (*Player, error) {
	if g.exists(id) {
		return &Player{
			ID:   id,
			Dice: g.cups[id],
			Hand: g.hands[id],
		}, nil
	}
	return nil, errNotExist()
}

//Add adds a player to the game iff they aren't already present.
// It throws an error if they are.
func (g *Game) Add(id string) error {
	if g.exists(id) {
		return errExists()
	}
	g.players = append(g.players, id)
	g.cups[id] = g.Rules.Dice
	return nil
}

//sets a player's hand to zero and restarts the round
// the next valid bidder is the starting player
func (g *Game) forfeit(id string) error {
	g.cups[id] = 0
	n, err := g.next()
	if err != nil {
		return err
	}
	return g.round(n)
}

//Forfeit sets a player's hand size to 0 and resets the round
// returns an error if the player already lost or doesn't exist
func (g *Game) Forfeit(id string) error {
	if !g.exists(id) {
		return errNotExist()
	}
	if g.cups[id] == 0 {
		return errAlreadyLost()
	}

	return g.forfeit(id)
}

/*Bid changes the current bid to the new bid iff:
* - The proposing bidder exists
* - The proposing bidder matches the current bidder
* - The proposed bid face is within range [0,6]
* - The proposed bid quantity is within range [0,D] where D is the number of dice
* 	currently in the game
* - The proposed bid increases either the face or quantity of the last bid without
*	decreasing either
*
* A successful bid updates the current bid, notes the bidder and advances the round
*  to the next bidder.
*
* Returns nil if sucessful, otherwise error
 */
func (g *Game) Bid(id string, quantity uint, face uint) error {
	switch {
	case !g.exists(id):
		return errNotExist()
	case id != g.players[g.bidder]:
		return errOutOfTurn()
	case face < 0, face > 6:
		return errBidFace()
	case quantity < 0, quantity > g.dice():
		return errBidQuantity()
	case Monotonic(quantity, face) <= g.bid:
		q, f := Factorize(g.bid)
		return errBidTooLow(q, f)
	}

	n, err := g.next()
	if err != nil {
		return err
	}

	g.prev = g.bidder
	g.bid = Monotonic(quantity, face)
	g.bidder = n
	return nil
}

//LiarResult contains information about the result of a call of Liar
type LiarResult struct {
	Lying     bool
	AccuserID string
	AccusedID string
}

/*Liar returns a LiarResult based on if the last bid is inconsistent
* with the actual game state.
*
* Consistency means that the number of wilds + Faces bid
* is greater than or equal to the quantity. If the bid is inconsistent
* then the result of true.
*
* Liar will subtract a die from either the calling player or previous bidder
* depending on whether the result was false or true respectively.
* Liar also attempts to start a new game round and may return errors related to that
 */
func (g *Game) Liar(id string) (LiarResult, error) {
	if !g.exists(id) {
		return LiarResult{}, errNotExist()
	}
	if id != g.players[g.bidder] {
		return LiarResult{}, errOutOfTurn()
	}

	q, f := Factorize(g.bid)

	//Make sure we don't double count if someone bet wilds
	faces := g.Wilds
	faces[f] = true

	//Now we validate the bid
	for i, v := range faces {
		if v {
			q -= g.count(i)
		}
	}

	result := LiarResult{
		Lying:     q < 0,
		AccusedID: g.players[g.prev],
		AccuserID: id,
	}
	var loser int
	if result.Lying {
		g.cups[g.players[g.prev]]--
		loser = g.prev
	} else {
		g.cups[g.players[g.bidder]]--
		loser = g.bidder
	}

	return result, g.round(loser)
}
