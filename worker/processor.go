package worker

//
//import (
//	"context"
//	"fmt"
//
//	"github.com/red-rocket-software/reminder-go/config"
//	todoModel "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
//	"github.com/red-rocket-software/reminder-go/internal/reminder/storage"
//	model "github.com/red-rocket-software/reminder-go/internal/user/domain"
//	"github.com/red-rocket-software/reminder-go/worker/mail"
//)
//
//type Worker struct {
//	todoStorage storage.ReminderRepo
//	ctx         context.Context
//	cfg         config.Config
//}
//
//func NewWorker(ctx context.Context, todoStorage storage.ReminderRepo, cfg config.Config) *Worker {
//	return &Worker{
//		todoStorage: todoStorage,
//		ctx:         ctx,
//		cfg:         cfg,
//	}
//}
//
//func (w *Worker) ProcessSendNotification() error {
//	remindsToNotify, err := w.todoStorage.GetRemindsForNotification(w.ctx)
//	if err != nil {
//		return fmt.Errorf("erorr to get reminds to notification, err: %v", err)
//	}
//
//	mailer := mail.NewGmailSender(w.cfg.Email.EmailSenderName, w.cfg.Email.EmailSenderAddress, w.cfg.Email.EmailSenderPassword)
//
//	var user model.User
//	for _, remind := range remindsToNotify {
//		user, err = w.todoStorage.GetUserByID(w.ctx, remind.UserID)
//		if err != nil {
//			return fmt.Errorf("erorr to get user, err: %v", err)
//		}
//
//		subject := "Reminder notification"
//		content := fmt.Sprintf(`Hello %s,<br/>
//	I wont to remember that you have something to do...<br/><p style="color: red">
//	%s <p/> deadline to %s<br/>
//	`, user.Name, remind.Description, remind.DeadlineAt)
//		to := []string{user.Email}
//
//		err = mailer.SendEmail(subject, content, to, nil, nil, nil)
//		if err != nil {
//			return fmt.Errorf("failed to send verify email: %w", err)
//		}
//
//		err = w.todoStorage.UpdateNotification(w.ctx, remind.ID, todoModel.NotificationDAO{Notificated: true})
//		if err != nil {
//			return fmt.Errorf("failed to update notificated status: %w", err)
//		}
//		fmt.Println("Email sent successful")
//	}
//
//	return nil
//}
//
//func (w *Worker) ProcessSendDeadlineNotification() error {
//	remindsToNotify, timeToDelete, err := w.todoStorage.GetRemindsForDeadlineNotification(w.ctx)
//	if err != nil {
//		return fmt.Errorf("erorr to get reminds to notification, err: %v", err)
//	}
//
//	mailer := mail.NewGmailSender(w.cfg.Email.EmailSenderName, w.cfg.Email.EmailSenderAddress, w.cfg.Email.EmailSenderPassword)
//
//	var user model.User
//	for _, remind := range remindsToNotify {
//		user, err = w.todoStorage.GetUserByID(w.ctx, remind.UserID)
//		if err != nil {
//			return fmt.Errorf("erorr to get user, err: %v", err)
//		}
//
//		subject := "Reminder notification"
//		content := fmt.Sprintf(`Hello %s,<br/>
//	I wont to remember that you have a deadline: <br/> <p style="color: red">
//	%s <p/> deadline to %s<br/>
//	`, user.Name, remind.Description, remind.DeadlineAt)
//		to := []string{user.Email}
//
//		err = mailer.SendEmail(subject, content, to, nil, nil, nil)
//		if err != nil {
//			return fmt.Errorf("failed to send verify email: %w", err)
//		}
//
//		err := w.todoStorage.UpdateNotifyPeriod(w.ctx, remind.ID, timeToDelete)
//		if err != nil {
//			return fmt.Errorf("failed to update deadline notification period")
//		}
//
//		fmt.Println("Deadline notification Email sent successful")
//	}
//
//	return nil
//}
