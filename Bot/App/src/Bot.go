package main

import (
	"./Chat"
	"./Game"
	"./Game/DB"
	"./Text"
	"./loger"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var bot, err = tgbotapi.NewBotAPI("1168689726:AAHvx5_NlWlRKQ-jJ6bB8GaVl7P480u1mZc")
var chats map[int64]*Chat.Chat

const (
	BOT_USER_NAME = "@game_bunker_bot"
)

type jsonVote struct {
	ChatId int64
	Vote   int
}

func main() {
	chats = make(map[int64]*Chat.Chat)
	onRestart()
	loger.LogErr(err)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	loger.LogErr(err)

	for update := range updates {
		if update.Message == nil {
			if update.CallbackQuery != nil {
				switch {
				case update.CallbackQuery.Message.Chat.IsPrivate():
					go callbackDialogQuery(update)
				default:
					go callbackChatQuery(update)
				}

			}
			continue
		}
		if update.Message.Chat.IsSuperGroup() || update.Message.Chat.IsGroup() {
			chat := chats[update.Message.Chat.ID]
			if chat == nil {
				go groupChat(update)
				continue
			}
			if chat.GetGame().GetGameStage() > Game.REGISTRATION {
				if deleteNotPlayerMsg(update, chat) {
					continue
				}

			}
			groupChat(update)
		} else {
			dialog(update)
		}
		if update.Message.Command() != "" {
			_, err := bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.MessageID,
			})
			loger.LogErr(err)
		}
	}
}

func deleteNotPlayerMsg(update tgbotapi.Update, chat *Chat.Chat) bool {
	players := chat.GetGame().GetPlayers()
	if isFromAdministrator(update) {
		return false
	}
	for i := 0; i < chat.GetGame().GetNumberOfPlayers(); i++ {
		if players[i].GetUserId() == update.Message.From.ID {
			return false
		}
	}
	_, err = bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))
	if err != nil {
		loger.LogErr(err)
	}
	return true
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
		chat := &Chat.Chat{}
		chat.SetLang(lang)
		game := &Game.Game{}
		game.SetGameStage(Game.GAME_IS_OVER)
		chat.SetGame(game)
		chats[chatID] = chat
	}
}

func callbackDialogQuery(update tgbotapi.Update) {
	var vote jsonVote
	err := json.Unmarshal([]byte(update.CallbackQuery.Data), &vote)
	if err != nil {
		loger.LogErr(err)
		return
	}
	chat := chats[vote.ChatId]
	alertWindow := tgbotapi.CallbackConfig{}
	alertWindow.ShowAlert = true
	alertWindow.CacheTime = 1
	alertWindow.CallbackQueryID = update.CallbackQuery.ID
	switch chat.GetGame().GetGameStage() {
	case Game.SELECTING_CHARACTERISTICS:
		var player *Game.Player
		player = chat.GetGame().FindByID(update.CallbackQuery.From.ID)
		switch chat.GetLang() {
		case Chat.RU:
			alertWindow.Text = Text.CHARACTERISTIC_ALREADY_OPENED_RU
		case Chat.EN:
			alertWindow.Text = Text.CHARACTERISTIC_ALREADY_OPENED_EN
		}

		switch vote.Vote {
		case 0:
			if player.IsHealthOpen() {
				_, err := bot.AnswerCallbackQuery(alertWindow)
				if err != nil {
					loger.LogErr(err)

				}

				return
			}
			player.OpenHealth()

		case 1:
			if player.IsCharOpen() {
				_, err := bot.AnswerCallbackQuery(alertWindow)
				if err != nil {
					loger.LogErr(err)
				}
				return
			}
			player.OpenChar()

		case 2:
			if player.IsBagOpen() {
				_, err := bot.AnswerCallbackQuery(alertWindow)
				if err != nil {
					loger.LogErr(err)
				}
				return
			}
			player.OpenBag()

		case 3:
			if player.IsBioOpen() {
				_, err := bot.AnswerCallbackQuery(alertWindow)
				if err != nil {
					loger.LogErr(err)
				}
				return
			}
			player.OpenBio()

		case 4:
			if player.IsHobbyOpen() {
				_, err := bot.AnswerCallbackQuery(alertWindow)
				if err != nil {
					loger.LogErr(err)
				}
				return
			}
			player.OpenHobby()

		case 5:
			if player.IsPhobiaOpen() {
				_, err := bot.AnswerCallbackQuery(alertWindow)
				if err != nil {
					loger.LogErr(err)
				}
				return
			}
			player.OpenPhobia()

		case 6:
			if player.IsSkillOpen() {
				_, err := bot.AnswerCallbackQuery(alertWindow)
				if err != nil {
					loger.LogErr(err)
				}
				return
			}
			player.OpenSkill()

		}
		delMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID)
		_, err = bot.DeleteMessage(delMsg)
		loger.LogErr(err)
	case Game.VOTING:
		alertWindow.CallbackQueryID = update.CallbackQuery.ID
		if chat.GetGame().PlayersToKick[vote.Vote].GetUserId() != update.CallbackQuery.From.ID {
			chat.GetGame().PlayersToKick[vote.Vote].IncrementAgainstVotes()
			msg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
				chat.GetGame().FindByID(update.CallbackQuery.From.ID).MsgId())
			_, err = bot.DeleteMessage(msg)
			loger.LogErr(err)
		} else {

			switch chat.GetLang() {
			case Chat.EN:
				alertWindow.Text = Text.VOTE_AGAINST_YOURSELF_EN
			case Chat.RU:
				alertWindow.Text = Text.VOTE_AGAINST_YOURSELF_RU
			}
			_, err = bot.AnswerCallbackQuery(alertWindow)
			loger.LogErr(err)
		}
	case Game.DISCUSSION:

	}

}

func callbackChatQuery(update tgbotapi.Update) {
	chat := chats[update.CallbackQuery.Message.Chat.ID]
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
	switch update.CallbackQuery.Data {
	case Chat.EN:
		setChatLocalization(update)
		return
	case Chat.RU:
		setChatLocalization(update)
		return
	case "join":
		if chat.GetGame().GetGameStage() == Game.REGISTRATION {
			player := Game.Player{}
			player.SetUser(update.CallbackQuery.From)

			players := chat.GetGame().GetPlayers()
			for i := 0; i < chat.GetGame().GetNumberOfPlayers(); i++ {
				if player.GetUserId() == players[i].GetUserId() {
					msgText := ""
					switch chat.GetLang() {
					case Chat.EN:
						msgText = Text.ALREADY_REGISTRED_EN
					case Chat.RU:
						msgText = Text.ALREADY_REGISTRED_RU
					}
					msg = tgbotapi.NewMessage(int64(update.CallbackQuery.From.ID), msgText)
					_, err = bot.Send(msg)
					if err != nil {
						loger.LogErr(err)
					}
					return
				}
			}
			chat.GetGame().AddPlayer(player)
			msgText := ""
			numberOfPlayers := ""
			keyboard := tgbotapi.InlineKeyboardMarkup{}
			row := make([]tgbotapi.InlineKeyboardButton, 2)
			switch chat.GetLang() {
			case Chat.RU:
				numberOfPlayers = Text.NUMBER_OF_PLAYERS_RU
				msgText = Text.REGISTRATION_RU + "\n\n" + "Зарегистрировались:\n"
				row[0] = tgbotapi.NewInlineKeyboardButtonData(Text.JOIN_RU, "join")
				row[1] = tgbotapi.NewInlineKeyboardButtonData(Text.LEAVE_RU, "leave")
			case Chat.EN:
				numberOfPlayers = Text.NUMBER_OF_PLAYERS_EN
				msgText = Text.REGISTRATION_EN + "\n\n" + "Registered:\n"
				row[0] = tgbotapi.NewInlineKeyboardButtonData(Text.JOIN_EN, "join")
				row[1] = tgbotapi.NewInlineKeyboardButtonData(Text.LEAVE_EN, "leave")
			}
			players = chat.GetGame().GetPlayers()
			msgText += players[0].GetUserName()
			for i := 1; i < chat.GetGame().GetNumberOfPlayers(); i++ {
				msgText += ", " + players[i].GetUserName()
			}
			msgText += "\n\n" + numberOfPlayers + fmt.Sprintf("%d", chat.GetGame().GetNumberOfPlayers())
			msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				msgText)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			msg.ReplyMarkup = &keyboard
			_, err = bot.Send(msg)
			loger.LogErr(err)
		}
	case "leave":
		if chat.GetGame().GetGameStage() == Game.REGISTRATION {
			chat.GetGame().RemovePlayer(update.CallbackQuery.From.ID)
			msgText := ""
			numberOfPlayers := ""
			keyboard := tgbotapi.InlineKeyboardMarkup{}
			row := make([]tgbotapi.InlineKeyboardButton, 2)
			switch chat.GetLang() {
			case Chat.RU:
				numberOfPlayers = Text.NUMBER_OF_PLAYERS_RU
				msgText = Text.REGISTRATION_RU + "\n\n" + "Зарегистрировались:\n"
				row[0] = tgbotapi.NewInlineKeyboardButtonData(Text.JOIN_RU, "join")
				row[1] = tgbotapi.NewInlineKeyboardButtonData(Text.LEAVE_RU, "leave")
			case Chat.EN:
				numberOfPlayers = Text.NUMBER_OF_PLAYERS_EN
				msgText = Text.REGISTRATION_EN + "\n\n" + "Registered:\n"
				row[0] = tgbotapi.NewInlineKeyboardButtonData(Text.JOIN_EN, "join")
				row[1] = tgbotapi.NewInlineKeyboardButtonData(Text.LEAVE_EN, "leave")
			}
			players := chat.GetGame().GetPlayers()
			if chat.GetGame().GetNumberOfPlayers() > 0 {
				msgText += players[0].GetUserName()
			}
			for i := 1; i < chat.GetGame().GetNumberOfPlayers(); i++ {
				msgText += ", " + players[i].GetUserName()
			}
			msgText += "\n\n" + numberOfPlayers + fmt.Sprintf("%d", chat.GetGame().GetNumberOfPlayers())
			msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				msgText)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			msg.ReplyMarkup = &keyboard
			_, err = bot.Send(msg)
			loger.LogErr(err)
		}
	case "done":
		if !isBotAdmin(update.CallbackQuery.Message.Chat.ID) {
			_, err = bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "Test",
				ShowAlert:       true,
				URL:             "",
				CacheTime:       3,
			})
			loger.LogErr(err)
		} else {
			_, err := bot.DeleteMessage(
				tgbotapi.NewDeleteMessage(
					update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
			if err != nil {
				loger.LogErr(err)
			}
		}
	case "language":
		if isFromAdministrator(update) {
			keyboard := tgbotapi.InlineKeyboardMarkup{}
			row := make([]tgbotapi.InlineKeyboardButton, 2)
			row[0] = tgbotapi.NewInlineKeyboardButtonData("EN", "en")
			row[1] = tgbotapi.NewInlineKeyboardButtonData("RU", "ru")
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			msg.ReplyMarkup = keyboard
			msg.Text = "Select language:"
		}
	}
	_, err = bot.Send(msg)
	loger.LogErr(err)
}

func isFromAdministrator(update tgbotapi.Update) bool {

	var chatId int64
	var userId int

	if update.Message != nil {
		chatId = update.Message.Chat.ID
		userId = update.Message.From.ID
	} else {
		chatId = update.CallbackQuery.Message.Chat.ID
		userId = update.CallbackQuery.From.ID
	}

	member, err := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
		ChatID:             chatId,
		SuperGroupUsername: "",
		UserID:             userId,
	})

	if err != nil {
		loger.LogErr(err)
	}

	return member.IsAdministrator() || member.IsCreator()
}

func setChatLocalization(update tgbotapi.Update) {
	if isFromAdministrator(update) {
		_, err = bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
			chats[update.CallbackQuery.Message.Chat.ID].GetLangMsgId()))
		if err != nil {
			loger.LogErr(err)
		}
		db := DB.GetDataBase()

		_, err := db.Exec("INSERT or REPLACE INTO Localization (Chat_id, lang) VALUES($1,$2)",
			update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)

		if err != nil {
			loger.LogErr(err)
		}

		chats[update.CallbackQuery.Message.Chat.ID].SetLang(update.CallbackQuery.Data)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

		switch update.CallbackQuery.Data {

		case Chat.EN:
			msg.Text = Text.CHANGED_LANG_EN

		case Chat.RU:
			msg.Text = Text.CHANGED_LANG_RU

		}

		sentMsg, err := bot.Send(msg)
		if err != nil {
			loger.LogErr(err)
		}
		chats[update.CallbackQuery.Message.Chat.ID].SetLangMsgId(sentMsg.MessageID)
		_, err = bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
		})
		loger.LogErr(err)
	}
}

func dialog(update tgbotapi.Update) {
	command := update.Message.Text
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "")
	switch command {
	case Text.NEW_GAME:
		msg.Text = "123"
	case Text.START:
	case Text.RULES:
		msg.Text = "321"
	}
	_, err = bot.Send(msg)
	if err != nil {
		loger.LogErr(err)
	}
}

func formPlayerInfoMsg(chat *Chat.Chat) string {
	var msgText string
	db := DB.GetDataBase()
	for i := 0; i < chat.GetGame().GetNumberOfPlayers(); i++ {
		query, err := db.Query("SELECT Profession_"+chat.GetLang()+" FROM Profession WHERE id=$1",
			chat.GetGame().GetPlayers()[i].GetProfId())
		if err != nil {
			loger.LogErr(err)
		}
		query.Next()
		var reader string

		err = query.Scan(&reader)
		loger.LogErr(err)
		msgText += chat.GetGame().GetPlayers()[i].GetFullName() + ": " + reader
		if chat.GetGame().GetPlayers()[i].IsCharOpen() {
			query, err := db.Query("SELECT character_"+chat.GetLang()+" FROM Character WHERE id=$1",
				chat.GetGame().GetPlayers()[i].GetCharacterId())
			loger.LogErr(err)
			query.Next()
			err = query.Scan(&reader)
			loger.LogErr(err)
			msgText += ", " + reader
		}
		if chat.GetGame().GetPlayers()[i].IsBioOpen() {
			msgText += ", " + chat.GetGame().GetPlayers()[i].GetBioChar(chat.GetLang())
		}
		if chat.GetGame().GetPlayers()[i].IsHealthOpen() {
			query, err := db.Query("SELECT health_"+chat.GetLang()+" FROM Health WHERE id=$1",
				chat.GetGame().GetPlayers()[i].GetHealthId())
			loger.LogErr(err)
			query.Next()
			err = query.Scan(&reader)
			loger.LogErr(err)
			msgText += ", " + reader
		}
		if chat.GetGame().GetPlayers()[i].IsPhobiaOpen() {
			query, err := db.Query("SELECT phobias_"+chat.GetLang()+" FROM Phobias WHERE id=$1",
				chat.GetGame().GetPlayers()[i].GetPhobiasId())
			loger.LogErr(err)
			query.Next()
			err = query.Scan(&reader)
			loger.LogErr(err)
			msgText += ", " + reader
		}
		if chat.GetGame().GetPlayers()[i].IsBagOpen() {
			query, err := db.Query("SELECT baggage_"+chat.GetLang()+" FROM Baggage WHERE id=$1",
				chat.GetGame().GetPlayers()[i].GetBaggageId())
			loger.LogErr(err)
			query.Next()
			err = query.Scan(&reader)
			loger.LogErr(err)
			msgText += ", " + reader
		}
		if chat.GetGame().GetPlayers()[i].IsHobbyOpen() {
			query, err := db.Query("SELECT hobby_"+chat.GetLang()+" FROM Hobby WHERE id=$1",
				chat.GetGame().GetPlayers()[i].GetHobbyId())
			loger.LogErr(err)
			query.Next()
			err = query.Scan(&reader)
			loger.LogErr(err)
			msgText += ", " + reader
		}
		if chat.GetGame().GetPlayers()[i].IsSkillOpen() {
			query, err := db.Query("SELECT skills_"+chat.GetLang()+" FROM Skills WHERE id=$1",
				chat.GetGame().GetPlayers()[i].GetSckillId())
			loger.LogErr(err)
			query.Next()
			err = query.Scan(&reader)
			loger.LogErr(err)
			msgText += ", " + reader
		}
		if chat.GetGame().GetPlayers()[i].IsAlive() {
			switch chat.GetLang() {
			case Chat.EN:
				msgText += ", " + "alive\n"
			case Chat.RU:
				msgText += ", " + "живой\n"
			}
		} else {
			switch chat.GetLang() {
			case Chat.EN:
				msgText += ", " + "dead\n"
			case Chat.RU:
				msgText += ", " + "мертвый\n"
			}
		}
	}
	return msgText
}

func groupChat(update tgbotapi.Update) {
	var chat *Chat.Chat
	command := update.Message.Command()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if chats[update.Message.Chat.ID] == nil {
		chat = &Chat.Chat{}
		chat.SetLang(Chat.EN)
		game := &Game.Game{}
		game.SetGameStage(Game.GAME_IS_OVER)
		chat.SetGame(game)
		chats[update.Message.Chat.ID] = chat
	} else {
		chat = chats[update.Message.Chat.ID]
	}
	if !isBotAdmin(update.Message.Chat.ID) {
		settings(update)
		return
	}
	switch command {
	case Text.NEW_GAME:
		if chat.GetGame().GetGameStage() == Game.GAME_IS_OVER {
			chat.GetGame().SetGameStage(Game.REGISTRATION)
			keyboard := tgbotapi.InlineKeyboardMarkup{}
			row := make([]tgbotapi.InlineKeyboardButton, 2)
			switch chat.GetLang() {
			case Chat.RU:
				msg.Text = Text.REGISTRATION_RU
				row[0] = tgbotapi.NewInlineKeyboardButtonData(Text.JOIN_RU, "join")
				row[1] = tgbotapi.NewInlineKeyboardButtonData(Text.LEAVE_RU, "leave")
			case Chat.EN:
				msg.Text = Text.REGISTRATION_EN
				row[0] = tgbotapi.NewInlineKeyboardButtonData(Text.JOIN_EN, "join")
				row[1] = tgbotapi.NewInlineKeyboardButtonData(Text.LEAVE_EN, "leave")
			}
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			msg.ReplyMarkup = keyboard
			m, err := bot.Send(msg)
			if err != nil {
				loger.LogErr(err)
			}
			chat.SetRegistrationMsgId(m.MessageID)
			_, err = bot.PinChatMessage(tgbotapi.PinChatMessageConfig{
				ChatID:              update.Message.Chat.ID,
				MessageID:           m.MessageID,
				DisableNotification: false,
			})
			loger.LogErr(err)
			go func() {
				timer := time.NewTimer(60 * time.Second)
				select {
				case <-timer.C:
					msg = startGame(update, chat)
					_, err := bot.Send(msg)
					if err != nil {
						loger.LogErr(err)
					}
				}
			}()

			return
		}

	case Text.STOP:
		_, err = bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.Message.Chat.ID,
			chat.GetRegistrationMsgId()))
		chat.GetGame().FinishGame()
		loger.LogErr(err)
	case Text.START:
		if isFromAdministrator(update) {
			msg = startGame(update, chat)
		}

	case Text.RULES:
		msg = tgbotapi.NewMessage(int64(update.Message.From.ID), "")
		msg.Text = Text.TEST
	case Text.LANG:
		if isFromAdministrator(update) {
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

func formCharacteristicMsg(chatId int64, chat *Chat.Chat) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, "")
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	buttons := make([]tgbotapi.InlineKeyboardButton, 10)
	switch chat.GetLang() {
	case Chat.EN:
		msg.Text = Text.SELECT_CHARACTERISTIC_EN
	case Chat.RU:
		msg.Text = Text.SELECT_CHARACTERISTIC_RU
	}
	var vote jsonVote
	vote.ChatId = msg.ChatID

	for i := 0; i < 7; i++ {
		vote.Vote = i
		voteParsed, err := json.Marshal(vote)
		if err != nil {
			loger.LogErr(err)
			return tgbotapi.MessageConfig{}
		}
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData("", string(voteParsed))
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, buttons[i:i+1])
	}
	switch chat.GetLang() {
	case Chat.RU:
		buttons[0].Text = Text.HEALTH_RU
		buttons[1].Text = Text.CHARACTER_RU
		buttons[2].Text = Text.BAGGAGE_RU
		buttons[3].Text = Text.BIO_RU
		buttons[4].Text = Text.HOBBY_RU
		buttons[5].Text = Text.PHOBIA_RU
		buttons[6].Text = Text.SKILL_RU
	case Chat.EN:
		buttons[0].Text = Text.HEALTH_EN
		buttons[1].Text = Text.CHARACTER_EN
		buttons[2].Text = Text.BAGGAGE_EN
		buttons[3].Text = Text.BIO_EN
		buttons[4].Text = Text.HOBBY_EN
		buttons[5].Text = Text.PHOBIA_EN
		buttons[6].Text = Text.SKILL_EN
	}
	msg.ReplyMarkup = keyboard
	return msg
}

func formVotingMsg(chatId int64, chat *Chat.Chat) tgbotapi.MessageConfig {

	if chat.GetGame().PlayersToKick == nil {

		chat.GetGame().PlayersToKick = make([]*Game.Player, chat.GetGame().NumberOfAlivePlayers())
		index := 0
		for i := 0; i < chat.GetGame().GetNumberOfPlayers(); i++ {
			if chat.GetGame().GetPlayers()[i].IsAlive() {
				chat.GetGame().PlayersToKick[index] = &chat.GetGame().GetPlayers()[i]
				index++
			}
		}
	}

	msg := tgbotapi.NewMessage(chatId, "")
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	buttons := make([]tgbotapi.InlineKeyboardButton, len(chat.GetGame().PlayersToKick))
	for i, val := range chat.GetGame().PlayersToKick {

		var vote jsonVote
		vote.ChatId = msg.ChatID
		vote.Vote = i
		voteParsed, err := json.Marshal(vote)
		if err != nil {
			loger.LogErr(err)
			return tgbotapi.MessageConfig{}
		}
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i+1)+". "+
			val.GetFullName(), string(voteParsed))
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, buttons[i:i+1])

	}
	msg.ReplyMarkup = keyboard
	switch chat.GetLang() {
	case Chat.RU:
		msg.Text = Text.VOTE_AGAINST_RU
	case Chat.EN:
		msg.Text = Text.VOTE_AGAINST_EN
	}
	return msg
}

func startGame(update tgbotapi.Update, chat *Chat.Chat) tgbotapi.MessageConfig {
	if chat.GetGame().GetGameStage() != Game.REGISTRATION {
		return tgbotapi.MessageConfig{}
	}

	_, err = bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.Message.Chat.ID,
		chat.GetRegistrationMsgId()))
	loger.LogErr(err)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	err := chat.GetGame().NewGame()
	if err != nil {
		msg.Text = err.Error()
		loger.LogErr(err)
		return msg
	}
	players := chat.GetGame().GetPlayers()
	for i := 0; i < chat.GetGame().GetNumberOfPlayers(); i++ {
		sendProfile(players[i], chat.GetLang())
	}
	db := DB.GetDataBase()
	query, err := db.Query(
		strings.ReplaceAll(
			"SELECT catastrophe_%s, description_%s,destruction FROM Catastrophe WHERE id = $1",
			"%s", chat.GetLang()),
		chat.GetGame().GetCatastropheId())
	if err != nil {
		loger.LogErr(err)
		chat.GetGame().FinishGame()
		return tgbotapi.MessageConfig{}
	}

	if query.Next() {
		var (
			amenities    = ""
			capacity     = 0
			area         = 0
			timeInBunker = 0
			catastrophe  = ""
			description  = ""
			destruction  = 0
			alive        = 0
		)
		rand.Seed(time.Now().UnixNano())
		capacity = chat.GetGame().GetNumberOfPlayers() / 2
		timeInBunker = rand.Int()%17 + 3
		area = (capacity + rand.Int()%12) * 10
		err = query.Scan(&catastrophe, &description, &destruction)
		loger.LogErr(err)
		alive = (rand.Int()%(15-destruction) + 1) * 5
		destruction = (rand.Int()%(10+destruction) + 1) * 5

		switch chat.GetLang() {
		case Chat.RU:
			query, err = db.Query("SELECT amenities_ru FROM Bunker ORDER BY RANDOM() limit 1")
			loger.LogErr(err)
			query.Next()
			err = query.Scan(&amenities)
			loger.LogErr(err)
			msg.Text = fmt.Sprintf(Text.DESCRIPTION_RU, catastrophe, description, destruction, alive)
			msg.Text += fmt.Sprintf(Text.BUNKER_RU, capacity, area, timeInBunker, amenities)
		case Chat.EN:
			query, err = db.Query("SELECT amenities_en FROM Bunker ORDER BY RANDOM() limit 1")
			loger.LogErr(err)
			query.Next()
			err=query.Scan(&amenities)
			loger.LogErr(err)
			msg.Text = fmt.Sprintf(Text.DESCRIPTION_EN, catastrophe, description, destruction, alive)
			msg.Text += fmt.Sprintf(Text.BUNKER_EN, capacity, area, timeInBunker, amenities)
		}
		loger.LogErr(err)
		msg.ParseMode = "markdown"
	}

	_,err=bot.Send(msg)
	loger.LogErr(err)
	go GameLogic(msg.ChatID, chat)
	msg.ReplyMarkup = nil
	msg.Text = ""

	return msg
}

func GameLogic(chatId int64, chat *Chat.Chat) {
	closeCharacteristics := Game.NUMBER_OF_CHARACTERISTICS - 1
	game := chat.GetGame()
	for {
		if game.BunkerCap() < game.NumberOfAlivePlayers() {
			if closeCharacteristics-(game.NumberOfAlivePlayers()-game.BunkerCap()) > 0 {
				game.SetGameStage(Game.SELECTING_CHARACTERISTICS)
			} else {
				game.SetGameStage(Game.DISCUSSION)
			}
			switch game.GetGameStage() {
			case Game.SELECTING_CHARACTERISTICS:
				selectingCharacteristics(chat, chatId)
				closeCharacteristics--
			case Game.DISCUSSION:
				discussion(chat, chatId)
			}
		} else {
			msg := tgbotapi.NewMessage(chatId, "")
			text := ""
			for i := 0; i < chat.GetGame().GetNumberOfPlayers(); i++ {
				if chat.GetGame().GetPlayers()[i].IsAlive() {
					text += ", " + chat.GetGame().GetPlayers()[i].GetFullName()
				}
			}
			if len(text) > 1 {
				msg.Text = text[1:]
			}
			switch chat.GetLang() {
			case Chat.RU:
				msg.Text += ""
			case Chat.EN:
				msg.Text += ""
			}
			_, err = bot.Send(msg)
			loger.LogErr(err)
			game.FinishGame()
			break
		}
	}
}

func discussion(chat *Chat.Chat, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "")
	msg.Text = formPlayerInfoMsg(chat)
	_, err := bot.Send(msg)
	if err != nil {
		loger.LogErr(err)
	}

	timer := time.NewTimer(10 * time.Second)
	select {
	case <-timer.C:
		chat.GetGame().SetGameStage(Game.VOTING)
		voting(chat, chatId)
	}
}

func timeForVoting(chat *Chat.Chat, chatId int64) {
	msg := formVotingMsg(chatId, chat)
	for i := 0; i < chat.GetGame().GetNumberOfPlayers(); i++ {
		if chat.GetGame().GetPlayers()[i].IsAlive() {
			msg.ChatID = int64(chat.GetGame().GetPlayers()[i].GetUser().ID)
			msgReply, err := bot.Send(msg)
			loger.LogErr(err)
			chat.GetGame().GetPlayers()[i].SetMsgId(msgReply.MessageID)
		}
	}
	timer := time.NewTimer(10 * time.Second)
	select {
	case <-timer.C:
		for i := 0; i < chat.GetGame().GetNumberOfPlayers(); i++ {
			if chat.GetGame().GetPlayers()[i].IsAlive() {
				_, err := bot.DeleteMessage(tgbotapi.NewDeleteMessage(
					int64(chat.GetGame().GetPlayers()[i].GetUser().ID),
					chat.GetGame().GetPlayers()[i].MsgId()))
				loger.LogErr(err)
			}
		}
	}
}

func voting(chat *Chat.Chat, chatId int64) {
	timeForVoting(chat, chatId)
	msg := tgbotapi.NewMessage(chatId, "")
	chat.GetGame().Kick()
	if len(chat.GetGame().PlayersToKick) == 1 {
		msg.ChatID = chatId
		chat.GetGame().PlayersToKick[0].Kill()
		chat.GetGame().SetNumberOfAlivePlayers(chat.GetGame().NumberOfAlivePlayers() - 1)
		switch chat.GetLang() {
		case Chat.EN:
			msg.Text = fmt.Sprintf(Text.KICK_EN, chat.GetGame().PlayersToKick[0].GetFullName())
		case Chat.RU:
			msg.Text = fmt.Sprintf(Text.KICK_RU, chat.GetGame().PlayersToKick[0].GetFullName())
		}
		msg.ReplyMarkup = nil
		_, err := bot.Send(msg)
		loger.LogErr(err)
		chat.GetGame().PlayersToKick = nil
	} else {
		msg.Text = ""
		msg.ReplyMarkup = nil
		msg.ChatID = chatId
		for _, val := range chat.GetGame().PlayersToKick {
			msg.Text += val.GetFullName() + ", "
		}
		switch chat.GetLang() {
		case Chat.RU:
			msg.Text += Text.MORE_THAN_ONE_RU
		case Chat.EN:
			msg.Text += Text.MORE_THAN_ONE_EN
		}
		_, err := bot.Send(msg)
		loger.LogErr(err)
		timer := time.NewTimer(10 * time.Second)
		select {
		case <-timer.C:
			timeForVoting(chat, chatId)
			chat.GetGame().Kick()
			for _, val := range chat.GetGame().PlayersToKick {
				val.Kill()
				chat.GetGame().SetNumberOfAlivePlayers(chat.GetGame().NumberOfAlivePlayers() - 1)
				switch chat.GetLang() {
				case Chat.EN:
					msg.Text = fmt.Sprintf(Text.KICK_EN, val.GetFullName())
				case Chat.RU:
					msg.Text = fmt.Sprintf(Text.KICK_RU, val.GetFullName())
				}
				msg.ReplyMarkup = nil
				_, err := bot.Send(msg)
				loger.LogErr(err)
			}
			chat.GetGame().PlayersToKick = nil
		}
	}
}

func selectingCharacteristics(chat *Chat.Chat, chatId int64) {
	msg := formCharacteristicMsg(chatId, chat)
	players := chat.GetGame().GetPlayers()
	for i := 0; i < chat.GetGame().GetNumberOfPlayers(); i++ {
		if players[i].IsAlive() {
			msg.ChatID = int64(players[i].GetUser().ID)
			msgReply, err := bot.Send(msg)
			loger.LogErr(err)
			players[i].SetMsgId(msgReply.MessageID)
		}
	}
	timer := time.NewTimer(5 * time.Second)
	select {
	case <-timer.C:
		for i := 0; i < chat.GetGame().GetNumberOfPlayers(); i++ {
			if players[i].IsAlive() {
				_, err := bot.DeleteMessage(tgbotapi.NewDeleteMessage(int64(players[i].GetUser().ID),
					players[i].MsgId()))
				if err == nil {
					switch {
					case !players[i].IsHealthOpen():
						players[i].OpenHealth()
					case !players[i].IsPhobiaOpen():
						players[i].OpenPhobia()
					case !players[i].IsBioOpen():
						players[i].OpenBio()
					case !players[i].IsBagOpen():
						players[i].OpenBag()
					case !players[i].IsHobbyOpen():
						players[i].OpenHobby()
					case !players[i].IsSkillOpen():
						players[i].OpenSkill()
					case !players[i].IsCharOpen():
						players[i].OpenChar()
					}
				}
			}
		}
	}
}

func sendProfile(player Game.Player, lang string) {
	userId := player.GetUserId()
	var profile string
	db := DB.GetDataBase()
	s := strings.ReplaceAll(`SELECT Profession.profession_%s, Health.health_%s, Character.character_%s, 
								Baggage.baggage_%s,Hobby.hobby_%s, Phobias.phobias_%s, Skills.skills_%s
								FROM Profession, Health,Character,Baggage,Hobby,Phobias,Skills
								WHERE Profession.id = $1 and Health.id = $2 and
								Character.id = $3 and Baggage.id = $4 and Hobby.id = $5 and
								Phobias.id = $6 and Skills.id = $7`, "%s", lang)

	query, err := db.Query(s, player.GetProfId(), player.GetHealthId(), player.GetCharacterId(),
		player.GetBaggageId(), player.GetHobbyId(), player.GetPhobiasId(),
		player.GetSckillId())

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
		err := query.Scan(&profession, &health, &character, &baggage, &hobby, &phobias, &skills)
		biologicalCharacteristics = player.GetBioChar(lang)
		if err != nil {
			loger.LogErr(err)
		}
	}

	switch lang {

	case Chat.EN:
		profile = fmt.Sprintf(Text.PROFILE_EN, profession, health, character, baggage,
			biologicalCharacteristics, hobby, phobias, skills)
	case Chat.RU:
		profile = fmt.Sprintf(Text.PROFILE_RU, profession, health, character, baggage,
			biologicalCharacteristics, hobby, phobias, skills)
	}

	msg := tgbotapi.NewMessage(int64(userId), profile)
	_, err = bot.Send(msg)
	if err != nil {
		loger.LogErr(err)
	}
}

func settings(update tgbotapi.Update) tgbotapi.Message {
	botMember, err := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
		ChatID:             update.Message.Chat.ID,
		SuperGroupUsername: "",
		UserID:             bot.Self.ID,
	})

	if err != nil {
		loger.LogErr(err)
		return tgbotapi.Message{}
	}

	if !botMember.IsAdministrator() {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		keyboard := tgbotapi.InlineKeyboardMarkup{}
		row1 := make([]tgbotapi.InlineKeyboardButton, 1)
		row2 := make([]tgbotapi.InlineKeyboardButton, 1)

		row1[0] = tgbotapi.NewInlineKeyboardButtonData(Text.DONE, "done")
		row2[0] = tgbotapi.NewInlineKeyboardButtonData(Text.LANGUAGE, "language")

		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row1)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row2)
		msg.Text = Text.SETTINGS
		msg.ReplyMarkup = keyboard
		settingsMsg, err := bot.Send(msg)

		if err != nil {
			loger.LogErr(err)
		}
		return settingsMsg
	}
	return tgbotapi.Message{}
}

func isBotAdmin(chatId int64) bool {
	botMember, err := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
		ChatID:             chatId,
		SuperGroupUsername: "",
		UserID:             bot.Self.ID,
	})
	if err != nil {
		loger.LogErr(err)
	}
	return botMember.IsAdministrator()
}
