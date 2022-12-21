package constants

const (
	InternalServerErr            = "internal server error"
	UserNotFoundErr              = "user not found"
	TransactionNotCreatedErr     = "transaction failed to created"
	WebsocketFailedtoConnect     = "connection failed to connected"
	WebsocketFailedtoSendMessage = "message failed to sended"

	OK       = "OK"
	Sec0     = 0
	Sec180   = 180
	Sec300   = 300
	Sec600   = 600
	Sec1800  = 1800
	Sec3600  = 3600
	Sec18000 = 18000
	Sec36000 = 36000
	Sec86400 = 86400
)

const (
	POST   = "POST"
	GET    = "GET"
	PUT    = "PUT"
	DELETE = "DELETE"
)
