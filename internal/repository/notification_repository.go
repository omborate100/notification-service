package repository

import (
	"context"
	"encoding/json"

	"notification-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}
func (r *NotificationRepository) Create(
	notification *model.EmailNotification,
) (int64, error) {

	variablesJSON, err := json.Marshal(notification.Variables)

	if err != nil {
		return 0, err
	}

	query := `
	INSERT INTO email_notifications
	(
		template_id,
		recipient_email,
		subject,
		body,
		variables,
		status,
		provider
	)
	VALUES
	(
		$1,$2,$3,$4,$5,$6,$7
	)
	RETURNING id
	`

	var id int64

	err = r.db.QueryRow(
		context.Background(),
		query,
		notification.TemplateID,
		notification.RecipientEmail,
		notification.Subject,
		notification.Body,
		variablesJSON,
		notification.Status,
		notification.Provider,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *NotificationRepository) MarkSent(
	id int64,
	providerMessageID string,
) error {

	query := `
	UPDATE email_notifications
	SET
		status=$1,
		provider_message_id=$2,
		sent_at=NOW(),
		updated_at=NOW()
	WHERE id=$3
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		model.StatusSent,
		providerMessageID,
		id,
	)

	return err
}

func (r *NotificationRepository) MarkFailed(
	id int64,
	errorMessage string,
) error {

	query := `
	UPDATE email_notifications
	SET
		status=$1,
		error_message=$2,
		updated_at=NOW()
	WHERE id=$3
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		model.StatusFailed,
		errorMessage,
		id,
	)

	return err
}