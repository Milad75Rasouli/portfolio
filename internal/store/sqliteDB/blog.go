package sqlitedb

import (
	"context"
	"time"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
	"github.com/Milad75Rasouli/portfolio/internal/model"
	"github.com/Milad75Rasouli/portfolio/internal/store"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type BlogSqlite struct {
	dbPool *sqlitex.Pool
	logger *zap.Logger
}

func NewBlogSqlite(dbPool *sqlitex.Pool, logger *zap.Logger) *BlogSqlite {
	return &BlogSqlite{
		dbPool: dbPool,
		logger: logger,
	}
}

func (b BlogSqlite) parseToBlog(stmt *sqlite.Stmt) (model.Blog, error) {
	var (
		blog model.Blog
		err  error
	)
	blog.ID = stmt.GetInt64("id")
	blog.Title = stmt.GetText("title")
	blog.Body = stmt.GetText("body")
	blog.Caption = stmt.GetText("caption")
	blog.ImagePath = stmt.GetText("image_path")
	blog.CreatedAt, err = time.Parse(timeLayout, stmt.GetText("created_at"))
	if err != nil {
		return blog, err
	}
	blog.ModifiedAt, err = time.Parse(timeLayout, stmt.GetText("modified_at"))

	return blog, err
}

func (b BlogSqlite) parseToBlogWithCategory(stmt *sqlite.Stmt) (model.BlogWithCategory, error) {
	var (
		blog     model.BlogWithCategory
		category model.Category
		err      error
	)
	blog.Blog.ID = stmt.GetInt64("id")
	blog.Blog.Title = stmt.GetText("title")
	blog.Blog.Body = stmt.GetText("body")
	blog.Blog.Caption = stmt.GetText("caption")
	blog.Blog.ImagePath = stmt.GetText("image_path")
	category.ID = stmt.GetInt64("category_id")
	category.Name = stmt.GetText("category")
	blog.Category = append(blog.Category, category)
	blog.Blog.CreatedAt, err = time.Parse(timeLayout, stmt.GetText("created_at"))

	if err != nil {
		return blog, err
	}
	blog.Blog.ModifiedAt, err = time.Parse(timeLayout, stmt.GetText("modified_at"))

	return blog, err
}

func (b BlogSqlite) parseToCategory(stmt *sqlite.Stmt) (model.Category, error) {
	var (
		blog model.Category
		err  error
	)
	blog.ID = stmt.GetInt64("id")
	blog.Name = stmt.GetText("name")
	return blog, err
}

func (b BlogSqlite) parseToCategoryRelation(stmt *sqlite.Stmt) (model.Relation, error) {
	var (
		blog model.Relation
		err  error
	)
	blog.PostID = stmt.GetInt64("post_id")
	blog.CategoryID = stmt.GetInt64("category_id")
	return blog, err
}

func (b *BlogSqlite) CreateBlog(ctx context.Context, blog model.Blog) (int64, error) {
	var rowID int64
	conn := b.dbPool.Get(ctx)
	defer b.dbPool.Put(conn)

	stmt, err := conn.Prepare(`INSERT INTO post (title, body, caption, image_path,created_at,modified_at)
	VALUES ($1, $2, $3, $4, $5, $6);`)
	if err != nil {
		return rowID, errors.Errorf("unable to create the new blog %s", err.Error())
	}
	defer stmt.Finalize()
	stmtSelect, err := conn.Prepare(`SELECT last_insert_rowid();`)
	if err != nil {
		return 0, errors.Errorf("unable to prepare the select statement: %s", err.Error())
	}
	defer stmtSelect.Finalize()
	stmt.SetText("$1", blog.Title)
	stmt.SetText("$2", blog.Body)
	stmt.SetText("$3", blog.Caption)
	stmt.SetText("$4", blog.ImagePath)
	stmt.SetText("$5", blog.CreatedAt.Format(timeLayout))
	stmt.SetText("$6", blog.ModifiedAt.Format(timeLayout))

	_, err = stmt.Step()
	if err != nil {
		return rowID, err
	}
	if conn.Changes() == 0 {
		return rowID, store.BlogCreateError
	}
	hasRow, err := stmtSelect.Step()
	if err != nil {
		return rowID, err
	}

	if hasRow {
		rowID = conn.LastInsertRowID()
	}
	return rowID, err
}

func (b *BlogSqlite) GetBlogByID(ctx context.Context, id int64) (model.Blog, error) {
	var blog model.Blog
	conn := b.dbPool.Get(ctx)
	defer b.dbPool.Put(conn)
	stmt, err := conn.Prepare(`SELECT * FROM post WHERE id=$1 LIMIT 1;`)
	if err != nil {
		return blog, errors.Errorf("unable to get the post %s from id", err.Error())
	}
	defer stmt.Finalize()
	stmt.SetInt64("$1", id)

	var hasRow bool
	hasRow, err = stmt.Step()
	if hasRow == false {
		return blog, store.BlogNotFoundError
	}
	if err != nil {
		return blog, err
	}
	blog, err = b.parseToBlog(stmt)
	return blog, err
}

func (b *BlogSqlite) GetAllBlog(ctx context.Context) ([]model.Blog, error) {
	var blog []model.Blog
	conn := b.dbPool.Get(ctx)
	defer b.dbPool.Put(conn)
	stmt, err := conn.Prepare(`SELECT * FROM post;`)
	if err != nil {
		return blog, errors.Errorf("unable to get all blog %s", err.Error())
	}
	defer stmt.Finalize()

	for {
		var (
			swapBlog model.Blog
			hasRow   bool
		)
		hasRow, err = stmt.Step()
		if hasRow == false {
			break
		}
		swapBlog, err = b.parseToBlog(stmt)
		if err != nil {
			return blog, errors.Errorf("getting the blog from database error %s", err.Error())
		}
		blog = append(blog, swapBlog)
	}
	return blog, err
}

func (u *BlogSqlite) DeleteBlogByID(ctx context.Context, id int64) error {
	conn := u.dbPool.Get(ctx)
	defer u.dbPool.Put(conn)
	stmt, err := conn.Prepare(`DELETE FROM post WHERE id=$1;`)
	if err != nil {
		return err
	}
	defer stmt.Finalize()
	stmt.SetInt64("$1", id)
	_, err = stmt.Step()
	return err
}

func (u *BlogSqlite) UpdateBlogByID(ctx context.Context, blog model.Blog) error {
	conn := u.dbPool.Get(ctx)
	defer u.dbPool.Put(conn)
	var s string
	s = `UPDATE post
	SET title=$1, body=$2, caption=$3, image_path=$4
	WHERE id=$5;`
	stmt, err := conn.Prepare(s)
	if err != nil {
		return err
	}
	stmt.SetText("$1", blog.Title)
	stmt.SetText("$2", blog.Body)
	stmt.SetText("$3", blog.Caption)
	stmt.SetText("$4", blog.ImagePath)
	stmt.SetInt64("$5", blog.ID)
	defer stmt.Finalize()
	var hasRow bool
	hasRow, err = stmt.Step()
	if hasRow {
		return store.BlogNotFoundError
	}
	return err
}

/******************* Category *******************/

func (b *BlogSqlite) CreateCategory(ctx context.Context, category model.Category) (int64, error) {
	var rowID int64
	conn := b.dbPool.Get(ctx)
	defer b.dbPool.Put(conn)

	stmt, err := conn.Prepare(`INSERT INTO category (name) VALUES ($1);`)
	if err != nil {
		return rowID, errors.Errorf("unable to create the new category %s", err.Error())
	}
	defer stmt.Finalize()
	stmtSelect, err := conn.Prepare(`SELECT last_insert_rowid();`)
	if err != nil {
		return 0, errors.Errorf("unable to prepare the select statement: %s", err.Error())
	}
	defer stmtSelect.Finalize()
	stmt.SetText("$1", category.Name)

	_, err = stmt.Step()
	if err != nil {
		return rowID, err
	}
	if conn.Changes() == 0 {
		return rowID, store.CategoryCreateError
	}
	hasRow, err := stmtSelect.Step()
	if err != nil {
		return rowID, err
	}

	if hasRow {
		rowID = conn.LastInsertRowID()
	}
	return rowID, err
}

func (b *BlogSqlite) GetCategoryByID(ctx context.Context, id int64) (model.Category, error) {
	var blog model.Category
	conn := b.dbPool.Get(ctx)
	defer b.dbPool.Put(conn)
	stmt, err := conn.Prepare(`SELECT * FROM category WHERE id=$1 LIMIT 1;`)
	if err != nil {
		return blog, errors.Errorf("unable to get the category %s from id", err.Error())
	}
	defer stmt.Finalize()
	stmt.SetInt64("$1", id)

	var hasRow bool
	hasRow, err = stmt.Step()
	if hasRow == false {
		return blog, store.CategoryNotFoundError
	}
	if err != nil {
		return blog, err
	}
	blog, err = b.parseToCategory(stmt)
	return blog, err
}

func (b *BlogSqlite) GetAllCategory(ctx context.Context) ([]model.Category, error) {
	var blog []model.Category
	conn := b.dbPool.Get(ctx)
	defer b.dbPool.Put(conn)
	stmt, err := conn.Prepare(`SELECT * FROM category;`)
	if err != nil {
		return blog, errors.Errorf("unable to get all category %s", err.Error())
	}
	defer stmt.Finalize()

	for {
		var (
			swapCategory model.Category
			hasRow       bool
		)
		hasRow, err = stmt.Step()
		if hasRow == false {
			break
		}
		swapCategory, err = b.parseToCategory(stmt)
		if err != nil {
			return blog, errors.Errorf("getting the category from database error %s", err.Error())
		}
		blog = append(blog, swapCategory)
	}
	return blog, err
}

func (u *BlogSqlite) DeleteCategoryByID(ctx context.Context, id int64) error {
	conn := u.dbPool.Get(ctx)
	defer u.dbPool.Put(conn)
	stmt, err := conn.Prepare(`DELETE FROM category WHERE id=$1;`)
	if err != nil {
		return err
	}
	defer stmt.Finalize()
	stmt.SetInt64("$1", id)
	_, err = stmt.Step()
	return err
}

func (u *BlogSqlite) UpdateCategoryByID(ctx context.Context, blog model.Category) error {
	conn := u.dbPool.Get(ctx)
	defer u.dbPool.Put(conn)
	var s string
	s = `UPDATE category SET name=$1 WHERE id=$2;`
	stmt, err := conn.Prepare(s)
	if err != nil {
		return err
	}
	stmt.SetText("$1", blog.Name)
	stmt.SetInt64("$2", blog.ID)
	var hasRow bool
	hasRow, err = stmt.Step()
	if hasRow {
		return store.BlogNotFoundError
	}
	return err
}

/******************* Post & Category Relations *******************/

func (b *BlogSqlite) CreateCategoryRelation(ctx context.Context, Relation model.Relation) error {
	conn := b.dbPool.Get(ctx)
	defer b.dbPool.Put(conn)

	stmt, err := conn.Prepare(`INSERT INTO post_category_relation 
		(category_id, post_id) VALUES ($1,$2);`)
	if err != nil {
		return errors.Errorf("unable to create the new category relation %s", err.Error())
	}
	defer stmt.Finalize()

	stmt.SetInt64("$1", Relation.CategoryID)
	stmt.SetInt64("$2", Relation.PostID)

	_, err = stmt.Step()
	if err != nil {
		return err
	}
	if conn.Changes() == 0 {
		return store.CategoryRelationCreateError
	}

	return err
}
func (b *BlogSqlite) GetCategoryRelationAllByPostID(ctx context.Context, id int64) ([]model.Relation, error) {

	var relation []model.Relation
	conn := b.dbPool.Get(ctx)
	defer b.dbPool.Put(conn)
	stmt, err := conn.Prepare(`SELECT * FROM post_category_relation where post_id=$1;`)
	if err != nil {
		return relation, errors.Errorf("unable to get all category relation %s", err.Error())
	}
	defer stmt.Finalize()
	stmt.SetInt64("$1", id)
	var times int
	for {
		var (
			swapCategoryRelation model.Relation
			hasRow               bool
		)
		hasRow, err = stmt.Step()
		if hasRow == false {
			break
		}
		times++
		swapCategoryRelation, err = b.parseToCategoryRelation(stmt)
		if err != nil {
			return relation, errors.Errorf("getting the category from database error %s", err.Error())
		}
		relation = append(relation, swapCategoryRelation)
	}

	if times == 0 {
		return relation, store.CategoryRelationNotFoundError
	}
	return relation, err
}

func (b *BlogSqlite) GetCategoryRelationAllByCategoryID(ctx context.Context, id int64) ([]model.Relation, error) {
	var relation []model.Relation
	conn := b.dbPool.Get(ctx)
	defer b.dbPool.Put(conn)
	stmt, err := conn.Prepare(`SELECT * FROM post_category_relation where category_id=$1;`)
	if err != nil {
		return relation, errors.Errorf("unable to get all category relation %s", err.Error())
	}
	defer stmt.Finalize()
	stmt.SetInt64("$1", id)

	var times int
	for {
		var (
			swapCategoryRelation model.Relation
			hasRow               bool
		)
		hasRow, err = stmt.Step()
		if hasRow == false {
			break
		}
		times++
		swapCategoryRelation, err = b.parseToCategoryRelation(stmt)
		if err != nil {
			return relation, errors.Errorf("getting the category from database error %s", err.Error())
		}
		relation = append(relation, swapCategoryRelation)
	}
	if times == 0 {
		return relation, store.CategoryRelationNotFoundError
	}
	return relation, err
}

func (b *BlogSqlite) DeleteCategoryRelationAllByPostID(ctx context.Context, id int64) error {
	conn := b.dbPool.Get(ctx)
	defer b.dbPool.Put(conn)
	stmt, err := conn.Prepare(`DELETE FROM post_category_relation WHERE post_id=$1;`)
	if err != nil {
		return err
	}
	defer stmt.Finalize()
	stmt.SetInt64("$1", id)
	var hasRow bool
	hasRow, err = stmt.Step()
	if hasRow == true {
		return store.CategoryNotFoundError
	}
	return err
}

func (b *BlogSqlite) DeleteCategoryRelationAllByCategoryID(ctx context.Context, id int64) error {
	conn := b.dbPool.Get(ctx)
	defer b.dbPool.Put(conn)
	stmt, err := conn.Prepare(`DELETE FROM post_category_relation WHERE category_id=$1;`)
	if err != nil {
		return err
	}
	defer stmt.Finalize()
	stmt.SetInt64("$1", id)
	var hasRow bool
	hasRow, err = stmt.Step()
	if hasRow == true {
		return store.CategoryNotFoundError
	}
	return err
}

/******************* practical part *******************/

// This won't return blog.body
func (b *BlogSqlite) GetAllPostsWithCategory(ctx context.Context) ([]model.BlogWithCategory, error) {
	var (
		blogWithCategory []model.BlogWithCategory
	)
	conn := b.dbPool.Get(ctx)
	defer b.dbPool.Put(conn)
	stmt, err := conn.Prepare(`SELECT p.id as id,p.title as title, 
	p.caption as caption,p.image_path as image_path,
	p.created_at as created_at, p.modified_at as modified_at ,
	c.id as category_id , c.name as category FROM post as p
	LEFT JOIN post_category_relation as pc ON pc.post_id = p.id
	LEFT JOIN category as c ON pc.category_id = c.id;`)
	if err != nil {
		return blogWithCategory, errors.Errorf("unable to get all blogWithCategory %s", err.Error())
	}
	defer stmt.Finalize()

	temp := make(map[model.Blog][]model.Category)
	for {
		var (
			swapBlog model.BlogWithCategory
			hasRow   bool
		)
		hasRow, err = stmt.Step()
		if hasRow == false {
			break // TODO: return a notFoundError in case of having no item
		}
		swapBlog, err = b.parseToBlogWithCategory(stmt)
		if err != nil {
			return blogWithCategory, errors.Errorf("getting the blogWithCategory from database error %s", err.Error())
		}
		category, ok := temp[swapBlog.Blog]
		if ok == false {
			temp[swapBlog.Blog] = swapBlog.Category
		} else {
			category = append(category, swapBlog.Category...)
			temp[swapBlog.Blog] = category
		}
	}
	for b, c := range temp {
		blogWithCategory = append(blogWithCategory, model.BlogWithCategory{Blog: b, Category: c})
	}
	return blogWithCategory, err
}
