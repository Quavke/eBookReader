package repositories_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/Quavke/eBookReader/pkg/models"
	"github.com/Quavke/eBookReader/pkg/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Вспомогательная функция для создания mock базы данных
func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)
	
	cleanup := func() {
		db.Close()
	}
	
	return gormDB, mock, cleanup
}

func TestBookRepo_GetAll(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)
	now := time.Now()
	
	testBooks := []*models.Book{
		{
			Model: gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
			Title:    "Test Book 1",
			Content:  "This is a long enough content string to pass validation... 1",
			AuthorID: uint(123),
		},
		{
			Model: gorm.Model{ID: 2, CreatedAt: now, UpdatedAt: now},
			Title:    "Test Book 2",
			Content:  "This is a long enough content string to pass validation... 2",
			AuthorID: uint(124),
		},
	}

	// GORM автоматически добавляет COUNT запрос для пагинации
	countQuery := regexp.QuoteMeta(`SELECT count(*) FROM "books" WHERE "books"."deleted_at" IS NULL`)
	mock.ExpectQuery(countQuery).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	// Тест успешного получения всех книг
	query := regexp.QuoteMeta(`SELECT * FROM "books" WHERE "books"."deleted_at" IS NULL`)
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "content", "author_id",
	})
	for _, book := range testBooks {
		rows.AddRow(
			book.Model.ID,
			book.Model.CreatedAt,
			book.Model.UpdatedAt,
			book.Model.DeletedAt,
			book.Title,
			book.Content,
			book.AuthorID,
		)
	}
	mock.ExpectQuery(query).WillReturnRows(rows)

	pag := &models.Pagination{Limit: 10, Page: 1, Sort: "title"}
	result, err := repo.GetAll(pag)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	books := result.Rows.([]models.Book)
	assert.Len(t, books, 2)
	assert.Equal(t, uint(1), books[0].ID)
	assert.Equal(t, "Test Book 1", books[0].Title)
	assert.Equal(t, "This is a long enough content string to pass validation... 1", books[0].Content)
	assert.Equal(t, uint(123), books[0].AuthorID)
	assert.Equal(t, uint(2), books[1].ID)
	assert.Equal(t, "Test Book 2", books[1].Title)
	assert.Equal(t, "This is a long enough content string to pass validation... 2", books[1].Content)
	assert.Equal(t, uint(124), books[1].AuthorID)
	
	assert.NoError(t, mock.ExpectationsWereMet())

	// Тест с некорректными данными: пустой результат
	mock.ExpectQuery(countQuery).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "content", "author_id",
	}))
	
	pagEmpty := &models.Pagination{Limit: 10, Page: 2, Sort: "title"}
	resultEmpty, err := repo.GetAll(pagEmpty)
	
	assert.Error(t, err)
	assert.Nil(t, resultEmpty)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepo_GetAll_InvalidPagination(t *testing.T) {
	// Тест с некорректными параметрами пагинации
	invalidPag := &models.Pagination{Limit: 0, Page: 0, Sort: ""}
	
	// Проверяем, что используются значения по умолчанию
	assert.Equal(t, uint(10), invalidPag.GetLimit())
	assert.Equal(t, uint(1), invalidPag.GetPage())
	assert.Equal(t, "Id desc", invalidPag.GetSort())
}

func TestBookRepo_GetByID(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)

	now := time.Now().UTC()
	testBook := models.Book{
		Model: gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		Title:     "Book get by id 1",
		Content:   "Test content for book 1",
		AuthorID:  123,
	}

	// Тест успешного получения книги по ID
	query := regexp.QuoteMeta(`SELECT * FROM "books" WHERE id = $1 AND "books"."deleted_at" IS NULL ORDER BY "books"."id" LIMIT $2`)
	
	row := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "content", "author_id",
	})
	
	row.AddRow(
		testBook.ID,
		testBook.CreatedAt,
		testBook.UpdatedAt,
		testBook.DeletedAt,
		testBook.Title,
		testBook.Content,
		testBook.AuthorID,
	)

	mock.ExpectQuery(query).WithArgs(1, 1).WillReturnRows(row)

	book, err := repo.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, uint(1), book.ID)
	assert.Equal(t, "Book get by id 1", book.Title)
	assert.Equal(t, "Test content for book 1", book.Content)
	assert.Equal(t, uint(123), book.AuthorID)

	assert.NoError(t, mock.ExpectationsWereMet())

	// Тест с некорректными данными: книга не найдена
	mock.ExpectQuery(query).WithArgs(999, 1).WillReturnError(gorm.ErrRecordNotFound)

	bookNotFound, err := repo.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, bookNotFound)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepo_GetByID_InvalidID(t *testing.T) {
	gormDB, _, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)

	// Тест с некорректным ID: 0
	book, err := repo.GetByID(0)
	assert.Error(t, err)
	assert.Nil(t, book)
}

func TestBookRepo_Create(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)

	testBook := &models.Book{
		Title:    "Test book. Create",
		Content:  "Test test test create",
		AuthorID: 123,
	}

	// Тест успешного создания книги
	mock.ExpectBegin()
	
	query := regexp.QuoteMeta(`INSERT INTO "books" ("created_at","updated_at","deleted_at","title","content","author_id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)

	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(testBook)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), testBook.ID)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Тест с некорректными данными: пустой заголовок
	invalidBook := &models.Book{
		Title:    "",
		Content:  "Some content",
		AuthorID: 123,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(query).WillReturnError(gorm.ErrInvalidData)
	mock.ExpectRollback()

	err = repo.Create(invalidBook)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepo_Create_InvalidData(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)

	// Тест с некорректными данными: слишком короткий контент
	invalidBook := &models.Book{
		Title:    "Valid Title",
		Content:  "Short", // Меньше 10 символов
		AuthorID: 123,
	}

	mock.ExpectBegin()
	query := regexp.QuoteMeta(`INSERT INTO "books" ("created_at","updated_at","deleted_at","title","content","author_id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)
	mock.ExpectQuery(query).WillReturnError(gorm.ErrInvalidData)
	mock.ExpectRollback()

	err := repo.Create(invalidBook)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepo_Update(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)

	now := time.Now().UTC()
	existingBook := models.Book{
		Model: gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		Title:    "The Go Programming Language",
		Content:  "Original content",
		AuthorID: 123,
	}
	
	updatedBook := &models.Book{
		Title:    "Clean Code",
		Content:  "Updated content",
	}

	// Тест успешного обновления книги
	mock.ExpectBegin()

	query := regexp.QuoteMeta(`SELECT * FROM "books" WHERE id = $1 AND "books"."deleted_at" IS NULL ORDER BY "books"."id" LIMIT $2`)
	
	row := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "content", "author_id",
	})
	
	row.AddRow(
		existingBook.ID,
		existingBook.CreatedAt,
		existingBook.UpdatedAt,
		existingBook.DeletedAt,
		existingBook.Title,
		existingBook.Content,
		existingBook.AuthorID,
	)

	mock.ExpectQuery(query).WithArgs(1, 1).WillReturnRows(row)

	// GORM генерирует UPDATE с другим порядком полей и добавляет deleted_at IS NULL
	updateQuery := regexp.QuoteMeta(`UPDATE "books" SET "updated_at"=$1,"title"=$2,"content"=$3 WHERE "books"."deleted_at" IS NULL AND "id" = $4`)
	mock.ExpectExec(updateQuery).
		WithArgs(
			sqlmock.AnyArg(), // updated_at
			updatedBook.Title,
			updatedBook.Content,
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	
	err := repo.Update(updatedBook, 1)
	
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Тест с некорректными данными: книга не найдена
	mock.ExpectBegin()
	mock.ExpectQuery(query).WithArgs(999, 1).WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	err = repo.Update(updatedBook, 999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepo_Update_InvalidData(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)

	existingBook := models.Book{
		Model: gorm.Model{ID: 1, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		Title:    "Existing Book",
		Content:  "Existing content",
		AuthorID: 123,
	}

	// Тест с некорректными данными: пустой заголовок
	invalidBook := &models.Book{
		Title:    "",
		Content:  "Valid content",
	}

	mock.ExpectBegin()
	query := regexp.QuoteMeta(`SELECT * FROM "books" WHERE id = $1 AND "books"."deleted_at" IS NULL ORDER BY "books"."id" LIMIT $2`)
	
	row := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "content", "author_id",
	})
	
	row.AddRow(
		existingBook.ID,
		existingBook.CreatedAt,
		existingBook.UpdatedAt,
		existingBook.DeletedAt,
		existingBook.Title,
		existingBook.Content,
		existingBook.AuthorID,
	)

	mock.ExpectQuery(query).WithArgs(1, 1).WillReturnRows(row)

	// GORM генерирует UPDATE только для измененных полей
	updateQuery := regexp.QuoteMeta(`UPDATE "books" SET "updated_at"=$1,"content"=$2 WHERE "books"."deleted_at" IS NULL AND "id" = $3`)
	mock.ExpectExec(updateQuery).
		WithArgs(
			sqlmock.AnyArg(), // updated_at
			invalidBook.Content,
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	
	err := repo.Update(invalidBook, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepo_Delete(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)

	// Тест успешного удаления книги
	// GORM использует soft delete - UPDATE вместо DELETE и может добавлять транзакции
	mock.ExpectBegin()
	query := regexp.QuoteMeta(`UPDATE "books" SET "deleted_at"=$1 WHERE id = $2 AND "books"."deleted_at" IS NULL`)
	mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Тест с некорректными данными: книга не найдена
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), 999).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = repo.Delete(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no book found with id 999")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepo_Delete_InvalidID(t *testing.T) {
	gormDB, _, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)

	// Тест с некорректным ID: 0
	err := repo.Delete(0)
	assert.Error(t, err)
}

func TestBookRepo_IsBelongsTo(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)

	now := time.Now().UTC()
	testBook := models.Book{
		Model: gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		Title:    "Test Book",
		Content:  "Test content",
		AuthorID: 123,
	}

	// Тест успешной проверки принадлежности книги автору
	query := regexp.QuoteMeta(`SELECT * FROM "books" WHERE (id = $1 AND author_id = $2) AND "books"."deleted_at" IS NULL ORDER BY "books"."id" LIMIT $3`)
	
	row := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "content", "author_id",
	})
	
	row.AddRow(
		testBook.ID,
		testBook.CreatedAt,
		testBook.UpdatedAt,
		testBook.DeletedAt,
		testBook.Title,
		testBook.Content,
		testBook.AuthorID,
	)

	mock.ExpectQuery(query).WithArgs(1, 123, 1).WillReturnRows(row)

	belongs, err := repo.IsBelongsTo(1, 123)
	
	assert.NoError(t, err)
	assert.True(t, belongs)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Тест с некорректными данными: книга не принадлежит автору
	// Когда книга не найдена, GORM возвращает ошибку, а не пустой результат
	mock.ExpectQuery(query).WithArgs(1, 999, 1).WillReturnError(gorm.ErrRecordNotFound)

	belongs, err = repo.IsBelongsTo(1, 999)
	
	assert.Error(t, err)
	assert.False(t, belongs)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepo_IsBelongsTo_InvalidData(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)

	// Тест с некорректными данными: несуществующий ID книги
	query := regexp.QuoteMeta(`SELECT * FROM "books" WHERE (id = $1 AND author_id = $2) AND "books"."deleted_at" IS NULL ORDER BY "books"."id" LIMIT $3`)
	mock.ExpectQuery(query).WithArgs(999, 123, 1).WillReturnError(gorm.ErrRecordNotFound)

	belongs, err := repo.IsBelongsTo(999, 123)
	
	assert.Error(t, err)
	assert.False(t, belongs)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тесты граничных случаев
func TestBookRepo_EdgeCases(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)

	// Тест с максимальными значениями
	maxBook := &models.Book{
		Title:    string(make([]byte, 400)), // Максимальная длина заголовка
		Content:  string(make([]byte, 1000)), // Длинный контент
		AuthorID: 999999999,
	}

	mock.ExpectBegin()
	query := regexp.QuoteMeta(`INSERT INTO "books" ("created_at","updated_at","deleted_at","title","content","author_id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)
	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(maxBook)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Тест с минимальными значениями
	minBook := &models.Book{
		Title:    "A", // Минимальная длина заголовка
		Content:  "1234567890", // Минимальная длина контента
		AuthorID: 1,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	mock.ExpectCommit()

	err = repo.Create(minBook)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тесты производительности
func TestBookRepo_Performance(t *testing.T) {
	gormDB, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewGormBookRepo(gormDB)

	// Тест с большим количеством книг
	largePagination := &models.Pagination{Limit: 1000, Page: 1, Sort: "id"}
	
	// GORM автоматически добавляет COUNT запрос для пагинации
	countQuery := regexp.QuoteMeta(`SELECT count(*) FROM "books" WHERE "books"."deleted_at" IS NULL`)
	mock.ExpectQuery(countQuery).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1000))
	
	query := regexp.QuoteMeta(`SELECT * FROM "books" WHERE "books"."deleted_at" IS NULL`)
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "content", "author_id",
	})
	
	// Создаем 1000 тестовых записей
	for i := 1; i <= 1000; i++ {
		rows.AddRow(
			i,
			time.Now(),
			time.Now(),
			nil,
			fmt.Sprintf("Book %d", i),
			fmt.Sprintf("Content for book %d", i),
			uint(i%100+1),
		)
	}
	
	mock.ExpectQuery(query).WillReturnRows(rows)

	result, err := repo.GetAll(largePagination)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Rows.([]models.Book), 1000)
	assert.NoError(t, mock.ExpectationsWereMet())
}