package moneytubemodel

type UploadDataTemplate struct {
	Label string `bson:"_id"`

	UploadDataFields `bson:",inline" mapstructure:",squash"`
}
