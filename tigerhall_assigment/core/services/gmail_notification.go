package services

import "log"

type MailNotification struct {
	//mail over which smtp is registered
	MailId string
	//SMTP or anything
	Provide string
	//Could be fetched from vault
	Cred string
}
type MailNotificationRequest struct {
	Body     string
	Template string
	Header   string
	To       []string
}

func NewMailNotification() Notification {
	return &MailNotification{}
}

func (mail *MailNotification) SendNotification(request interface{}) {
	mailReq, ok := request.(MailNotificationRequest)
	if !ok {
		log.Printf("[Mail notification] Request not of valid type")
	}
	for _, curEmail := range mailReq.To {
		log.Printf("[Mail notification] Sent request for following mail = %s", curEmail)
	}

}
