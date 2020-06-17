package Game

import "./DB"

type Player struct {
	userId                     int
	userName                   string
	professionId               uint8
	characterId                uint8
	biologicalCharacteristicId uint8
	healthId                   uint8
	skillId                    uint8
	baggageId                  uint8
	actionId                   []uint8
	hobbyId                    uint8
	phobiasId                  uint8
	alive                      bool
}

func (p *Player) SetUserId(i int) {
	p.userId = i
}

func (p Player) GetUserId() int {
	return p.userId
}

func (p *Player) SetUserName(i string) {
	p.userName = i
}

func (p Player) GetUserName() string {
	return p.userName
}

func (p Player) GetProfId() uint8 {
	return p.professionId
}

func (p Player) GetCharacterId() uint8 {
	return p.characterId
}
func (p Player) GetBioCharId() uint8 {
	return p.biologicalCharacteristicId
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
	p.alive = false
}

func (p Player) GetActionIds() []uint8 {
	return p.actionId
}

func (player *Player) GenPlayer() {

	db := DB.GetDataBase()
	query, err := db.Query(`SELECT Profession_en.id, character_en.id, hobby_en.id, phobias_en.id, skills_en.id,
							  	  health_en.id, biological_characteristics_en.id, baggage_en.id 
								  FROM Profession_en, character_en, hobby_en, phobias_en, skills_en,
								  health_en, biological_characteristics_en, baggage_en 
								  ORDER BY random() limit 1`)
	if err != nil {
		panic(err)
	}
	if query.Next() {
		query.Scan(&player.professionId, &player.characterId, &player.hobbyId,
			&player.phobiasId, &player.skillId, &player.healthId,
			&player.biologicalCharacteristicId, &player.baggageId)
	}
	query, err = db.Query("SELECT action_en.id FROM action_en ORDER BY RANDOM() limit 3")
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
								  Biological_characteristics_profit.profit, Hobby_profit.profit, Phobias_profit.profit, Skills_profit.profit
								  FROM Profession_profit JOIN Health_profit ON Profession_profit.catastrophe_id=Health_profit.catastrophe_id  
								  JOIN  Character_profit ON Character_profit.catastrophe_id = Profession_profit.catastrophe_id
								  JOIN Baggage_profit ON Baggage_profit.catastrophe_id = Profession_profit.catastrophe_id
								  JOIN Biological_characteristics_profit ON Biological_characteristics_profit.catastrophe_id = Profession_profit.catastrophe_id
								  JOIN Hobby_profit ON Hobby_profit.catastrophe_id = Profession_profit.catastrophe_id
								  JOIN Phobias_profit ON Phobias_profit.catastrophe_id = Phobias_profit.catastrophe_id
								  JOIN Skills_profit ON Skills_profit.catastrophe_id = Skills_profit.catastrophe_id
								  WHERE Profession_profit.catastrophe_id = $1 and Profession_profit.profession_id = $2 and
 								  Health_profit.health_id = $3 and Character_profit.character_id = $4 and
								  Baggage_profit.baggage_id = $5 and Biological_characteristics_profit.characteristics_id = $6 and
								  Hobby_profit.hobby_id = $7 and Phobias_profit.phobia_id = $8 and Skills_profit.skills_id = $9`,
		catastropheId, p.professionId, p.healthId, p.characterId, p.baggageId, p.biologicalCharacteristicId,
		p.hobbyId, p.phobiasId, p.skillId)
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
