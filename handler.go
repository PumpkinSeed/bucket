package odatas

type Handler struct {
	state string
}

type Configuration struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	BucketUsername string `json:"bucket_username"`
	BucketPassword string `json:"bucket_password"`
}

func New(c *Configuration) Handler {
	return Handler{}
}