package sessions

import (
	"context"
	"fmt"
	"log/slog"
	"olimotracker/pkg/db"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, session *Session) (*uuid.UUID, error)
	GetByID(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (*SessionResponse, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*SessionResponse, error)
	GetByCategoryID(ctx context.Context, categoryID uuid.UUID, userID uuid.UUID) ([]*SessionResponse, error)
	Update(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, session *Session) (*uuid.UUID, error)
	Delete(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) error
	GetMinutesByCategoryForUser(ctx context.Context, userID *uuid.UUID) ([]CategoryMinutes, error)
	GetMinutesBySessionForUser(ctx context.Context, userID *uuid.UUID) ([]SessionsMinutes, error)
}

type repository struct {
	db db.DBClient
	l  *slog.Logger
}

func NewRepository(db db.DBClient, l *slog.Logger) Repository {
	return &repository{db: db, l: l}
}

func (r *repository) Create(ctx context.Context, session *Session) (*uuid.UUID, error) {
	var id uuid.UUID

	builder := squirrel.Insert(db.SessionsTable).
		Columns(db.SessionsCategoryIDColumn, db.SessionsUserIDColumn, db.SessionsDurationColumn, db.SessionsNotesColumn).
		Values(session.CategoryID, session.UserID, session.Duration, session.Note).
		Suffix("RETURNING " + db.SessionsIDColumn).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error building query: ", "err", err)
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		r.l.Error("error scanning row: ", "err", err)
		return nil, err
	}

	return &id, nil
}

func (r *repository) GetByID(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (*SessionResponse, error) {
	var session SessionResponse

	builder := squirrel.Select(fmt.Sprintf("s.%v, s.%v, s.%v, s.%v, s.%v, s.%v, c.%v, c.%v", db.SessionsIDColumn, db.SessionsUserIDColumn, db.SessionsCategoryIDColumn, db.SessionsDurationColumn, db.SessionsNotesColumn, db.SessionsCreatedAtColumn, db.CategoriesTitleColumn, db.CategoriesColorColumn)).
		From(db.SessionsTable + " s").
		LeftJoin(fmt.Sprintf("%v c ON c.%v = s.%v", db.CategoriesTable, db.CategoriesIDColumn, db.SessionsCategoryIDColumn)).
		Where(squirrel.Eq{"s." + db.SessionsIDColumn: sessionID, "s." + db.SessionsUserIDColumn: userID}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error building query: ", "err", err)
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&session.ID, &session.UserID, &session.CategoryID, &session.Duration, &session.Note, &session.CreatedAt, &session.CategoryTitle, &session.CategoryColor)
	if err != nil {
		r.l.Error("error scanning row: ", "err", err)
		return nil, err
	}

	return &session, nil
}

func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*SessionResponse, error) {
	var sessions []*SessionResponse

	builder := squirrel.Select(fmt.Sprintf("s.%v, s.%v, s.%v, s.%v, s.%v, s.%v, c.%v, c.%v", db.SessionsIDColumn, db.SessionsUserIDColumn, db.SessionsCategoryIDColumn, db.SessionsDurationColumn, db.SessionsNotesColumn, db.SessionsCreatedAtColumn, db.CategoriesTitleColumn, db.CategoriesColorColumn)).
		From(db.SessionsTable + " s").
		LeftJoin(fmt.Sprintf("%v c ON c.%v = s.%v", db.CategoriesTable, db.CategoriesIDColumn, db.SessionsCategoryIDColumn)).
		Where(squirrel.Eq{"s." + db.SessionsUserIDColumn: userID}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error building query: ", "err", err)
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		r.l.Error("error querying database: ", "err", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var session SessionResponse
		err := rows.Scan(&session.ID, &session.UserID, &session.CategoryID, &session.Duration, &session.Note, &session.CreatedAt, &session.CategoryTitle, &session.CategoryColor)
		if err != nil {
			r.l.Error("error scanning row: ", "err", err)
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

func (r *repository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID, userID uuid.UUID) ([]*SessionResponse, error) {
	var sessions []*SessionResponse

	builder := squirrel.Select(fmt.Sprintf("s.%v, s.%v, s.%v, s.%v, s.%v, s.%v, c.%v, c.%v", db.SessionsIDColumn, db.SessionsUserIDColumn, db.SessionsCategoryIDColumn, db.SessionsDurationColumn, db.SessionsNotesColumn, db.SessionsCreatedAtColumn, db.CategoriesTitleColumn, db.CategoriesColorColumn)).
		From(db.SessionsTable + " s").
		LeftJoin(fmt.Sprintf("%v c ON c.%v = s.%v", db.CategoriesTable, db.CategoriesIDColumn, db.SessionsCategoryIDColumn)).
		Where(squirrel.Eq{
			"s." + db.SessionsCategoryIDColumn: categoryID,
			"s." + db.SessionsUserIDColumn:     userID,
		}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error building query: ", "err", err)
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		r.l.Error("error querying database: ", "err", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var session SessionResponse
		err := rows.Scan(&session.ID, &session.UserID, &session.CategoryID, &session.Duration, &session.Note, &session.CreatedAt, &session.CategoryTitle, &session.CategoryColor)
		if err != nil {
			r.l.Error("error scanning row: ", "err", err)
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

func (r *repository) Update(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, session *Session) (*uuid.UUID, error) {
	var id uuid.UUID

	builder := squirrel.Update(db.SessionsTable).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{db.SessionsIDColumn: sessionID, db.SessionsUserIDColumn: userID}).
		Suffix("RETURNING " + db.CategoriesIDColumn)

	if session.CategoryID != nil {
		builder = builder.Set(db.SessionsCategoryIDColumn, session.CategoryID)
	}

	if session.Duration > 0 {
		builder = builder.Set(db.SessionsDurationColumn, session.Duration)
	}

	if session.Note != nil {
		builder = builder.Set(db.SessionsNotesColumn, session.Note)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error building query: ", "err", err)
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		r.l.Error("error scanning row: ", "err", err)
		return nil, err
	}

	return &id, nil
}

func (r *repository) Delete(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) error {
	builder := squirrel.Delete(db.SessionsTable).
		Where(squirrel.Eq{db.SessionsIDColumn: sessionID, db.SessionsUserIDColumn: userID}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error building query: ", "err", err)
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		r.l.Error("error executing query: ", "err", err)
		return err
	}

	return nil
}

func (r *repository) GetMinutesByCategoryForUser(ctx context.Context, userID *uuid.UUID) ([]CategoryMinutes, error) {
	builder := squirrel.Select(fmt.Sprintf("c.%v, c.%v, c.%v, SUM(s.%v) AS minutes", db.CategoriesIDColumn, db.CategoriesTitleColumn, db.CategoriesColorColumn, db.SessionsDurationColumn)).
		From(db.SessionsTable + " s").
		Where(squirrel.Eq{"s." + db.SessionsUserIDColumn: *userID}).
		Join(fmt.Sprintf("%v c ON s.%v = c.%v", db.CategoriesTable, db.SessionsCategoryIDColumn, db.CategoriesIDColumn)).
		PlaceholderFormat(squirrel.Dollar).
		GroupBy(fmt.Sprintf("c.%v, c.%v, c.%v", db.CategoriesIDColumn, db.CategoriesTitleColumn, db.CategoriesColorColumn))

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error building query: ", "err", err)
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		r.l.Error("error executing query: ", "err", err)
		return nil, err
	}
	defer rows.Close()

	var result []CategoryMinutes
	for rows.Next() {
		var categoryTitle string
		var categoryID uuid.UUID
		var categoryColor string
		var minutes int
		if err := rows.Scan(&categoryID, &categoryTitle, &categoryColor, &minutes); err != nil {
			r.l.Error("error scanning row: ", "err", err)
			return nil, err
		}
		result = append(result, CategoryMinutes{
			CategoryID:    categoryID,
			CategoryTitle: categoryTitle,
			CategoryColor: categoryColor,
			Minutes:       minutes,
		})
	}
	if err := rows.Err(); err != nil {
		r.l.Error("error iterating rows: ", "err", err)
		return nil, err
	}

	return result, nil
}

func (r *repository) GetMinutesBySessionForUser(ctx context.Context, userID *uuid.UUID) ([]SessionsMinutes, error) {
	builder := squirrel.Select(fmt.Sprintf("SUM(%v), DATE(%v) AS day", db.SessionsDurationColumn, db.SessionsCreatedAtColumn)).
		From(db.SessionsTable).
		Where(squirrel.Eq{db.SessionsUserIDColumn: userID}).
		GroupBy("DATE(" + db.SessionsCreatedAtColumn + ")").
		OrderBy("day").
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error building query: ", "err", err)
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		r.l.Error("error executing query: ", "err", err)
		return nil, err
	}
	defer rows.Close()

	var result []SessionsMinutes
	for rows.Next() {
		var minutes int
		var date time.Time
		if err := rows.Scan(&minutes, &date); err != nil {
			r.l.Error("error scanning row: ", "err", err)
			return nil, err
		}
		result = append(result, SessionsMinutes{
			Minutes: minutes,
			Date:    date,
		})
	}

	return result, nil
}
