package model

import (
	"context"
)

type Form struct {
	Name                string
	Email               string
	Phone               string
	Plan                string
	Period              string
	OnlineService       bool
	LargerStorage       bool
	CustomizableProfile bool
}

func (m *Model) AddForm(ctx context.Context, input Form) error {
	m.log.InfoContext(ctx, "start AddForm")

	sqlStatement := `
		INSERT INTO form (name, email, phone, plan, period, online_service, larger_storage, customizable_profile)
		VALUES ( ?, ?, ?, ?, ?, ?, ?, ? );
	`

	if _, err := m.db.ExecContext(
		ctx,
		sqlStatement,
		input.Name,
		input.Email,
		input.Phone,
		input.Plan,
		input.Period,
		input.OnlineService,
		input.LargerStorage,
		input.CustomizableProfile,
	); err != nil {
		m.log.ErrorContext(ctx, "fail AddForm", "error", err)
		return err
	}

	m.log.InfoContext(ctx, "success AddForm")
	return nil
}
