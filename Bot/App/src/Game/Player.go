package Game

import (
	"../loger"
	"./DB"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"time"
)

type GameInfo struct {
	age          int
	sex          bool
	professionId uint8
	characterId  uint8
	isCharOpen   bool
	isBioOpen    bool
	healthId     uint8
	isHealOpen   bool
	skillId      uint8
	isSkillOpen  bool
	baggageId    uint8
	isBagOpen    bool
	actionId     []uint8
	hobbyId      uint8
	isHobbyOpen  bool
	phobiasId    uint8
	isPhobiaOpen bool
	alive        bool
}

type Player struct {
	user         *tgbotapi.User
	againstVotes int
	msgId        int
	GameInfo
}

const (
	NUMBER_OF_CHARACTERISTICS = 8
)

func (player Player) MsgId() int {
	return player.msgId
}

func (player *Player) SetMsgId(msgId int) {
	player.msgId = msgId
}

func (player Player) AgainstVotes() int {
	return player.againstVotes
}

func (player *Player) NullifyVotes() {
	player.againstVotes = 0
}

func (player *Player) IncrementAgainstVotes() {
	player.againstVotes++
}

func (player *Player) IsPhobiaOpen() bool {
	return player.isPhobiaOpen
}

func (player *Player) OpenPhobia() {
	player.isPhobiaOpen = true
}

func (player *Player) IsHobbyOpen() bool {
	return player.isHobbyOpen
}

func (player *Player) OpenHobby() {
	player.isHobbyOpen = true
}

func (player *Player) IsBagOpen() bool {
	return player.isBagOpen
}

func (player *Player) OpenBag() {
	player.isBagOpen = true
}

func (player *Player) IsSkillOpen() bool {
	return player.isSkillOpen
}

func (player *Player) OpenSkill() {
	player.isSkillOpen = true
}

func (player *Player) IsHealthOpen() bool {
	return player.isHealOpen
}

func (player *Player) OpenHealth() {
	player.isHealOpen = true
}

func (player *Player) IsBioOpen() bool {
	return player.isBioOpen
}

func (player *Player) OpenBio() {
	player.isBioOpen = true
}

func (player *Player) IsCharOpen() bool {
	return player.isCharOpen
}

func (player *Player) OpenChar() {
	player.isCharOpen = true
}

func (player *Player) SetUser(user *tgbotapi.User) {
	player.user = user
}

func (player Player) GetUser() *tgbotapi.User {
	return player.user
}

func (player Player) GetUserId() int {
	return player.user.ID
}

func (player Player) GetUserName() string {
	if player.user.UserName == "" {
		return player.GetFullName()
	}
	return player.user.UserName
}

func (player Player) GetFullName() string {
	return player.user.FirstName + " " + player.user.LastName
}

func (player Player) GetProfId() uint8 {
	return player.professionId
}

func (player Player) GetCharacterId() uint8 {
	return player.characterId
}
func (player Player) GetBioChar(lang string) string {
	bio := ""
	switch lang {
	case "ru":
		if player.sex {
			bio = "мужчина, "
		} else {
			bio = "женщина, "
		}
		bio += fmt.Sprintf("%d", player.age) + " лет"
	case "en":
		if player.sex {
			bio = "man, "
		} else {
			bio = "woman, "
		}
		bio += fmt.Sprintf("%d", player.age) + " years old"
	}
	return bio
}
func (player Player) GetHealthId() uint8 {
	return player.healthId
}
func (player Player) GetSckillId() uint8 {
	return player.skillId
}
func (player Player) GetBaggageId() uint8 {
	return player.baggageId
}

func (player Player) GetHobbyId() uint8 {
	return player.hobbyId
}
func (player Player) GetPhobiasId() uint8 {
	return player.phobiasId
}

func (player Player) IsAlive() bool {
	return player.alive
}

func (player *Player) Kill() {
	player.againstVotes = 0
	player.alive = false
}

func (player Player) GetActionIds() []uint8 {
	return player.actionId
}

func (player *Player) GenPlayer() {
	rand.Seed(time.Now().UnixNano())
	db := DB.GetDataBase()
	query, err := db.Query(` SELECT COUNT(*) FROM Profession UNION ALL
								SELECT COUNT(*)	FROM Health UNION ALL
								SELECT COUNT(*) FROM Hobby	UNION ALL
								SELECT COUNT(*)	FROM Phobias UNION ALL
								SELECT COUNT(*)	FROM Skills UNION ALL
								SELECT COUNT(*)	FROM Baggage UNION ALL
								SELECT COUNT(*)	From Character`)

	if err != nil {
		loger.LogErr(err)
	}

	if query.Next() {
		err=query.Scan(&player.professionId)
		loger.LogErr(err)
		player.professionId = uint8(rand.Intn(int(player.professionId)-1)) + 1
	}

	if query.Next() {
		err=query.Scan(&player.healthId)
		loger.LogErr(err)
		player.healthId = uint8(rand.Intn(int(player.healthId)-1)) + 1
	}

	if query.Next() {
		err=query.Scan(&player.hobbyId)
		loger.LogErr(err)
		player.hobbyId = uint8(rand.Intn(int(player.hobbyId)-1)) + 1
	}

	if query.Next() {
		err=query.Scan(&player.phobiasId)
		loger.LogErr(err)
		player.phobiasId = uint8(rand.Intn(int(player.phobiasId)-1)) + 1
	}

	if query.Next() {
		err=query.Scan(&player.skillId)
		loger.LogErr(err)
		player.skillId = uint8(rand.Intn(int(player.skillId)-1)) + 1
	}

	if query.Next() {
		err=query.Scan(&player.baggageId)
		loger.LogErr(err)
		player.baggageId = uint8(rand.Intn(int(player.baggageId)-1)) + 1
	}
	if query.Next() {
		err=query.Scan(&player.characterId)
		loger.LogErr(err)
		player.characterId = uint8(rand.Intn(int(player.characterId)-1)) + 1
	}

	if query.Next() {
		err=query.Scan(&player.professionId)
		loger.LogErr(err)
		player.professionId = uint8(rand.Intn(int(player.professionId)-1)) + 1
	}

	player.age = rand.Intn(90)
	if player.age < 18 {
		player.age += int(rand.Uint32() % 20)
	}
	player.sex = rand.Int()%2 == 0
	query, err = db.Query("SELECT Action.id FROM Action ORDER BY RANDOM() limit 2")
	if err != nil {
		panic(err)
	}
	player.actionId = make([]uint8, 3)
	for i := 0; query.Next() && i < 3; i++ {
		err=query.Scan(&player.actionId[i])
		loger.LogErr(err)
	}

	player.alive = true
}

func (player Player) countProfit(catastropheId uint8) int {
	var profit int
	profit = 0
	db := DB.GetDataBase()
	query, err := db.Query(`SELECT Profession_profit.profit, Health_profit.profit, Character_profit.profit, Baggage_profit.profit,
								Hobby_profit.profit, Phobias_profit.profit, Skills_profit.profit
								FROM Profession_profit JOIN Health_profit ON Profession_profit.catastrophe_id=Health_profit.catastrophe_id  
								JOIN  Character_profit ON Character_profit.catastrophe_id = Profession_profit.catastrophe_id
								JOIN Baggage_profit ON Baggage_profit.catastrophe_id = Profession_profit.catastrophe_id
								JOIN Hobby_profit ON Hobby_profit.catastrophe_id = Profession_profit.catastrophe_id
								JOIN Phobias_profit ON Phobias_profit.catastrophe_id = Phobias_profit.catastrophe_id
								JOIN Skills_profit ON Skills_profit.catastrophe_id = Skills_profit.catastrophe_id
								WHERE Profession_profit.catastrophe_id = $1 and Profession_profit.profession_id = $2 and
								Health_profit.health_id = $3 and Character_profit.character_id = $4 and
								Baggage_profit.baggage_id = $5 and Hobby_profit.hobby_id = $6 and Phobias_profit.phobia_id = $7 and Skills_profit.skill_id = $8`, catastropheId, player.professionId, player.healthId, player.characterId, player.baggageId, player.hobbyId, player.phobiasId, player.skillId)
	if err != nil {
		panic(err)
	}
	var profits []int
	tmp, _ := query.Columns()
	profits = make([]int, len(tmp))
	if query.Next() {
		err=query.Scan(&profits[0], &profits[1], &profits[2], &profits[3], &profits[4], &profits[5], &profits[6], &profits[7])
		loger.LogErr(err)
		for i := range profits {
			profit += profits[i]
		}

		return profit
	}
	return 0
}
