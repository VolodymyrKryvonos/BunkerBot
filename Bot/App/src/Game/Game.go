package Game

import (
	"./DB"
	"fmt"
)

type Game struct {
	gameStage            int
	catastropheId        uint8
	numberOfPlayers      int
	numberOfAlivePlayers int
	bunkerCap            int
	players              []Player
	PlayersToKick        []*Player
}

func (g *Game) BunkerCap() int {
	return g.bunkerCap
}

func (g *Game) SetBunkerCap(bunkerCap int) {
	g.bunkerCap = bunkerCap
}

const (
	GAME_IS_OVER              = -1
	REGISTRATION              = 0
	SELECTING_CHARACTERISTICS = 1
	VOTING                    = 2
	DISCUSSION                = 3
)

func (g *Game) NumberOfAlivePlayers() int {
	return g.numberOfAlivePlayers
}

func (g *Game) SetNumberOfAlivePlayers(numberOfAlivePlayers int) {
	g.numberOfAlivePlayers = numberOfAlivePlayers
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
	if g.numberOfPlayers <= 0 {
		g.gameStage=GAME_IS_OVER
		return fmt.Errorf("Not enough players. The minimum number of players is 6. ")
	}

	g.gameStage = 1
	g.numberOfAlivePlayers = g.numberOfPlayers
	g.bunkerCap = g.numberOfAlivePlayers/2
	for i := 0; i < g.numberOfPlayers; i++ {
		g.players[i].GenPlayer()
	}
	db := DB.GetDataBase()
	query, err := db.Query("SELECT id FROM Catastrophe ORDER BY random() limit 1")
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

func (g Game) FindByID(id int) *Player {
	for i := 0; i < g.numberOfPlayers; i++ {
		if g.players[i].user.ID == id {
			return &g.players[i]
		}
	}
	return &Player{}
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

func (g *Game) Kick() {
	playersToKick := make([]*Player, 0, g.numberOfAlivePlayers)
	maxVotes := 0
	for i := 0; i < len(g.PlayersToKick); i++ {
		if g.players[i].againstVotes > maxVotes {
			maxVotes = g.players[i].againstVotes
		}
	}

	for i := 0; i < len(g.PlayersToKick); i++ {
		if g.players[i].againstVotes == maxVotes {
			playersToKick = append(playersToKick, &g.players[i])
		}
	}
	g.PlayersToKick = playersToKick
}

func (g *Game) FinishGame() {
	g.players = nil
	g.numberOfPlayers = 0
	g.gameStage = -1
}
