package socketserver

import (
	"errors"
	"sync"
	"time"

	questionService "github.com/acha-bill/quizzer_backend/packages/dblayer/question"
	"github.com/labstack/gommon/log"
)

var (
	searchingMutex sync.Mutex
	gameManager    *GManager
	gManagerOnce   sync.Once
)

var (
	// ErrNotEnoughQuestions is returned if there are not enough questions to launch a game
	ErrNotEnoughQuestions         = errors.New("not enough questions")
	ErrPlayerAlreadyInAnotherGame = errors.New("player already in another game")
	ErrAlreadySearching           = errors.New("already searching")
)

// ServerManager is the game manager
type GManager struct {
	searching []*WsConnection
	games     []*Game
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
		for _, game := range mgr.games {
			if game.Active == condition {
				res = append(res, game)
			}
		}
	} else {
		res = mgr.games
	}

	return res
}

// FindPlayerGame finds the game that contains the specified player
func (mgr *GManager) FindPlayerGame(player *WsConnection) *Game {
	for _, game := range mgr.games {
		for _, p := range game.Players {
			if p == player {
				return game
			}
		}
	}
	return nil
}

// NewGame creates and starts a new game between 2 players and a specified number of questions.
// If there are not upto `length` questions available, it returns with error.
// If any of the members is already in a game, it returns with error.
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
	mgr.games = append(mgr.games, g)

	// Tell players game is about to start
	ServerManager().WriteConnection(player1, NewSocketResponseOpponentFound(player2.Context.User.Username))
	ServerManager().WriteConnection(player2, NewSocketResponseOpponentFound(player1.Context.User.Username))

	// wait a bit for players to prepare
	time.Sleep(3 * time.Second)

	g.Start()
	return nil
}

// AddSearcher adds a player to the searching queue
// If the player is already searching, it returns with error.
func (mgr *GManager) AddSearcher(player *WsConnection) error {
	searchingMutex.Lock()
	defer searchingMutex.Unlock()

	isAlreadySearching := false
	for _, p := range mgr.searching {
		if p == player {
			isAlreadySearching = true
			break
		}
	}
	if isAlreadySearching {
		return ErrAlreadySearching
	}
	mgr.searching = append(mgr.searching, player)
	return nil
}

// GetPair returns the 1st two players in the searching queue.
// It nil for both if the pair cannot be formed.
func (mgr *GManager) GetPair() (player1 *WsConnection, player2 *WsConnection) {
	if len(mgr.searching) < 2 {
		return
	}
	player1, player2 = mgr.searching[0], mgr.searching[1]
	mgr.searching = mgr.searching[2:]
	return
}

const opponentFoundType = "opponentFound"

// SocketResponseOpponentFound represents the user found by search
type SocketResponseOpponentFound struct {
	Type     string `json:"Type"`
	Username string `json:"username"`
}

// NewSocketResponseOpponentFound creates a new NewSocketResponseOpponentFound
func NewSocketResponseOpponentFound(username string) SocketResponseOpponentFound {
	return SocketResponseOpponentFound{
		Type:     opponentFoundType,
		Username: username,
	}
}
