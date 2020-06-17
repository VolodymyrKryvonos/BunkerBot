package main

import (
	"./Game"
	"./Game/DB"
	"./Text"
	"./loger"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"time"
)

var bot, err = tgbotapi.NewBotAPI("1168689726:AAHvx5_NlWlRKQ-jJ6bB8GaVl7P480u1mZc")
var games map[int64]*Game.Game

const botUserName = "@game_bunker_bot"

func main() {
	games = make(map[int64]*Game.Game)
	onRestart()
	loger.LogErr(err)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	loger.LogErr(err)

	for update := range updates {
		if update.Message == nil {
			if update.CallbackQuery != nil {
				callbackQuery(update)
			}
			continue
		}
		if update.Message.Chat.IsSuperGroup() || update.Message.Chat.IsGroup() {
			groupChat(update)
		} else {
			dialog(update)
		}
	}
}

func onRestart() {
	db := DB.GetDataBase()
	query, err := db.Query("SELECT * FROM Localization")
	if err != nil {
		loger.LogErr(err)
		return
	}

	for query.Next() {
		var chatID int64
		var lang string
		err := query.Scan(&chatID, &lang)
		if err != nil {
			loger.LogErr(err)
			continue
		}
		game := &Game.Game{}
		game.SetLang(lang)
		game.SetGameStage(-1)
		games[chatID] = game
	}
}

func callbackQuery(update tgbotapi.Update) {
	setChatLocalization(update)
}

func isFromAdministrator(msg *tgbotapi.Message) bool {

	administrators, err := bot.GetChatAdministrators(tgbotapi.ChatConfig{
		ChatID:             msg.Chat.ID,
		SuperGroupUsername: "",
	})

	if err != nil {
		loger.LogErr(err)
	}

	for i := range administrators {
		if administrators[i].User.ID == msg.From.ID {
			return true
		}
	}

	return false
}

func setChatLocalization(update tgbotapi.Update) {
	if isFromAdministrator(update.CallbackQuery.Message) {
		db := DB.GetDataBase()

		t0 := time.Now()
		_, err := db.Exec("INSERT or replace INTO Localization (Chat_id, lang) VALUES($1,$2)",
			update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)

		fmt.Println(t0.Sub(time.Now()))
		if err != nil {
			loger.LogErr(err)
		}

		games[update.CallbackQuery.Message.Chat.ID].SetLang(update.CallbackQuery.Data)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

		switch update.CallbackQuery.Data {

		case Game.EN:
			msg.Text = Text.CHANGED_LANG_EN

		case Game.RU:
			msg.Text = Text.CHANGED_LANG_RU

		}

		_, err = bot.Send(msg)
		if err != nil {
			loger.LogErr(err)
		}
	}
}

func dialog(update tgbotapi.Update) {
	command := update.Message.Text
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "")
	switch command {
	case Text.NEW_GAME:
		msg.Text = "123"
	case Text.JOIN:

	case Text.LEAVE:
	case Text.START:
	case Text.RULES:
		msg.Text = "321"
	}
	_, err = bot.Send(msg)
	if err != nil {
		loger.LogErr(err)
	}
}

func groupChat(update tgbotapi.Update) {
	var game *Game.Game
	command := update.Message.Text
	command = strings.ReplaceAll(command, botUserName, "")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if games[update.Message.Chat.ID] == nil {
		game = &Game.Game{}
		game.SetLang(Game.EN)
		games[update.Message.Chat.ID] = game
		games[update.Message.Chat.ID].SetGameStage(-1)
	} else {
		game = games[update.Message.Chat.ID]
	}

	switch command {
	case Text.NEW_GAME:
		if game.GetGameStage() < 0 {
			game.SetGameStage(0)
			keyboard := tgbotapi.InlineKeyboardMarkup{}
			row := make([]tgbotapi.InlineKeyboardButton, 2)
			switch game.GetLang() {
			case Game.RU:
				msg.Text = Text.REGISTRATION_RU
				row[0] = tgbotapi.NewInlineKeyboardButtonData(Text.JOIN_RU, "join")
				row[1] = tgbotapi.NewInlineKeyboardButtonData(Text.LEAVE_RU, "leave")
			case Game.EN:
				msg.Text = Text.REGISTRATION_EN
				row[0] = tgbotapi.NewInlineKeyboardButtonData(Text.JOIN_EN, "join")
				row[1] = tgbotapi.NewInlineKeyboardButtonData(Text.LEAVE_EN, "leave")
			}
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			msg.ReplyMarkup = keyboard
		} else {
			switch game.GetLang() {
			case Game.RU:
				msg.Text = Text.GAME_ALREADY_STARTED_RU
			case Game.EN:
				msg.Text = Text.GAME_ALREADY_STARTED_EN
			}
		}
	case Text.JOIN:
		if game.GetGameStage() == 0 {
			var player Game.Player
			player.SetUserId(update.Message.From.ID)
			userName := update.Message.From.UserName
			if userName == "" {
				userName = update.Message.From.FirstName + " " + update.Message.From.LastName
			}
			player.SetUserName(userName)
			players := game.GetPlayers()
			for i := 0; i < len(players); i++ {
				if player.GetUserId() == players[i].GetUserId() {
					msg = tgbotapi.NewMessage(int64(update.Message.From.ID), userName+" already in the game")
					_, err = bot.Send(msg)
					if err != nil {
						loger.LogErr(err)
					}
					return
				}
			}
			game.AddPlayer(player)
			msg.Text = player.GetUserName() + " registered to the game"
		}
	case Text.LEAVE:
		if game.GetGameStage() == 0 {

		}
	case Text.START:
		if game.GetGameStage() == 0 {
			err := game.NewGame()
			if err != nil {
				msg.Text = err.Error()
				_, err = bot.Send(msg)
				if err != nil {
					loger.LogErr(err)
				}
				return
			}
			players := game.GetPlayers()
			for i := 0; i < game.GetNumberOfPlayers(); i++ {
				sendProfile(players[i], game.GetLang())
			}
		}
	case Text.RULES:
		msg = tgbotapi.NewMessage(int64(update.Message.From.ID), "")
		msg.Text = Text.TEST
	case Text.LANG:
		if isFromAdministrator(update.Message) {
			keyboard := tgbotapi.InlineKeyboardMarkup{}
			row := make([]tgbotapi.InlineKeyboardButton, 2)
			row[0] = tgbotapi.NewInlineKeyboardButtonData("EN", "en")
			row[1] = tgbotapi.NewInlineKeyboardButtonData("RU", "ru")
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			msg.ReplyMarkup = keyboard
			msg.Text = "Select language:"
		} else {
			_, err := bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.MessageID,
			})
			if err != nil {
				loger.LogErr(err)
			}
		}
	}
	_, err = bot.Send(msg)
	if err != nil {
		loger.LogErr(err)
	}
}

func sendProfile(player Game.Player, lang string) {
	userId := player.GetUserId()
	var profile string
	db := DB.GetDataBase()
	s := strings.ReplaceAll(`SELECT Profession_%s.profession_name, health_%s.health, character_%s.character, baggage_%s.baggage,
							   biological_characteristics_%s.characteristics, hobby_%s.hobby, phobias_%s.phobias, skills_%s.skills
							   FROM Profession_%s, health_%s,character_%s,baggage_%s,biological_characteristics_%s,hobby_%s,phobias_%s,skills_%s
							   WHERE Profession_%s.id = $1 and health_%s.id = $2 and
 							   character_%s.id = $3 and baggage_%s.id = $4 and
							   biological_characteristics_%s.id = $5 and hobby_%s.id = $6 and
							   phobias_%s.id = $7 and skills_%s.id = $8`, "%s", lang)

	query, err := db.Query(s, 1, 1, 1, 1, 1, 1, 1, 1)
	if err != nil {
		loger.LogErr(err)
		return
	}

	var (
		profession                string
		health                    string
		character                 string
		baggage                   string
		biologicalCharacteristics string
		hobby                     string
		phobias                   string
		skills                    string
	)

	if query.Next() {
		err := query.Scan(&profession, &health, &character, &baggage, &biologicalCharacteristics, &hobby, &phobias, &skills)
		if err != nil {
			loger.LogErr(err)
		}
	}

	switch lang {

	case Game.EN:
		profile = fmt.Sprintf(Text.PROFILE_EN, profession, health, character, baggage, biologicalCharacteristics, hobby, phobias, skills)
	case Game.RU:
		profile = fmt.Sprintf(Text.PROFILE_RU, profession, health, character, baggage, biologicalCharacteristics, hobby, phobias, skills)
	}

	msg := tgbotapi.NewMessage(int64(userId), profile)
	_, err = bot.Send(msg)
	if err != nil {
		loger.LogErr(err)
	}
}
