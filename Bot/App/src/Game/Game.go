package Game

import (
	"./DB"
	"fmt"
)

type Game struct {
	alive                int
	destruction          int
	gameStage            int
	catastropheId        uint8
	numberOfPlayers      int
	numberOfAlivePlayers int
	players              []Player
}

const (
	GAME_IS_OVER=-1
	REGISTRATION = 0
	SELECTING_CHARACTERISTICS=1
	VOTING = 2
	DISCUSSION = 3
)

func (g *Game) NumberOfAlivePlayers() int {
	return g.numberOfAlivePlayers
}

func (g *Game) SetNumberOfAlivePlayers(numberOfAlivePlayers int) {
	g.numberOfAlivePlayers = numberOfAlivePlayers
}

func (g Game) GetAlive() int {
	return g.alive
}

func (g *Game) SetAlive(n int) {
	g.alive = n
}

func (g *Game) SetDestruction(n int) {
	g.destruction = n
}

func (g Game) GetDestruction() int {
	return g.destruction
}

func (g Game) GetCatastropheId() uint8 {
	return g.catastropheId
}

func (g Game) GetNumberOfPlayers() int {
	return g.numberOfPlayers
}

func (g *Game) IncrementNumberOfVoted() {
	g.numberOfPlayers++
}

func (g *Game) SetGameStage(n int) {
	g.gameStage = n
}

func (g Game) GetGameStage() int {
	return g.gameStage
}

func (g *Game) AddPlayer(p Player) {
	if g.players == nil {
		g.players = make([]Player, 0, 24)
	}
	g.players = append(g.players, p)
	g.numberOfPlayers++
}

func (g Game) GetPlayers() []Player {
	return g.players
}

func (g *Game) NewGame() error {
	if g.numberOfPlayers < 0 {
		return fmt.Errorf("Not enough players. The minimum number of players is 6. ")
	}

	g.gameStage = 1
	g.numberOfAlivePlayers = g.numberOfPlayers
	for i := 0; i < g.numberOfPlayers; i++ {
		g.players[i].GenPlayer()
	}
	db := DB.GetDataBase()
	query, err := db.Query("SELECT id FROM Catastrophe_en ORDER BY random() limit 1")
	if err != nil {
		return err
	}
	query.Next()
	err = query.Scan(&g.catastropheId)
	if err != nil {
		return err
	}
	return nil
}

func (g Game) CountProfit() int {
	var profit = 0
	for i := range g.players {
		if g.players[i].IsAlive() {
			profit += g.players[i].countProfit(g.catastropheId)
		}
	}
	return profit
}

func (g Game) FindByID(id int) Player {
	for i := 0; i < g.numberOfPlayers; i++ {
		if g.players[i].user.ID == id {
			return g.players[i]
		}
	}
	return Player{}
}

func (g *Game) RemovePlayer(p int) {
	for i := 0; i < g.numberOfPlayers; i++ {
		if g.players[i].user.ID == p {
			g.players = append(g.players[:i], g.players[i+1:]...)
			g.numberOfPlayers--
			return
		}
	}
}

func (g *Game) FinishGame() {
	g.players = nil
	g.numberOfPlayers = 0
	g.gameStage = -1
}
