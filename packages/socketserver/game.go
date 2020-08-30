package socketserver

import (
	"errors"
	"time"

	"github.com/acha-bill/quizzer_backend/models"
)

// RoundResult represents the result of a game round
type RoundResult struct {
	QuestionIndex int
	Question      string
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
	Winnner      string
}

var (
	ErrGameIsStillRunning = errors.New("game is still running. Try again with force option")
)

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

// PrematureLoose looses the player because he disconnected from the server or quit the game.
func (game *Game) PrematureLoose(player *WsConnection) {
	game.Active = false
	var winner *WsConnection
	for _, p := range game.Players {
		if p != player {
			winner = p
			break
		}
	}
	if winner == nil {
		return
	}
	response := NewSocketResponseGameFinished(game)
	game.Winnner = winner.Context.User.Username
	response.Winner = game.Winnner
	broadcast(game, response)
	_ = game.Close(true)
}

// SetRoundResult sets the result submitted by a player for a particular round.
func (game *Game) SetRoundResult(player *WsConnection, questionIndex int, answer string, timeReceived time.Time) {
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
	roundResult := game.RoundResults[round]
	for _, player := range game.Players {
		ServerManager().WriteConnection(player, NewSocketResponseRoundResult(roundResult))
	}
	go nextRound(game)
}

// sends a message to all players
func broadcast(game *Game, msg interface{}) {
	for _, player := range game.Players {
		ServerManager().WriteConnection(player, msg)
	}
}

// Close closes the game and removes all players
func (game *Game) Close(force ...bool) error {
	f := len(force) > 0 && force[0]

	if game.Active && !f {
		return ErrGameIsStillRunning
	}
	game.Active = false
	game.Players = []*WsConnection{}
	return nil
}

// nextRound starts a new round of the game
func nextRound(game *Game) {
	round := game.Cursor
	if round >= len(game.Questions) {
		game.Active = false
		broadcast(game, NewSocketResponseGameFinished(game))
		return
	}

	// prepare result
	roundResult := &RoundResult{
		QuestionIndex: round,
		Question:      game.Questions[round].Question,
		Answers:       make(map[*WsConnection]string),
		Times:         make(map[*WsConnection]time.Time),
		Scores:        make(map[*WsConnection]float64),
	}
	game.RoundResults = append(game.RoundResults, roundResult)

	timeSent := time.Now()
	//send question to players
	broadcast(game, NewSocketResponseQuestion(timeSent.Unix(), round, game.Questions[round].Question, game.Questions[round].Answers, game.Questions[round].CorrectAnswer))
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
	if !g.Active {
		return
	}
	g.SetRoundResult(wsConnection, answer.Round, answer.Answer, time.Now())
}

// handleQuitMessage handles a quit message
func handleQuitMessage(connection *WsConnection, _ SocketMessageQuit) {
	g := GameManager().FindPlayerGame(connection)
	if !g.Active {
		return
	}
	g.PrematureLoose(connection)
}

const responseQuestionType = "question"
const responseGameFinishedType = "gameFinished"
const responseRoundResultType = "roundResult"

// SocketResponseQuestion represents a question
type SocketResponseQuestion struct {
	Type          string   `json:"type"`
	Time          int64    `json:"time"`
	Round         int      `json:"round"`
	Question      string   `json:"question"`
	Answers       []string `json:"answers"`
	CorrectAnswer string   `json:"correctAnswer"`
}

type Result struct {
	Answer string    `json:"string"`
	Time   time.Time `json:"time"`
	Score  float64   `json:"score"`
}

// SocketResponseRoundResult is the result of players of a round
type SocketResponseRoundResult struct {
	Type     string             `json:"type"`
	Round    int                `json:"round"`
	Question string             `json:"question"`
	Results  map[string]*Result `json:"results"`
}

// NewSocketResponseRoundResult returns a NewSocketResponseRoundResult
func NewSocketResponseRoundResult(result *RoundResult) SocketResponseRoundResult {
	results := make(map[string]*Result)
	for player, answer := range result.Answers {
		results[player.Context.User.Username] = &Result{
			Answer: answer,
		}
	}
	for player, t := range result.Times {
		results[player.Context.User.Username].Time = t
	}
	for player, score := range result.Scores {
		results[player.Context.User.Username].Score = score
	}
	return SocketResponseRoundResult{
		Type:     responseRoundResultType,
		Round:    result.QuestionIndex,
		Question: result.Question,
		Results:  results,
	}
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

// SocketMessageQuit represents a quit game message
type SocketMessageQuit struct {
}

// SocketResponseGameFinished the response returned when a round is finished.
type SocketResponseGameFinished struct {
	Type         string                      `json:"type"`
	RoundResults []SocketResponseRoundResult `json:"roundResults"`
	Totals       map[string]float64          `json:"totals"`
	Winner       string                      `json:"winner"`
}

// NewSocketResponseGameFinished returns a new NewSocketResponseGameFinished
func NewSocketResponseGameFinished(game *Game) SocketResponseGameFinished {
	var roundResults []SocketResponseRoundResult
	totals := make(map[string]float64)
	for _, roundResult := range game.RoundResults {
		round := NewSocketResponseRoundResult(roundResult)
		roundResults = append(roundResults, round)
		for username, result := range round.Results {
			totals[username] += result.Score
		}
	}

	maxScore := 0.0
	winner := ""
	for username, totalScore := range totals {
		if totalScore > maxScore {
			winner = username
		}
	}
	game.Winnner = winner
	return SocketResponseGameFinished{
		Type:         responseGameFinishedType,
		RoundResults: roundResults,
		Totals:       totals,
		Winner:       winner,
	}
}
