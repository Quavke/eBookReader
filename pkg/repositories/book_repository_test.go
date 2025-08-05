package repositories_test

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBookRepo_GetAll(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: db,
    }), &gorm.Config{})
    assert.NoError(t, err)

    repo := repositories.NewGormBookRepo(gormDB)

    now := time.Now().UTC()
    testBooks := []*models.Book{
        {
            ID:        1,
            Title:     "The Go Programming Language",
						Content:   "Test id 1",
            CreatedAt: now,
            UpdatedAt: now,
            Author: models.Author{
                Firstname: "Alan",
                Lastname:  "Doe",
                Birthday:  now,
								CreatedAt: now,
								UpdatedAt: now,
            },
        },
        {
            ID:        2,
            Title:     "Clean Code",
						Content:   "Test id 2",
            CreatedAt: now,
            UpdatedAt: now,
            Author: models.Author{
                Firstname: "Robert",
                Lastname:  "Martin",
                Birthday:  now,
								CreatedAt: now,
								UpdatedAt: now,
            },
        },
    }

    // Правильный запрос для GetAll
    query := regexp.QuoteMeta(`SELECT * FROM "books"`)
    
    // Создаем строки с данными (все поля в одной таблице)
    rows := sqlmock.NewRows([]string{
        "id", "title", "content", "created_at", "updated_at",
        "author_firstname", "author_lastname", "author_birthday", "author_created_at", "author_updated_at",
    })
    
    // Добавляем тестовые данные
    for _, book := range testBooks {
        rows.AddRow(
            book.ID,
            book.Title,
						book.Content,
            book.CreatedAt,
            book.UpdatedAt,
            book.Author.Firstname,
            book.Author.Lastname,
            book.Author.Birthday,
						book.Author.CreatedAt,
						book.Author.UpdatedAt,
        )
    }

    mock.ExpectQuery(query).WillReturnRows(rows)

    books, err := repo.GetAll()

    assert.NoError(t, err)
    assert.Len(t, books, 2)
    
    // Проверяем первую книгу
    assert.Equal(t, 1, books[0].ID)
    assert.Equal(t, "The Go Programming Language", books[0].Title)
    assert.Equal(t, "Alan", books[0].Author.Firstname)
    assert.Equal(t, "Doe", books[0].Author.Lastname)
    
    // Проверяем вторую книгу
    assert.Equal(t, 2, books[1].ID)
    assert.Equal(t, "Clean Code", books[1].Title)
    assert.Equal(t, "Robert", books[1].Author.Firstname)
    assert.Equal(t, "Martin", books[1].Author.Lastname)

    assert.NoError(t, mock.ExpectationsWereMet())
}