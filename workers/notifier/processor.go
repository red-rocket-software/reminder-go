package notifier

import (
	"context"
	"fmt"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/pkg/firestore"
	"github.com/red-rocket-software/reminder-go/workers/notifier/mail"
)

type Worker struct {
	todoStorage domain.TodoRepository
	fireClient  firestore.Client
	ctx         context.Context
	cfg         config.Config
}

func NewWorker(ctx context.Context, todoStorage domain.TodoRepository, fireClient firestore.Client, cfg config.Config) *Worker {
	return &Worker{
		todoStorage: todoStorage,
		fireClient:  fireClient,
		ctx:         ctx,
		cfg:         cfg,
	}
}

func (w *Worker) ProcessSendNotification() error {
	remindsToNotify, err := w.todoStorage.GetRemindsForNotification(w.ctx)
	if err != nil {
		return fmt.Errorf("erorr to get reminds to notification, err: %v", err)
	}

	mailer := mail.NewGmailSender(w.cfg.Email.EmailSenderName,
		w.cfg.Email.EmailSenderAddress,
		w.cfg.Email.EmailSenderPassword,
		w.cfg.Email.SMTPAuthAddress,
		w.cfg.Email.SMTPServerAddress)

	for _, remind := range remindsToNotify {
		user, err := w.fireClient.GetUser(remind.UserID)
		if err != nil {
			return fmt.Errorf("erorr to get user, err: %v", err)
		}

		subject := "Reminder notification"
		content := fmt.Sprintf(`Hello %s,<br/>
	I wont to remember that you have something to do...<br/><p style="color: red">
	%s <p/> deadline to %s<br/>
	`, user.UserInfo.DisplayName, remind.Description, remind.DeadlineAt)
		to := []string{user.Email}

		err = mailer.SendEmail(subject, content, to, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to send notifier email: %w", err)
		}

		err = w.todoStorage.UpdateNotification(w.ctx, remind.ID, domain.NotificationDAO{Notificated: true})
		if err != nil {
			return fmt.Errorf("failed to update notificated status: %w", err)
		}
		fmt.Println("Email sent successful")
	}

	return nil
}

func (w *Worker) ProcessSendDeadlineNotification() error {
	remindsToNotify, timeToDelete, err := w.todoStorage.GetRemindsForDeadlineNotification(w.ctx)
	if err != nil {
		return fmt.Errorf("erorr to get reminds to notification, err: %v", err)
	}

	mailer := mail.NewGmailSender(w.cfg.Email.EmailSenderName,
		w.cfg.Email.EmailSenderAddress,
		w.cfg.Email.EmailSenderPassword,
		w.cfg.Email.SMTPAuthAddress,
		w.cfg.Email.SMTPServerAddress,
	)

	for _, remind := range remindsToNotify {
		user, err := w.fireClient.GetUser(remind.UserID)
		if err != nil {
			return fmt.Errorf("erorr to get user, err: %v", err)
		}

		subject := "Reminder notification"
		content := fmt.Sprintf(`Hello %s,<br/>
	I wont to remember that you have a deadline: <br/> <p style="color: red">
	%s <p/> deadline to %s<br/>
	`, user.DisplayName, remind.Description, remind.DeadlineAt)
		to := []string{user.Email}

		err = mailer.SendEmail(subject, content, to, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to send verify email: %w", err)
		}

		err = w.todoStorage.UpdateNotifyPeriod(w.ctx, remind.ID, timeToDelete)
		if err != nil {
			return fmt.Errorf("failed to update deadline notification period")
		}

		fmt.Println("Deadline notification Email sent successful")
	}

	return nil
}
