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

type Worker struct {
	ticker             time.Ticker
	todoStorage        domain.TodoRepository
	fireClient         firestore.Client
	ctx                context.Context
	cfg                config.Config
	log                logging.Logger
	userEmailQueue     chan EmailJob
	deadlineEmailQueue chan EmailJob
	timeToDelete       chan string
}

type EmailJob struct {
	To       []string
	Subject  string
	Content  string
	RemindID int
}

func NewWorker(ctx context.Context, todoStorage domain.TodoRepository, fireClient firestore.Client, ticker time.Ticker, cfg config.Config, log logging.Logger) *Worker {
	userQueueChan := make(chan EmailJob)
	deadlineQueueChan := make(chan EmailJob)
	return &Worker{
		todoStorage:        todoStorage,
		fireClient:         fireClient,
		ctx:                ctx,
		cfg:                cfg,
		ticker:             ticker,
		log:                log,
		userEmailQueue:     userQueueChan,
		deadlineEmailQueue: deadlineQueueChan,
	}
}

func (w *Worker) Run() {
	//w.ProcessReceiveUserNotification()
	w.ProcessReceiveDeadlineNotification()

	var wg sync.WaitGroup

	select {
	case <-w.deadlineEmailQueue:
		w.ProcessSendDeadlineNotification(&wg)
	}

	//go w.ProcessSendUserNotification(&wg)

	//wg.Wait()
}

func (w *Worker) ProcessReceiveUserNotification() {
	reminds, err := w.receiveUserTodoForNotify()
	if err != nil {
		w.log.Error(err)
		return
	}

	if len(reminds) == 0 {
		w.log.Info("Nothing send to user notification...")
		return
	}

	var wg sync.WaitGroup

	for _, remind := range reminds {
		wg.Add(1)

		go func() {
			defer wg.Done()

			user, err := w.fireClient.GetUser(remind.UserID)
			if err != nil {
				w.log.Error("error to get user for notification")
				return
			}

			w.userEmailQueue <- EmailJob{
				RemindID: remind.ID,
				To:       []string{user.Email},
				Subject:  "Reminder notification",
				Content: fmt.Sprintf(`Hello %s,<br/>
					 I wont to remember that you have something to do...<br/><p style="color: red">
					 %s <p/> deadline to %s<br/>
					 `, user.UserInfo.DisplayName, remind.Description, remind.DeadlineAt),
			}
		}()

		wg.Wait()
	}
}

func (w *Worker) ProcessReceiveDeadlineNotification() {
	reminds, timeToDelete, err := w.receiveDeadlineTodoForNotify()
	if err != nil {
		w.log.Error(err)
		return
	}

	if len(reminds) == 0 {
		w.log.Info("Nothing send to deadline notify...")
		return
	}

	w.timeToDelete <- timeToDelete

	var wg sync.WaitGroup

	for _, remind := range reminds {
		wg.Add(1)

		go func() {
			defer wg.Done()

			user, err := w.fireClient.GetUser(remind.UserID)
			if err != nil {
				w.log.Error("error to get user for notification")
				return
			}

			w.deadlineEmailQueue <- EmailJob{
				RemindID: remind.ID,
				To:       []string{user.Email},
				Subject:  "Reminder notification",
				Content: fmt.Sprintf(`Hello %s,<br/>
						I wont to remember that you have a deadline: <br/> <p style="color: red">
						%s <p/> deadline to %s<br/>
						`, user.DisplayName, remind.Description, remind.DeadlineAt),
			}
		}()

		wg.Wait()
	}
}

func (w *Worker) ProcessSendUserNotification(wg *sync.WaitGroup) {
	defer wg.Done()

	mailer := mail.NewGmailSender(w.cfg.Email.EmailSenderName,
		w.cfg.Email.EmailSenderAddress,
		w.cfg.Email.EmailSenderPassword,
		w.cfg.Email.SMTPAuthAddress,
		w.cfg.Email.SMTPServerAddress,
	)

	for job := range w.userEmailQueue {
		err := mailer.SendEmail(job.Subject, job.Content, job.To, nil, nil, nil)
		if err != nil {
			w.log.Error("error failed to send notifier email:", err)
			return
		}

		err = w.todoStorage.UpdateNotification(w.ctx, job.RemindID, domain.NotificationDAO{Notificated: true})
		if err != nil {
			w.log.Error("failed to update notificated status:", err)
			return
		}

		w.log.Infof("Email sent successful, remind id:%d", job.RemindID)
	}
}

func (w *Worker) ProcessSendDeadlineNotification(wg *sync.WaitGroup) {
	//defer wg.Done()

	mailer := mail.NewGmailSender(w.cfg.Email.EmailSenderName,
		w.cfg.Email.EmailSenderAddress,
		w.cfg.Email.EmailSenderPassword,
		w.cfg.Email.SMTPAuthAddress,
		w.cfg.Email.SMTPServerAddress,
	)

	for job := range w.deadlineEmailQueue {
		err := mailer.SendEmail(job.Subject, job.Content, job.To, nil, nil, nil)
		if err != nil {
			w.log.Error("error failed to send notifier email:", err)
			return
		}

		err = w.todoStorage.UpdateNotification(w.ctx, job.RemindID, domain.NotificationDAO{Notificated: true})
		if err != nil {
			w.log.Error(err)
			return
		}

		w.log.Infof("Email sent successful, remind id:%d", job.RemindID)
	}
}

func (w *Worker) receiveUserTodoForNotify() ([]domain.NotificationRemind, error) {
	remindsToNotify, err := w.todoStorage.GetRemindsForNotification(w.ctx)
	if err != nil {
		return nil, fmt.Errorf("erorr to get user reminds to notification, err: %v", err)
	}

	return remindsToNotify, nil
}

func (w *Worker) receiveDeadlineTodoForNotify() ([]domain.NotificationRemind, string, error) {
	remindsToNotify, timeToDelete, err := w.todoStorage.GetRemindsForDeadlineNotification(w.ctx)
	if err != nil {
		return nil, "", fmt.Errorf("erorr to get dealine reminds to notification, err: %v", err)
	}

	return remindsToNotify, timeToDelete, nil
}
