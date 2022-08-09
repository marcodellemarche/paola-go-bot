package telegram

type UserId = int64

type UserStatus struct {
	text string
}

func userStatusNew(text string) UserStatus {
	return UserStatus{
		text,
	}
}

type StatusMap = map[UserId]UserStatus

type StatusUpdate struct {
	id   UserId
	text string
}

func statusUpdateNew(id UserId, text string) StatusUpdate {
	return StatusUpdate{
		id,
		text,
	}
}
