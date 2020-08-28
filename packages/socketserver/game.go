package socketserver

import (
	"time"

	"github.com/acha-bill/quizzer_backend/models"
)

// RoundResult represents the result of a game round
type RoundResult struct {
	QuestionIndex int
	Answers       map[*WsConnection]string
	Times         map[*WsConnection]time.Time
	Scores        map[*WsConnection]float64
}

// Game represents the game between players
type Game struct {
	Active       bool
	Players      []*WsConnection
	Questions    []*models.Question
	Cursor       int
	RoundResults []*RoundResult
	RoundTimes   []time.Time
}

// newGame creates a new game
func newGame(player1 *WsConnection, player2 *WsConnection, questions []*models.Question) *Game {
	g := &Game{
		Active:     false,
		Players:    []*WsConnection{player1, player2},
		Questions:  questions,
		Cursor:     0,
		RoundTimes: make([]time.Time, len(questions)),
	}
	return g
}

// Start starts the game
func (game *Game) Start() {
	game.Active = true
	go nextRound(game)
}

// SetRoundResult sets the result submitted by a player for a particular round.
func (game *Game) SetRoundResult(player *WsConnection, questionIndex int, answer string, timeReceived time.Time) {
	if !game.Active {
		return
	}
	roundResult := game.RoundResults[questionIndex]
	roundResult.Answers[player] = answer
	roundResult.QuestionIndex = questionIndex
	roundResult.Times[player] = timeReceived
	timeTakenSecs := timeReceived.Sub(game.RoundTimes[questionIndex]).Seconds()
	score := 10 - timeTakenSecs
	if score < 0 {
		score = 0
	}
	roundResult.Scores[player] = score

	isRoundComplete := true
	for _, player := range game.Players {
		if _, ok := roundResult.Answers[player]; !ok {
			isRoundComplete = false
			break
		}
	}
	if isRoundComplete {
		go finalizeAndGoToNextRound(game, questionIndex)
	}
}

// finalizeAndGoToNextRound broadcasts the result of the current round and starts the next round
func finalizeAndGoToNextRound(game *Game, round int) {
	for _, player := range game.Players {
		ServerManager().WriteConnection(player, game.RoundResults[round])
	}
	go nextRound(game)
}

// nextRound starts a new round of the game
func nextRound(game *Game) {
	round := game.Cursor
	if round >= len(game.Questions) {
		game.Active = false
		// game finished!
		return
	}

	// prepare result
	roundResult := &RoundResult{
		QuestionIndex: round,
		Answers:       make(map[*WsConnection]string),
		Times:         make(map[*WsConnection]time.Time),
		Scores:        make(map[*WsConnection]float64),
	}
	game.RoundResults = append(game.RoundResults, roundResult)

	timeSent := time.Now()
	//send question to players
	for _, player := range game.Players {
		ServerManager().WriteConnection(player, NewSocketResponseQuestion(timeSent.Unix(), round, game.Questions[round].Question, game.Questions[round].Answers, game.Questions[round].CorrectAnswer))
	}
	game.RoundTimes[round] = timeSent

	time.AfterFunc(11*time.Second, func() {
		// nobody has given an answer after 11 secs
		if len(game.RoundResults[round].Answers) == 0 {
			go finalizeAndGoToNextRound(game, round)
		}
	})
	game.Cursor++
}

// handleAnswerMessage handles the answer message for a question
func handleAnswerMessage(wsConnection *WsConnection, answer SocketMessageAnswer) {
	//find game that is running with this connection
	g := GameManager().FindPlayerGame(wsConnection)
	g.SetRoundResult(wsConnection, answer.Round, answer.Answer, time.Now())
}

const responseQuestionType = "question"

// SocketResponseQuestion represents a question
type SocketResponseQuestion struct {
	Type          string   `json:"type"`
	Time          int64    `json:"time"`
	Round         int      `json:"round"`
	Question      string   `json:"question"`
	Answers       []string `json:"answers"`
	CorrectAnswer string   `json:"correctAnswer"`
}

// NewSocketResponseQuestion returns a new NewSocketResponseQuestion
func NewSocketResponseQuestion(time int64, round int, question string, answers []string, correctAnswer string) SocketResponseQuestion {
	return SocketResponseQuestion{
		Type:          responseQuestionType,
		Time:          time,
		Round:         round,
		Question:      question,
		Answers:       answers,
		CorrectAnswer: correctAnswer,
	}
}

// SocketMessageAnswer represents the answer to a specific round of the questions.
type SocketMessageAnswer struct {
	Round    int    `json:"round"`
	Answer   string `json:"answer"`
	Question string `json:"question"`
}
