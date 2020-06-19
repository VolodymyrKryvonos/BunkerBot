package Chat

import "../Game"

type Chat struct {
	game              *Game.Game
	lang              string
	registrationMsgId int
	langMsgId         int
}

const (
	EN = "en"
	RU = "ru"
)

func (chat *Chat) SetGame(game *Game.Game) {
	chat.game = game
}

func (chat Chat) GetGame() *Game.Game {
	return chat.game
}

func (chat *Chat) SetLangMsgId(msg int) {
	chat.langMsgId = msg
}

func (chat Chat) GetLangMsgId() int {
	return chat.langMsgId
}

func (chat *Chat) SetRegistrationMsgId(msg int) {
	chat.registrationMsgId = msg
}

func (chat Chat) GetRegistrationMsgId() int {
	return chat.registrationMsgId
}

func (chat *Chat) SetLang(n string) {
	chat.lang = n
}

func (chat *Chat) GetLang() string {
	return chat.lang
}
