# Socket documentation

Before a game can be played, the users must be connected to the socket server on the endpoint.
```
    ws://<host>:<api>
```
Communication between the client and the server will occure via the socket instantiated on this endpoint.
e.g
```
    url = 'ws://loclhost:8081'
    ws = new Websocket(url)
```

## Messaging

### Authentication
The `first` message that must be sent over the socket to the server is an authentication message.
```
message = {
    type: 'auth',
    authMessage: {
        token: string
    }
}
ws.emit(socketMessageAuth)
```
The token is gotten from logging in with username and password. `/login`

The response of `auth` message is an `authResponse`.
```
authResponse = {
    type: 'auth',
    error: ''  // omitted if no error is returned
}
```

### Ping
A ping is simply used to test a connection to the server.
```
message = {
    type: 'ping',
    pingMessage: {} //can be ignored
}
```
Response is a `pong`. You can use a timeout with this ping request to assert a successful connection to the server.
```
pingResponse = {
    type: 'pong'
    pong: 'pong'
}
```

### Opponent Found

When an opponent is found after a `search opponent` request, the server will respond with 
```
opponentFound = {
    type: 'opponentFound',
    username: 'string' // The username of the opponent
}
```
The client should navigate to the game page once this message is received because shortly afterwards,
the server will start the game and begin sending questions.

### Question

When an opponent is found, the server will wait a few seconds for the both clients to navigate to the game page and become **ready** to play.
Then, the server will start the game and begin issuing questions.

```
message = {
    type: 'question',
    questionMessage: {
        time: int,        // The time the question was sent
        round: int,         // The question round. e.g round 4 of 10
        rounds: int,        // The total number of rounds. e.g 10
        question: stirng,   // The question itself
        answers: [string],  // the answer options of the question
        correctAnswer: string // the correct answer. An element of `answers 
    }
}
```

### Answer
When a question is received, the client will respond with an answer.
```
message = {
    type: 'answer',
    answerMessage: {
        round: int,
        answer: string,
        question: string
    }
}
```
N.B The client should respond immediatly he has the answer as the time of the response will determine the score.


### Rounds
A round gets finished when,

 - **both** clients send the `answers` within the time limit. `10 seconds`
 - The 10 seconds timeout elapses.
 
When a round is over the server will send a `round result` to both clients. 
```
roundResult = {
    type: 'roundResult'
    round: int,
    question: string,
    username1: {
        answer: string,
        time: number,
        score: float
    },
    username2: {
        answer: string,
        time: number,
        score: float
    }
}
```


### Game finished

When the game finishes, i.e all rounds hve been completed, the server will respond with `game finishsed`
```
gameFinished = {
    type: 'gameFinished',
    rounds: []roundResult
    totals: {
        username1: float,
        username2: float
    }
}
```

### Leaving a game
A client can leave the game if the game is not finished.
The client must send a `quit` message.

```
message = {
    type: 'quit'
    quitMessage: {} //can be ignored
}
```
If the game is not finished, the other player wins.
A `gameFinished` response will be received.

### Unexpected quit

If a client abruptly ends the connection to the sever, he will loose the game.
A `gameFinished` response will be received.




