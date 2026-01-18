package am

const (
	FailureReply = "am.Failure"
	SuccessReply = "am.Success"

	OutcomeSuccess = "SUCCESS"
	OutcomeFailure = "FAILURE"

	ReplyHandlerPrefix  = "REPLY_"
	ReplyNameHandler    = ReplyHandlerPrefix + "NAME"
	ReplyOutcomeHandler = ReplyHandlerPrefix + "OUTCOME"
)
