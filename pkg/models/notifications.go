package models

type ContactNotificationInput struct {
	Event       string
	EntityType  string
	EntityName  string
	ContactType string
	ContactName string
	Standing    *float32
	Labels      []string
	OldLabels   *[]string
	OldStanding *float32
}
