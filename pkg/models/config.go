package models

type CharacterInput struct {
	ID    int    `yaml:"id"`
	Token string `yaml:"token"`
	Name  string `yaml:"name"`
}

type NotificationInput struct {
	Alliances    []NotificationAlliancesInput  `yaml:"alliances"`
	Corporations []NotificationCoporationInput `yaml:"corporations"`
	Characters   []NotificationCharacterInput  `yaml:"characters"`
}

type NotificationAlliancesInput struct {
	Character int32 `yaml:"character_id"`
	ID        int32 `yaml:"id"`
}

type NotificationCoporationInput struct {
	Character int32 `yaml:"character_id"`
	ID        int32 `yaml:"id"`
}

type NotificationCharacterInput struct {
	ID    int32  `yaml:"id"`
	Token string `yaml:"token"`
}
