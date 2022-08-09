package telegram

type UserId = int64

type UserStatus struct {
	command string
	day uint8
	month uint8
}

func userStatusNew(command string, day uint8, month uint8) UserStatus {
	return UserStatus{
		command,
		day,
		month,
	}
}

type StatusMap = map[UserId]UserStatus

type StatusUpdateCommand struct {
	id   UserId
	command string
}

func statusUpdateCommandNew(id UserId, command string) StatusUpdateCommand {
	return StatusUpdateCommand{
		id,
		command,
	}
}

type StatusUpdateMonth struct {
	id   UserId
	month uint8
}

func statusUpdateMonthNew(id UserId, month uint8) StatusUpdateMonth {
	return StatusUpdateMonth{
		id,
		month,
	}
}
