/*
file: nnc.go
author: alemedeiros <alexandre.n.medeiros _at_ gmail.com>

A n-sized noughts and crosses game library.

It is a generalization of noughts and crosses, with a n x n board.
To win, you have to fill a line, column or diagonal with your symbol.
*/

// Package nnc implements a n-sized noughts and crosses game.
package nnc

import "errors"

// Empty is an unplayed square;
// Cross is a 'X';
// Nought is a 'O';
const (
	Empty  byte = ' '
	Cross  byte = 'X'
	Nought byte = 'O'
)

// A Game is a game board, use New function to initialize a Game.
type Game struct {
	board      [][]byte
	size       int
	currPlayer byte
}

// Board method returns a copy of the current state of the board.
func (g Game) Board() (board [][]byte) {
	board = make([][]byte, g.size)

	for i := range board {
		board[i] = make([]byte, g.size)
		copy(board[i], g.board[i])
	}

	return
}

// New function Initializes a game structure with a sz-sized board.
// First player is always Cross.
func New(sz int) (g Game) {
	// Allocate a new Game structure
	g = Game{
		board:      make([][]byte, sz),
		size:       sz,
		currPlayer: Cross, // First player is Cross
	}

	for i := range g.board {
		g.board[i] = make([]byte, sz)
		for j := range g.board[i] {
			g.board[i][j] = Empty
		}
	}

	return
}

// Play method checks if the coordinates are inside the board and if it is the
// given player's turn.
//
// Return true and winner (Empty means draw) if the move ended the game.
func (g *Game) Play(x, y int, player byte) (done bool, winner byte, err error) {
	// Validation check
	if g.currPlayer != player {
		return false, Empty, errors.New("not player's turn")
	}
	if x < 0 || g.size <= x || y < 0 || g.size <= y {
		return false, Empty, errors.New("invalid position")
	}

	// Move is valid, do it!
	g.board[x][y] = player

	// Check if move ended the game
	if isDone, winner := g.isDone(); isDone {
		return true, winner, nil
	}

	g.updateTurn()

	return false, Empty, nil
}

// PlayAI method checks if is the given player's turn, if so, it makes a move as
// that player.
//
// Return true and winner (Empty means draw) if the move ended the game.
func (g *Game) PlayAI(player byte) (done bool, winner byte, err error) {
	// Validation check
	if g.currPlayer != player {
		return false, Empty, errors.New("not player's turn")
	}

	// TODO: AI code here

	g.updateTurn()

	return false, Empty, nil
}

// updateTurn method updates whose turn it is.
//
// Assumes the turn was correctly set before call.
func (g *Game) updateTurn() error {
	switch g.currPlayer {
	case Cross:
		g.currPlayer = Nought
	case Nought:
		g.currPlayer = Cross
	default:
		return errors.New("invalid player turn value")
	}

	return nil
}

// isDone method determines if the game is over, and if it is, its winner.
// If winner is Empty, the it was a draw.
func (g Game) isDone() (done bool, winner byte) {
	winner = Empty
	done = true

	// TODO: Check for winner

	// Check for draw
outerFor:
	for i := range g.board {
		for _, p := range g.board[i] {
			if p == Empty {
				done = false
				break outerFor
			}
		}
	}

	return
}

// Outcome calculates the outcome function for a player (Nought/Cross) for the
// current game.
func (g Game) Outcome(player byte) (sum int) {
	if player != Nought && player != Cross {
		return
	}

	for i, sz := 0, g.size; i < sz; i++ {
		// Lines
		linit, lsum := Empty, 0
		for j := 0; j < sz; j++ {
			// Empty squares don't change the outcome function.
			if g.board[i][j] == Empty {
				continue
			}

			// Initialize initial symbol.
			if linit == Empty {
				linit = g.board[i][j]
			}

			// Different symbols means line sum is 0.
			if g.board[i][j] != linit {
				lsum = 0
				break
			}

			if g.board[i][j] == player {
				lsum += 1 // Increment for player
			} else {
				lsum -= 1 // Decrement for opponent
			}
		}

		// Colums
		cinit, csum := Empty, 0
		for j := 0; j < sz; j++ {
			// Empty squares don't change the outcome function.
			if g.board[j][i] == Empty {
				continue
			}

			// Initialize initial symbol.
			if cinit == Empty {
				cinit = g.board[j][i]
			}

			// Different symbols means column sum is 0.
			if g.board[j][i] != cinit {
				csum = 0
				break
			}

			if g.board[j][i] == player {
				csum += 1 // Increment for player
			} else {
				csum -= 1 // Decrement for opponent
			}
		}
		sum += lsum + csum
	}

	// Diagonal
	dinit, dsum := Empty, 0
	for i, sz := 0, g.size; i < sz; i++ {

		// Empty squares don't change the outcome function.
		if g.board[i][i] == Empty {
			continue
		}

		// Initialize initial symbol.
		if dinit == Empty {
			dinit = g.board[i][i]
		}

		// Different symbols means diagonal sum is 0.
		if g.board[i][i] != dinit {
			dsum = 0
			break
		}

		if g.board[i][i] == player {
			dsum += 1 // Increment for player
		} else {
			dsum -= 1 // Decrement for opponent
		}
	}
	sum += dsum

	// Anti-Diagonal
	adinit, adsum := Empty, 0
	for i, sz := 0, g.size; i < sz; i++ {
		// Empty squares don't change the outcome function.
		if g.board[i][sz-1-i] == Empty {
			continue
		}

		// Initialize initial symbol.
		if adinit == Empty {
			adinit = g.board[i][sz-1-i]
		}

		// Different symbols means anti-diagonal sum is 0.
		if g.board[i][sz-1-i] != adinit {
			adsum = 0
			break
		}

		if g.board[i][sz-1-i] == player {
			adsum += 1 // Increment for player
		} else {
			adsum -= 1 // Decrement for opponent
		}
	}
	sum += adsum

	return
}
