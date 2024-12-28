package bot_commander

type Option func(bot *BotCommander)

func WithLogger(l Logger) Option {
	return func(bot *BotCommander) {
		bot.logger = l
	}
}

func WithTgAPI(api BotClient) Option {
	return func(bot *BotCommander) {
		bot.tg = api
	}
}

func WithEmailClient(ec EmailClient) Option {
	return func(bot *BotCommander) {
		bot.emailClient = ec
	}
}
