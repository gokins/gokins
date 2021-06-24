package hook

type (
	Reference struct {
		Name string
		Path string
		Sha  string
	}

	Commit struct {
		Sha     string
		Message string
		Link    string
	}

	User struct {
		UserName string
	}
)
