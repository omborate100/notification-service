 CREATE TABLE email_templates (
     id              BIGSERIAL PRIMARY KEY,
     template_key    VARCHAR(100) NOT NULL UNIQUE,
     template_name   VARCHAR(255) NOT NULL,
     subject_template TEXT NOT NULL,
     body_template    TEXT NOT NULL,
     description      TEXT,
     is_active        BOOLEAN NOT NULL DEFAULT TRUE,
     created_at       TIMESTAMP NOT NULL DEFAULT NOW(),
     updated_at       TIMESTAMP NOT NULL DEFAULT NOW()
 );

CREATE TABLE email_notifications (
    id BIGSERIAL PRIMARY KEY,
    template_id BIGINT NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    variables JSONB NOT NULL DEFAULT '{}'::jsonb,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    provider VARCHAR(50) NOT NULL DEFAULT 'SMTP',
    provider_message_id VARCHAR(255),
    error_message TEXT,
    sent_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_email_template
        FOREIGN KEY(template_id)
        REFERENCES email_templates(id)
);