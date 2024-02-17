package services

import (
	"log"

	"github.com/nitin/tigerhall/core/internal/repositiories"
)

const (
	Msg       = "some custom msg"
	queueSize = 1000
)

type SightingNotification struct {
	queue       chan int
	noOfWorkers int
	userRepo    repositiories.UserRepo
	notService  Notification
}

func NewSightingNotification(noOfWorkers int, userRepo repositiories.UserRepo, notService Notification) *SightingNotification {
	service := &SightingNotification{
		queue:       make(chan int, queueSize),
		noOfWorkers: noOfWorkers,
		userRepo:    userRepo,
		notService:  notService,
	}
	for count := 0; count < noOfWorkers; count++ {
		go service.sightingWorkers()
	}
	log.Println("[SightingNotification] Sighting Workers intialized!!!!")
	//Intilaise all the go routines over here
	return service
}

func (sighting *SightingNotification) NotifyUser(userId int) {
	sighting.queue <- userId
}

func (sighting *SightingNotification) sightingWorkers() {
	for curId := range sighting.queue {

		email := sighting.userRepo.UserById(curId).Email

		//send notification for email from over here
		sighting.notService.SendNotification(MailNotificationRequest{
			To:     []string{email},
			Body:   Msg,
			Header: " some custom header",
		})
	}
}
