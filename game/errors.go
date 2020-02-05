package game

import (
	"errors"
)

//ErrorNotExist thrown when a player with ID was searched for and not found
type ErrorNotExist error

//ErrorExists thrown when a player with ID was already present
type ErrorExists error

//ErrorAlreadyLost is thrown when a player attempts to forfeit while already at 0
type ErrorAlreadyLost error

//ErrorNoBidder is thrown when there is no valid next bidder available
// note that the player who bid last cannot bid again
type ErrorNoBidder error

//ErrorBidFace is thrown when the proposed bid face is above 6 or below 0
type ErrorBidFace error

//ErrorBidQuantity is thrown when the proposed bid quantity is below 0 or higher than the number of dice in the game
type ErrorBidQuantity error

//ErrorBidTooLow is thrown when the proposed bid is valid but not higher than the current one
// The returned error contains the current bid.
type ErrorBidTooLow struct {
	error
	Quantity uint
	Face     uint
}

//ErrorOutOfTurn is thrown when an action is attempted out of valid turn order
type ErrorOutOfTurn error

//ErrorNotEnoughPlayers is thrown when there are not enough players to start a game
type ErrorNotEnoughPlayers error

func errNotExist() ErrorNotExist {
	return errors.New("player does not exist")
}

func errExists() ErrorExists {
	return errors.New("player already exists")
}

func errAlreadyLost() ErrorAlreadyLost {
	return errors.New("player already lost")
}

func errNoBidder() ErrorNoBidder {
	return errors.New("no valid bidder found")
}

func errBidFace() ErrorBidFace {
	return errors.New("proposed bid face was invalid")
}

func errBidQuantity() ErrorBidQuantity {
	return errors.New("proposed bid quantity was invalid")
}

func errBidTooLow(quantity uint, face uint) ErrorBidTooLow {
	return ErrorBidTooLow{
		error:    errors.New("proposed bid was too low"),
		Quantity: quantity,
		Face:     face,
	}
}

func errOutOfTurn() ErrorOutOfTurn {
	return errors.New("action attempted by non-bidding player")
}

func errNotEnoughPlayers() ErrorNotEnoughPlayers {
	return errors.New("not enough players to start the round")
}
