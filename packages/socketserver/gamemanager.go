package socketserver

import (
	"errors"
	"sync"

	questionService "github.com/acha-bill/quizzer_backend/packages/dblayer/question"
	"github.com/labstack/gommon/log"
)

var (
	games        []*Game
	gameManager  *GManager
	gManagerOnce sync.Once
)

var (
	// ErrNotEnoughQuestions is returned if there are not enough questions to launch a game
	ErrNotEnoughQuestions         = errors.New("not enough questions")
	ErrPlayerAlreadyInAnotherGame = errors.New("player already in another game")
)

// ServerManager is the game manager
type GManager struct {
}

// GameManager returns the manager instance
func GameManager() *GManager {
	gManagerOnce.Do(func() {
		gameManager = &GManager{}
	})
	return gameManager
}

// Games returns the list of games
func (mgr *GManager) Games(isActive ...bool) []*Game {
	var res []*Game
	if len(isActive) > 0 {
		condition := isActive[0]
		for _, game := range games {
			if game.Active == condition {
				res = append(res, game)
			}
		}
	} else {
		res = games
	}

	return res
}

// FindPlayerGame finds the game that contains the specified player
func (mgr *GManager) FindPlayerGame(player *WsConnection) *Game {
	for _, game := range games {
		for _, p := range game.Players {
			if p == player {
				return game
			}
		}
	}
	return nil
}

// NewGame creates and starts a new game between 2 players and a specified number of questions
func (mgr *GManager) NewGame(player1 *WsConnection, player2 *WsConnection, length int) error {
	// check that players are not already in another game
	if found1, found2 := mgr.FindPlayerGame(player1), mgr.FindPlayerGame(player2); found1 != nil || found2 != nil {
		return ErrPlayerAlreadyInAnotherGame
	}
	// find questions
	questions, err := questionService.FindAll()
	if err != nil {
		log.Info(err)
		return err
	}
	if len(questions) < length {
		return ErrNotEnoughQuestions
	}

	questions = questions[:length]
	g := newGame(player1, player2, questions)
	games = append(games, g)
	g.Start()
	return nil
}
