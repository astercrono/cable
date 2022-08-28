package event

const ActionRefreshConnection string = "refresh"
const ActionCloseConnection string = "close"

func IsAction(msg string, action string) bool {
	return msg == action
}

func IsRefresh(msg string) bool {
	return msg == ActionRefreshConnection
}

func IsClose(msg string) bool {
	return msg == ActionCloseConnection
}
