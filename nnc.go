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
	count      int
	currPlayer byte
}

// Structure to save the move and its value.
type move struct {
	value, i, j int
}

// CurrentPlayer method returns the player that should play.
func (g Game) CurrentPlayer() byte {
	return g.currPlayer
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

// Get the minimum weighted playing position.
func min(a, b move) move {
	if a.value <= b.value {
		return a
	} else {
		return b
	}
}

// Get the maximum weighted playing position.
func max(a, b move) move {
	if a.value >= b.value {
		return a
	} else {
		return b
	}
}

// New function Initializes a game structure with a sz-sized board.
// First player is always Cross.
func New(sz int) (g Game) {
	// Allocate a new Game structure
	g = Game{
		board:      make([][]byte, sz),
		size:       sz,
		count:      sz * sz,
		currPlayer: Cross, // First player is Cross
	}

	// Initialize board.
	for i := range g.board {
		g.board[i] = make([]byte, sz)
		for j := range g.board[i] {
			g.board[i][j] = Empty
		}
	}

	return
}

// Return a copy of the current game.
func (g Game) copyGame() (ng Game) {
	// Allocate a new Game structure
	ng = Game{
		board:      make([][]byte, g.size),
		size:       g.size,
		count:      g.count,
		currPlayer: g.currPlayer,
	}

	// Copy board.
	for i := range ng.board {
		ng.board[i] = make([]byte, g.size)
		for j := range ng.board[i] {
			ng.board[i][j] = g.board[i][j]
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
	if g.board[x][y] != Empty {
		print("error position: ", x, " ", y, "\n")
		return false, Empty, errors.New("cell already played")
	}

	// Move is valid, do it!
	g.board[x][y] = player

	// Check if move ended the game
	isDone, winner := g.isDone()

	g.updateTurn()
	g.count -= 1

	return isDone, winner, nil
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

	// A value greater than the maximum value possible for a game.
	lim := g.size * g.size * 10

	// Serial alpha-beta pruning
	m := alphaBetaPruningSerial(*g, g.size*g.size, -lim, lim, -1, -1, player)

	//res := make(chan move)
	//prune := make(chan struct{})
	//defer close(prune)

	//go alphaBetaPruning(*g, g.size*g.size, -lim, lim, -1, -1, player, res, prune)

	//// Wait for result.
	//m := <-res

	return g.Play(m.i, m.j, player)
}

// Serial implementation of Alpha-Beta Pruning algorithm.
// TODO: Try not to copy the entire game structure
func alphaBetaPruningSerial(g Game, depth int, alpha, beta int, x, y int, player byte) move {
	// Check for depth limit or if game is over
	if depth == 0 {
		return move{g.outcome(player), x, y}
	}
	if done, _ := g.isDone(); done {
		return move{g.outcome(player), x, y}
	}

	// Check for whose turn it is
	if curr := g.currPlayer; curr == player {
		p := move{alpha, x, y}

		for i, l := range g.board {
			for j, e := range l {
				// Check for possible move
				if e != Empty {
					continue
				}

				// Generate updated game
				ng := g.copyGame()
				ng.Play(i, j, player)

				m := alphaBetaPruningSerial(ng, depth-1, alpha, beta, i, j, player)
				m.i = i
				m.j = j

				// Update alpha
				p = max(p, m)
				alpha = p.value

				// Beta cut-off
				if beta <= alpha {
					return m
				}
			}
		}
		return p
	} else {
		p := move{beta, x, y}

		for i, l := range g.board {
			for j, e := range l {
				// Check for possible move
				if e != Empty {
					continue
				}

				// Generate updated game
				ng := g.copyGame()
				ng.Play(i, j, curr)

				m := alphaBetaPruningSerial(ng, depth-1, alpha, beta, i, j, player)
				m.i = i
				m.j = j

				// Update beta
				p = min(p, m)
				beta = p.value

				// Alpha cut-off
				if beta <= alpha {
					return m
				}
			}
		}
		return p
	}
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
	var local bool
	var init byte

	// Check for winner
	for i, sz := 0, g.size; i < sz; i++ {
		// Lines
		local = true
		init = Empty
		for j := 0; j < sz && local; j++ {
			if j == 0 {
				init = g.board[i][j]
			}

			if g.board[i][j] == Empty || g.board[i][j] != init {
				local = false
			}
		}

		// Return if someone won
		if local {
			return local, init
		}

		// Columns
		local = true
		init = Empty
		for j := 0; j < sz && local; j++ {
			if j == 0 {
				init = g.board[j][i]
			}

			if g.board[j][i] == Empty || g.board[j][i] != init {
				local = false
			}
		}

		// Return if someone won
		if local {
			return local, init
		}
	}

	// Diagonal
	local = true
	init = Empty
	for i, sz := 0, g.size; i < sz && local; i++ {
		if i == 0 {
			init = g.board[i][i]
		}

		if g.board[i][i] == Empty || g.board[i][i] != init {
			local = false
		}
	}

	// Return if someone won
	if local {
		return local, init
	}

	// Anti-diagonal
	local = true
	init = Empty
	for i, sz := 0, g.size; i < sz && local; i++ {
		if i == 0 {
			init = g.board[i][sz-1-i]
		}

		if g.board[i][sz-1-i] == Empty || g.board[i][sz-1-i] != init {
			local = false
		}
	}

	// Return if someone won
	if local {
		return local, init
	}

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
func (g Game) outcome(player byte) (sum int) {
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

		if lsum == sz || csum == sz {
			return 3 * sz * sz
		} else if lsum == -sz || csum == -sz {
			return -(3 * sz * sz)
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

	if dsum == g.size {
		return 3 * g.size * g.size
	} else if dsum == -g.size {
		return -(3 * g.size * g.size)
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

	if adsum == g.size {
		return 3 * g.size * g.size
	} else if adsum == -g.size {
		return -(3 * g.size * g.size)
	}

	sum += adsum

	return
}
