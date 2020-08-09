package Game

import (
	"./DB"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
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

func (p Player) MsgId() int {
	return p.msgId
}

func (p *Player) SetMsgId(msgId int) {
	p.msgId = msgId
}

func (p Player) AgainstVotes() int {
	return p.againstVotes
}

func (p *Player) NullifyVotes() {
	p.againstVotes = 0
}

func (p *Player) IncrementAgainstVotes() {
	p.againstVotes++
}

func (p *Player) IsPhobiaOpen() bool {
	return p.isPhobiaOpen
}

func (p *Player) OpenPhobia() {
	p.isPhobiaOpen = true
}

func (p *Player) IsHobbyOpen() bool {
	return p.isHobbyOpen
}

func (p *Player) OpenHobby() {
	p.isHobbyOpen = true
}

func (p *Player) IsBagOpen() bool {
	return p.isBagOpen
}

func (p *Player) OpenBag() {
	p.isBagOpen = true
}

func (p *Player) IsSkillOpen() bool {
	return p.isSkillOpen
}

func (p *Player) OpenSkill() {
	p.isSkillOpen = true
}

func (p *Player) IsHealthOpen() bool {
	return p.isHealOpen
}

func (p *Player) OpenHealth() {
	p.isHealOpen = true
}

func (p *Player) IsBioOpen() bool {
	return p.isBioOpen
}

func (p *Player) OpenBio() {
	p.isBioOpen = true
}

func (p *Player) IsCharOpen() bool {
	return p.isCharOpen
}

func (p *Player) OpenChar() {
	p.isCharOpen = true
}

func (p *Player) SetUser(user *tgbotapi.User) {
	p.user = user
}

func (p Player) GetUser() *tgbotapi.User {
	return p.user
}

func (p Player) GetUserId() int {
	return p.user.ID
}

func (p Player) GetUserName() string {
	if p.user.UserName == "" {
		return p.GetFullName()
	}
	return p.user.UserName
}

func (p Player) GetFullName() string {
	return p.user.FirstName + " " + p.user.LastName
}

func (p Player) GetProfId() uint8 {
	return p.professionId
}

func (p Player) GetCharacterId() uint8 {
	return p.characterId
}
func (p Player) GetBioChar(lang string) string {
	bio := ""
	switch lang {
	case "ru":
		if p.sex {
			bio = "мужчина, "
		} else {
			bio = "женщина, "
		}
		bio += fmt.Sprintf("%d", p.age) + " лет"
	case "en":
		if p.sex {
			bio = "man, "
		} else {
			bio = "woman, "
		}
		bio += fmt.Sprintf("%d", p.age) + " years old"
	}
	return bio
}
func (p Player) GetHealthId() uint8 {
	return p.healthId
}
func (p Player) GetSckillId() uint8 {
	return p.skillId
}
func (p Player) GetBaggageId() uint8 {
	return p.baggageId
}

func (p Player) GetHobbyId() uint8 {
	return p.hobbyId
}
func (p Player) GetPhobiasId() uint8 {
	return p.phobiasId
}

func (p Player) IsAlive() bool {
	return p.alive
}

func (p *Player) Kill() {
	p.againstVotes = 0
	p.alive = false
}

func (p Player) GetActionIds() []uint8 {
	return p.actionId
}

func (player *Player) GenPlayer() {

	db := DB.GetDataBase()
	query, err := db.Query(`SELECT Profession.id, Character.id, Hobby.id, Phobias.id, Skills.id,
							  	  Health.id, Baggage.id 
								  FROM Profession, Character, Hobby, Phobias, Skills,
								  Health, Baggage 
								  ORDER BY random() limit 1`)
	if err != nil {
		panic(err)
	}
	if query.Next() {
		query.Scan(&player.professionId, &player.characterId, &player.hobbyId,
			&player.phobiasId, &player.skillId, &player.healthId, &player.baggageId)
	}
	player.age = int(rand.Uint32() % 90)
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
		query.Scan(&player.actionId[i])

	}

	player.alive = true
}

func (p Player) countProfit(catastropheId uint8) int {
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
Baggage_profit.baggage_id = $5 and Hobby_profit.hobby_id = $6 and Phobias_profit.phobia_id = $7 and Skills_profit.skill_id = &8`,catastropheId, p.professionId, p.healthId, p.characterId, p.baggageId, p.hobbyId, p.phobiasId, p.skillId)
	if err != nil {
		panic(err)
	}
	var profits []int
	tmp, _ := query.Columns()
	profits = make([]int, len(tmp))
	if query.Next() {
		query.Scan(&profits[0], &profits[1], &profits[2], &profits[3], &profits[4], &profits[5], &profits[6], &profits[7])
		for i := range profits {
			profit += profits[i]
		}

		return profit
	}
	return 0
}
