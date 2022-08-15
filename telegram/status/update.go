package status

type Update struct {
	Id   UserId
	Next NextFunc
	Args []string
}

func SetNext(userId UserId, next NextFunc, args ...string) {
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
