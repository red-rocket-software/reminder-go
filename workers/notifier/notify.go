package notifier

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/pkg/firestore"
	"github.com/red-rocket-software/reminder-go/pkg/logging"
	"github.com/red-rocket-software/reminder-go/workers/notifier/mail"
)

const (
	userMessage     string = "UserMessage"
	deadlineMessage string = "DeadlineMessage"
)

type Worker struct {
	todoStorage        domain.TodoRepository
	fireClient         firestore.Client
	ctx                context.Context
	cfg                config.Config
	log                logging.Logger
	mailer             mail.EmailSender
	userEmailQueue     chan EmailJob
	deadlineEmailQueue chan EmailJob
	timeToDelete       chan string
}

type EmailJob struct {
	To           []string
	Subject      string
	Content      string
	RemindID     int
	MessageType  string
	TimeToDelete string
}

func NewWorker(ctx context.Context, todoStorage domain.TodoRepository, fireClient firestore.Client, cfg config.Config, log logging.Logger) *Worker {
	userQueueChan := make(chan EmailJob)
	deadlineQueueChan := make(chan EmailJob)

	mailer := mail.NewGmailSender(cfg.Email.EmailSenderName,
		cfg.Email.EmailSenderAddress,
		cfg.Email.EmailSenderPassword,
		cfg.Email.SMTPAuthAddress,
		cfg.Email.SMTPServerAddress,
	)

	return &Worker{
		todoStorage:        todoStorage,
		fireClient:         fireClient,
		ctx:                ctx,
		cfg:                cfg,
		log:                log,
		userEmailQueue:     userQueueChan,
		deadlineEmailQueue: deadlineQueueChan,
		mailer:             mailer,
	}
}

func (w *Worker) Run() {
	ticker := time.NewTicker(time.Second * 5) // Workers run every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.ProcessUserNotification()
		}
	}
}

func (w *Worker) ProcessUserNotification() {
	reminds, err := w.receiveUserTodoForNotify()
	if err != nil {
		w.log.Error(err)
		return
	}

	deadlineReminds, t, err := w.receiveDeadlineTodoForNotify()
	if err != nil {
		w.log.Error(err)
		return
	}
	reminds = append(reminds, deadlineReminds...)

	if len(reminds) == 0 {
		w.log.Info("User notification: Nothing to send...")
		return
	}

	var wg sync.WaitGroup

	wg.Add(len(reminds))
	for _, remind := range reminds {
		go func(remind domain.NotificationRemind) {
			//defer wg.Done()

			user, err := w.fireClient.GetUser(remind.UserID)
			if err != nil {
				w.log.Error("error getting user for notification:", err)
				return
			}

			job := EmailJob{
				RemindID: remind.ID,
				To:       []string{user.Email},
				Subject:  "Reminder notification",
				Content: fmt.Sprintf(`Hello %s,<br/>
						I want to remind you that you have something to do...<br/><p style="color: red">
						%s <p/> deadline to %s<br/>
						`, user.UserInfo.DisplayName, remind.Description, remind.DeadlineAt),
				TimeToDelete: t,
				MessageType:  remind.MessageType,
			}

			w.userEmailQueue <- job
		}(remind)
	}

	go func() {
		wg.Wait()
		close(w.userEmailQueue)
	}()

	go w.ProcessSendUserNotification(&wg)
}

func (w *Worker) ProcessSendUserNotification(wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range w.userEmailQueue {
		err := w.mailer.SendEmail(job.Subject, job.Content, job.To, nil, nil, nil)
		if err != nil {
			w.log.Error("failed to send notifier email:", err)
			continue // Continue processing other jobs even if one fails
		}

		if job.MessageType == userMessage {
			err = w.todoStorage.UpdateNotification(w.ctx, job.RemindID, domain.NotificationDAO{Notificated: true})
			if err != nil {
				w.log.Error("failed to update notificated status:", err)
				continue // Continue processing other jobs even if one fails
			}
		}

		if job.MessageType == deadlineMessage {
			err = w.todoStorage.UpdateNotifyPeriod(w.ctx, job.RemindID, job.TimeToDelete)
			if err != nil {
				w.log.Error("failed to update deadline notification period")
				return
			}
		}

		w.log.Infof("Email sent successfully, remind id:%d", job.RemindID)
	}
}

func (w *Worker) receiveUserTodoForNotify() ([]domain.NotificationRemind, error) {
	remindsToNotify, err := w.todoStorage.GetRemindsForNotification(w.ctx)
	if err != nil {
		return nil, fmt.Errorf("erorr to get user reminds to notification, err: %v", err)
	}

	for i, _ := range remindsToNotify {
		remindsToNotify[i].MessageType = userMessage
	}

	return remindsToNotify, nil
}

func (w *Worker) receiveDeadlineTodoForNotify() ([]domain.NotificationRemind, string, error) {
	remindsToNotify, timeToDelete, err := w.todoStorage.GetRemindsForDeadlineNotification(w.ctx)
	if err != nil {
		return nil, "", fmt.Errorf("erorr to get dealine reminds to notification, err: %v", err)
	}

	for i, _ := range remindsToNotify {
		remindsToNotify[i].MessageType = deadlineMessage
	}

	return remindsToNotify, timeToDelete, nil
}
