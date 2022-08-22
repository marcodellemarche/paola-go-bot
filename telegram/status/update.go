package status

type Update struct {
	Id   UserId
	Next NextCommand
	Args []string
}

func SetNext(userId UserId, next NextCommand, args ...string) {
	c <- Update{
		userId,
		next,
		args,
	}
}

func ResetNext(userId UserId) {
	c <- Update{
		userId,
		nil,
		[]string{},
	}
}
