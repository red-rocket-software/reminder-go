package worker

import (
	"context"
	"fmt"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/red-rocket-software/reminder-go/internal/storage"
	"github.com/red-rocket-software/reminder-go/worker/mail"
)

type Worker struct {
	db  storage.ReminderRepo
	ctx context.Context
	cfg config.Config
}

func NewWorker(ctx context.Context, db storage.ReminderRepo, cfg config.Config) *Worker {
	return &Worker{
		db:  db,
		ctx: ctx,
		cfg: cfg,
	}
}

func (w *Worker) ProcessSendNotification() error {
	remindsToNotify, err := w.db.GetRemindsForNotification(w.ctx)
	if err != nil {
		return fmt.Errorf("erorr to get reminds to notification, err: %v", err)
	}

	mailer := mail.NewGmailSender(w.cfg.Email.EmailSenderName, w.cfg.Email.EmailSenderAddress, w.cfg.Email.EmailSenderPassword)

	var user model.User
	for _, remind := range remindsToNotify {
		user, err = w.db.GetUserByID(w.ctx, remind.UserID)
		if err != nil {
			return fmt.Errorf("erorr to get user, err: %v", err)
		}

		subject := "Reminder notification"
		content := fmt.Sprintf(`Hello %s,<br/>
	I wont to remember that you have something to do...<br/>
	%s, deadline to %s<br/>
	`, user.Name, remind.Description, remind.DeadlineAt)
		to := []string{user.Email}

		err = mailer.SendEmail(subject, content, to, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to send verify email: %w", err)
		}

		err = w.db.UpdateNotification(w.ctx, remind.ID, model.NotificationDAO{Notificated: true})
		if err != nil {
			return fmt.Errorf("failed to update notificated status: %w", err)
		}
		fmt.Println("Email sent successful")
	}

	return nil
}

func (w *Worker) ProcessSendDeadlineNotification() error {
	remindsToNotify, timeToDelete, err := w.db.GetRemindsForDeadlineNotification(w.ctx)
	if err != nil {
		return fmt.Errorf("erorr to get reminds to notification, err: %v", err)
	}

	mailer := mail.NewGmailSender(w.cfg.Email.EmailSenderName, w.cfg.Email.EmailSenderAddress, w.cfg.Email.EmailSenderPassword)

	var user model.User
	for _, remind := range remindsToNotify {
		user, err = w.db.GetUserByID(w.ctx, remind.UserID)
		if err != nil {
			return fmt.Errorf("erorr to get user, err: %v", err)
		}

		subject := "Reminder notification"
		content := fmt.Sprintf(`Hello %s,<br/>
	I wont to remember that you have a deadline in two hours <br/>
	%s, deadline to %s<br/>
	`, user.Name, remind.Description, remind.DeadlineAt)
		to := []string{user.Email}

		err = mailer.SendEmail(subject, content, to, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to send verify email: %w", err)
		}

		err := w.db.UpdateNotifyPeriod(w.ctx, remind.ID, timeToDelete)
		if err != nil {
			return fmt.Errorf("failed to update deadline notification period")
		}

		fmt.Println("Deadline notification Email sent successful")
	}

	return nil
}
