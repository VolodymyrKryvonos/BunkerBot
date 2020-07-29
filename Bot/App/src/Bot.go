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
			bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.MessageID,
			})
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
	switch chat.GetGame().GetGameStage() {
	case Game.SELECTING_CHARACTERISTICS:
		switch vote.Vote {
		case 0:

		case 1:
		case 2:
		case 3:
		case 4:
		case 5:
		case 6:
		}
	case Game.VOTING:
		if chat.GetGame().GetPlayers()[vote.Vote].GetUserId() != update.CallbackQuery.From.ID {
			chat.GetGame().GetPlayers()[vote.Vote].IncrementAgainstVotes()
			msg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
				chat.GetGame().FindByID(update.CallbackQuery.From.ID).MsgId())
			bot.DeleteMessage(msg)
		} else {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
			switch chat.GetLang() {
			case Chat.EN:
				msg.Text = Text.VOTE_AGAINST_YOURSELF_EN
			case Chat.RU:
				msg.Text = Text.VOTE_AGAINST_YOURSELF_RU
			}
			bot.Send(msg)
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
			bot.Send(msg)
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
			bot.Send(msg)
		}
	case "done":
		if !isBotAdmin(update.CallbackQuery.Message.Chat.ID) {
			bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "Test",
				ShowAlert:       true,
				URL:             "",
				CacheTime:       3,
			})
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
	bot.Send(msg)

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

		_, err := db.Exec("INSERT or replace INTO Localization (Chat_id, lang) VALUES($1,$2)",
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

		sendedMsg, err := bot.Send(msg)
		if err != nil {
			loger.LogErr(err)
		}
		chats[update.CallbackQuery.Message.Chat.ID].SetLangMsgId(sendedMsg.MessageID)
		bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
		})
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
		query, err := db.Query("SELECT profession_name FROM Profession_"+chat.GetLang()+" WHERE id=$1", chat.GetGame().GetPlayers()[i].GetProfId())
		loger.LogErr(err)
		query.Next()
		var reader string

		query.Scan(&reader)
		msgText += chat.GetGame().GetPlayers()[i].GetFullName() + ": " + reader
		if chat.GetGame().GetPlayers()[i].IsCharOpen() {
			query, err := db.Query("SELECT character FROM character_"+chat.GetLang()+" WHERE id=$1", chat.GetGame().GetPlayers()[i].GetCharacterId())
			loger.LogErr(err)
			query.Next()
			query.Scan(&reader)
			msgText += ", " + reader
		}
		if chat.GetGame().GetPlayers()[i].IsBioOpen() {
			msgText += ", " + chat.GetGame().GetPlayers()[i].GetBioChar(chat.GetLang())
		}
		if chat.GetGame().GetPlayers()[i].IsHealthOpen() {
			query, err := db.Query("SELECT health FROM health_"+chat.GetLang()+" WHERE id=$1", chat.GetGame().GetPlayers()[i].GetHealthId())
			loger.LogErr(err)
			query.Next()
			query.Scan(&reader)
			msgText += ", " + reader
		}
		if chat.GetGame().GetPlayers()[i].IsPhobiaOpen() {
			query, err := db.Query("SELECT phobias FROM phobias_"+chat.GetLang()+" WHERE id=$1", chat.GetGame().GetPlayers()[i].GetPhobiasId())
			loger.LogErr(err)
			query.Next()
			query.Scan(&reader)
			msgText += ", " + reader
		}
		if chat.GetGame().GetPlayers()[i].IsBagOpen() {
			query, err := db.Query("SELECT baggage FROM baggage_"+chat.GetLang()+" WHERE id=$1", chat.GetGame().GetPlayers()[i].GetBaggageId())
			loger.LogErr(err)
			query.Next()
			query.Scan(&reader)
			msgText += ", " + reader
		}
		if chat.GetGame().GetPlayers()[i].IsHobbyOpen() {
			query, err := db.Query("SELECT hobby FROM hobby_"+chat.GetLang()+" WHERE id=$1", chat.GetGame().GetPlayers()[i].GetHobbyId())
			loger.LogErr(err)
			query.Next()
			query.Scan(&reader)
			msgText += ", " + reader
		}
		if chat.GetGame().GetPlayers()[i].IsSkillOpen() {
			query, err := db.Query("SELECT skills FROM skills_"+chat.GetLang()+" WHERE id=$1", chat.GetGame().GetPlayers()[i].GetSckillId())
			loger.LogErr(err)
			query.Next()
			query.Scan(&reader)
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
	command := update.Message.Text
	command = strings.ReplaceAll(command, BOT_USER_NAME, "")
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
			bot.PinChatMessage(tgbotapi.PinChatMessageConfig{
				ChatID:              update.Message.Chat.ID,
				MessageID:           m.MessageID,
				DisableNotification: false,
			})

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
		bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, chat.GetRegistrationMsgId()))
		chat.GetGame().FinishGame()

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

func nextStage() {

}

func formCharacteristicMsg(msg tgbotapi.MessageConfig, chat *Chat.Chat) tgbotapi.MessageConfig {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	buttons := make([]tgbotapi.InlineKeyboardButton, 10)
	msg.Text = "Test"
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

func formVotingMsg(msg tgbotapi.MessageConfig, chat *Chat.Chat) tgbotapi.MessageConfig {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	buttons := make([]tgbotapi.InlineKeyboardButton, chat.GetGame().NumberOfAlivePlayers())
	for i := 0; i < chat.GetGame().NumberOfAlivePlayers(); i++ {
		var vote jsonVote
		vote.ChatId = msg.ChatID
		vote.Vote = i
		voteParsed, err := json.Marshal(vote)
		if err != nil {
			loger.LogErr(err)
			return tgbotapi.MessageConfig{}
		}
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i+1)+". "+
			chat.GetGame().GetPlayers()[i].GetFullName(), string(voteParsed))
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, buttons[i:i+1])
	}
	msg.ReplyMarkup = keyboard
	msg.Text = "Test"
	return msg
}

func startGame(update tgbotapi.Update, chat *Chat.Chat) tgbotapi.MessageConfig {
	if chat.GetGame().GetGameStage() != Game.REGISTRATION {
		return tgbotapi.MessageConfig{}
	}
	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, chat.GetRegistrationMsgId()))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if chat.GetGame().GetGameStage() == Game.REGISTRATION {
		err := chat.GetGame().NewGame()
		if err != nil {
			msg.Text = err.Error()
			if err != nil {
				loger.LogErr(err)
			}
			return msg
		}
		players := chat.GetGame().GetPlayers()
		for i := 0; i < chat.GetGame().GetNumberOfPlayers(); i++ {
			sendProfile(players[i], chat.GetLang())
		}
		db := DB.GetDataBase()
		query, err := db.Query(fmt.Sprintf("SELECT catastrophe_name, description,destruction FROM Catastrophe_%s WHERE id = $1",
			chat.GetLang()),
			chat.GetGame().GetCatastropheId())
		if err != nil {
			loger.LogErr(err)
		}

		if query.Next() {
			var (
				amenities   = ""
				capacity    = 0
				area        = 0
				time        = 0
				catastrophe = ""
				description = ""
				destruction = 0
				alive       = 0
			)
			capacity = chat.GetGame().GetNumberOfPlayers() / 2
			time = rand.Int()%17 + 3
			area = (capacity + rand.Int()%12) * 10
			query.Scan(&catastrophe, &description, &destruction)
			alive = (rand.Int()%(15-destruction) + 1) * 5
			destruction = (rand.Int()%(10+destruction) + 1) * 5
			chat.GetGame().SetAlive(alive)
			chat.GetGame().SetDestruction(destruction)

			switch chat.GetLang() {
			case Chat.RU:
				query, err = db.Query("SELECT amenities FROM bunker_ru ORDER BY RANDOM() limit 1")
				query.Next()
				query.Scan(&amenities)
				msg.Text = fmt.Sprintf(Text.DESCRIPTION_RU, catastrophe, description, destruction, alive)
				msg.Text += fmt.Sprintf(Text.BUNKER_RU, capacity, area, time, amenities)
			case Chat.EN:
				query, err = db.Query("SELECT amenities FROM bunker_en ORDER BY RANDOM() limit 1")
				query.Next()
				query.Scan(&amenities)
				msg.Text = fmt.Sprintf(Text.DESCRIPTION_EN, catastrophe, description, destruction, alive)
				msg.Text += fmt.Sprintf(Text.BUNKER_EN, capacity, area, time, amenities)
			}
			loger.LogErr(err)
			msg.ParseMode = "markdown"
		}
	}
	bot.Send(msg)

	//msg = formCharacteristicMsg(msg, chat)
	//for i := 0; i < chat.GetGame().NumberOfAlivePlayers(); i++ {
	//	msg.ChatID = int64(chat.GetGame().GetPlayers()[i].GetUser().ID)
	//	msgReply, err := bot.Send(msg)
	//	loger.LogErr(err)
	//	chat.GetGame().GetPlayers()[i].SetMsgId(msgReply.MessageID)
	//}
	msg.ReplyMarkup = nil

	msg.Text = formPlayerInfoMsg(chat)
	msg.ChatID = update.Message.Chat.ID
	return msg
}

func sendProfile(player Game.Player, lang string) {
	userId := player.GetUserId()
	var profile string
	db := DB.GetDataBase()
	s := strings.ReplaceAll(`SELECT Profession_%s.profession_name, health_%s.health, character_%s.character, 
							   baggage_%s.baggage,hobby_%s.hobby, phobias_%s.phobias, skills_%s.skills
							   FROM Profession_%s, health_%s,character_%s,baggage_%s,hobby_%s,phobias_%s,skills_%s
							   WHERE Profession_%s.id = $1 and health_%s.id = $2 and
 							   character_%s.id = $3 and baggage_%s.id = $4 and hobby_%s.id = $5 and
							   phobias_%s.id = $6 and skills_%s.id = $7`, "%s", lang)

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
		err := query.Scan(&profession, &health, &character, &baggage, &hobby, &phobias, &skills)
		biologicalCharacteristics = player.GetBioChar(lang)
		if err != nil {
			loger.LogErr(err)
		}
	}

	switch lang {

	case Chat.EN:
		profile = fmt.Sprintf(Text.PROFILE_EN, profession, health, character, baggage, biologicalCharacteristics, hobby, phobias, skills)
	case Chat.RU:
		profile = fmt.Sprintf(Text.PROFILE_RU, profession, health, character, baggage, biologicalCharacteristics, hobby, phobias, skills)
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
