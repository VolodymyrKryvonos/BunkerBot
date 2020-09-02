package Text

const (
	TEST = "test"

	CHANGED_LANG_EN = "The bot language has been changed to English \xF0\x9F\x87\xBA\xF0\x9F\x87\xB8"
	CHANGED_LANG_RU = "Язык бота был изменен на русский \xF0\x9F\x87\xB7\xF0\x9F\x87\xBA"

	PROFILE_EN = "\xF0\x9F\x91\xB7Profession: %s\n\xE2\x9D\xA4Health status: %s\n\xF0\x9F\x91\xA5Character: %s\n\xF0\x9F\x92\xBCBaggage: %s\n" +
		"\xF0\x9F\x91\xABBiological characteristics: %s\n\xF0\x9F\x8F\x82Hobby: %s\n\xF0\x9F\x90\x8DPhobia: %s\n\xF0\x9F\x92\xAAPersonal skill: %s\n"
	PROFILE_RU = "\xF0\x9F\x91\xB7Профессия: %s \n\xE2\x9D\xA4Состояние здоровья: %s\n\xF0\x9F\x91\xA5Характер: %s\n\xF0\x9F\x92\xBCБагаж: %s\n" +
		"\xF0\x9F\x91\xABБиологические характеристики: %s \n\xF0\x9F\x8F\x82Хобби: %s \n\xF0\x9F\x90\x8DФобия: %s \n\xF0\x9F\x92\xAAЛичные навыки: %s"

	HEALTH_EN = "\xE2\x9D\xA4Health status"
	HEALTH_RU = "\xE2\x9D\xA4Состояние здоровья"

	CHARACTER_EN = "\xF0\x9F\x91\xA5Character"
	CHARACTER_RU = "\xF0\x9F\x91\xA5Характер"

	BAGGAGE_EN = "\xF0\x9F\x92\xBCBaggage"
	BAGGAGE_RU = "\xF0\x9F\x92\xBCБагаж"

	BIO_EN = "\xF0\x9F\x91\xABBiological characteristics"
	BIO_RU = "\xF0\x9F\x91\xABБиологические характеристики"

	HOBBY_EN = "\xF0\x9F\x8F\x82Hobby"
	HOBBY_RU = "\xF0\x9F\x8F\x82Хобби"

	PHOBIA_EN = "\xF0\x9F\x90\x8DPhobia"
	PHOBIA_RU = "\xF0\x9F\x90\x8DФобия"

	SKILL_EN = "\xF0\x9F\x92\xAAPersonal skill"
	SKILL_RU = "\xF0\x9F\x92\xAAЛичные навыки"

	REGISTRATION_EN = "Registration is open \xF0\x9F\x9A\xAA"
	REGISTRATION_RU = "Ведётся набор в игру \xF0\x9F\x9A\xAA"

	JOIN_EN = "Join"
	JOIN_RU = "Присоединиться"

	LEAVE_EN = "Leave \xF0\x9F\x9A\xB6"
	LEAVE_RU = "Покинуть \xF0\x9F\x9A\xB6"

	ALREADY_REGISTRED_EN = "You are already in the game \xE2\x9D\x97"
	ALREADY_REGISTRED_RU = "Вы уже в игре \xE2\x9D\x97"

	NUMBER_OF_PLAYERS_RU = "Количество игроков: "
	NUMBER_OF_PLAYERS_EN = "Number of players: "

	SETTINGS = "Grant administrator rights to the bot.\n\nTo start the game give me the following administrator rights: \n\xE2\x9C\x85 delete messages \n\xE2\x9C\x85 block users \n\xE2\x9C\x85 pin messages"
	LANGUAGE = "\xF0\x9F\x87\xB7\xF0\x9F\x87\xBA Language \xF0\x9F\x87\xBA\xF0\x9F\x87\xB8"
	DONE     = "Done \xF0\x9F\x91\x8D"

	DESCRIPTION_RU = "*%s*\n\n%s\n\nРазрушение: %v%%\nВыживших: %v%%\n"
	DESCRIPTION_EN = "*%s*\n\n%s\n\nDestruction: %v%%\nSurvivors: %v%%\n"

	BUNKER_RU = "Вместимость бункера: %v чел.\nПлощадь бункера: %v м2,\nВремя пребывания: %v мес.\n%s"
	BUNKER_EN = "Bunker capacity: %v people. \nBunker area: %v m2,\nStay time: %v months \n%s"

	VOTE_AGAINST_YOURSELF_RU = "Вы не можете голосовать против себя \xE2\x9D\x97"
	VOTE_AGAINST_YOURSELF_EN = "You cannot vote against yourself \xE2\x9D\x97"

	SELECT_CHARACTERISTIC_RU = "Выберите характеристику, которую хотите раскрыть"
	SELECT_CHARACTERISTIC_EN = "Select the characteristic you want to reveal"

	CHARACTERISTIC_ALREADY_OPENED_RU = "Вы уже вскрыли эту характеристику"
	CHARACTERISTIC_ALREADY_OPENED_EN = "You have already revealed this characteristic"

	KICK_RU = "Вы решили оставить %s за стенами бункера"
	KICK_EN = "You decided to leave %s outside the walls of the bunker"

	VOTE_AGAINST_RU= "Выберите кого вы не возьмете с собой в бункер"
	VOTE_AGAINST_EN= "Choose who you won't take to the bunker"

	MORE_THAN_ONE_RU = "имеют наибольшее количество голосов"
	MORE_THAN_ONE_EN = "have the largest number of votes"
)
