package Game

import (
	"./DB"
	"fmt"
)

type Game struct {
	gameStage       int
	catastropheId   uint8
	numberOfPlayers int
	players         []Player
}

func (g Game) GetNumberOfPlayers() int {
	return g.numberOfPlayers
}

func (g *Game) SetGameStage(n int) {
	g.gameStage = n
}

func (g Game) GetGameStage() int {
	return g.gameStage
}

func (g *Game) AddPlayer(p Player) {
	if g.players == nil {
		g.players = make([]Player, 24)
		g.players[0] = p
		g.numberOfPlayers = 1
		return
	}
	g.players[g.numberOfPlayers] = p
	g.numberOfPlayers++
}

func (g Game) GetPlayers() []Player {
	return g.players
}

func (g *Game) NewGame() error {
	if g.numberOfPlayers < 0 {
		return fmt.Errorf("Not enough players. The minimum number of players is 6. ")
	}

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

func (g *Game) RemovePlayer(p int)  {
	for i:= range g.players {
		if g.players[i].userId==p{
			g.players=append(g.players[:i],g.players[i+1:]...)
			g.numberOfPlayers--
			return
		}
	}
}

func (g *Game) FinishGame()  {
	g.players=nil
	g.numberOfPlayers=0
	g.gameStage=-1
}