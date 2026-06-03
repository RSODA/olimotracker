package categories

import (
	"context"
	"log/slog"
	"olimotracker/pkg/db"

	"github.com/Masterminds/squirrel"
)

type Repository interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	GetCategoriesByUserID(ctx context.Context, userID string) ([]*Category, error)
	GetCategoryByID(ctx context.Context, categoryID string, userID string) (*Category, error)
	UpdateCategory(ctx context.Context, category *Category) (*Category, error)
	DeleteCategory(ctx context.Context, categoryID string, userID string) error
}

type repository struct {
	db db.DBClient
	l  *slog.Logger
}

func NewRepository(db db.DBClient, l *slog.Logger) Repository {
	return &repository{db: db, l: l}
}

func (r *repository) CreateCategory(ctx context.Context, category *Category) (*Category, error) {
	builder := squirrel.Insert(db.CategoriesTable).
		Columns(db.CategoriesUserIDColumn, db.CategoriesColorColumn, db.CategoriesTitleColumn).
		Values(category.UserID, category.Color, category.Title).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING " + db.CategoriesIDColumn + ", " + db.CategoriesCreatedAtColumn)

	query, args, err := builder.ToSql()
	if err != nil {

		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&category.ID, &category.CreatedAt)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *repository) GetCategoriesByUserID(ctx context.Context, userID string) ([]*Category, error) {
	var res []*Category

	builder := squirrel.Select(db.CategoriesIDColumn, db.CategoriesUserIDColumn, db.CategoriesTitleColumn, db.CategoriesColorColumn, db.CategoriesCreatedAtColumn).
		From(db.CategoriesTable).
		Where(squirrel.Eq{db.CategoriesUserIDColumn: userID}).
		OrderBy(db.CategoriesCreatedAtColumn + " DESC").
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error creating query for GetCategoriesByUserID: ", "err", err, "query", query)
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category Category
		err = rows.Scan(&category.ID, &category.UserID, &category.Title, &category.Color, &category.CreatedAt)
		if err != nil {
			r.l.Error("error scanning row for GetCategoriesByUserID: ", "err", err)
			return nil, err
		}
		res = append(res, &category)
	}

	return res, nil
}

func (r *repository) GetCategoryByID(ctx context.Context, categoryID string, userID string) (*Category, error) {
	var res Category

	builder := squirrel.Select(db.CategoriesIDColumn, db.CategoriesUserIDColumn, db.CategoriesTitleColumn, db.CategoriesColorColumn, db.CategoriesCreatedAtColumn).
		From(db.CategoriesTable).
		Where(squirrel.Eq{db.CategoriesIDColumn: categoryID, db.CategoriesUserIDColumn: userID}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error creating query for GetCategoryByID: ", "err", err, "query", query)
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&res.ID, &res.UserID, &res.Title, &res.Color, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *repository) UpdateCategory(ctx context.Context, category *Category) (*Category, error) {
	builder := squirrel.Update(db.CategoriesTable).
		Where(squirrel.Eq{db.CategoriesIDColumn: category.ID}).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING " + db.CategoriesUserIDColumn + ", " + db.CategoriesCreatedAtColumn)

	if len(category.Title) > 0 {
		builder = builder.Set(db.CategoriesTitleColumn, category.Title)
	}
	if len(category.Color) > 0 {
		builder = builder.Set(db.CategoriesColorColumn, category.Color)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error creating query for UpdateCategory: ", "err", err, "query", query)
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&category.UserID, &category.CreatedAt)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *repository) DeleteCategory(ctx context.Context, categoryID string, userID string) error {
	builder := squirrel.Delete(db.CategoriesTable).
		Where(squirrel.Eq{db.CategoriesIDColumn: categoryID, db.CategoriesUserIDColumn: userID}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error creating query for DeleteCategory: ", "err", err, "query", query)
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}
