package repository

import (
	"context"
	"errors"

	"notification-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TemplateRepository struct {
	db *pgxpool.Pool
}

func NewTemplateRepository(db *pgxpool.Pool) *TemplateRepository {
	return &TemplateRepository{
		db: db,
	}
}

func (r *TemplateRepository) GetByTemplateKey(templateKey string) (*model.EmailTemplate, error) {

	query := `
	SELECT
		id,
		template_key,
		template_name,
		subject_template,
		body_template,
		description,
		is_active,
		created_at,
		updated_at
	FROM email_templates
	WHERE template_key = $1
	AND is_active = true
	`

	var template model.EmailTemplate

	err := r.db.QueryRow(
		context.Background(),
		query,
		templateKey,
	).Scan(
		&template.ID,
		&template.TemplateKey,
		&template.TemplateName,
		&template.SubjectTemplate,
		&template.BodyTemplate,
		&template.Description,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("template not found")
	}

	return &template, nil
}