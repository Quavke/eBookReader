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
    now := time.Now()
    testBooks := []*models.Book{
        {
        Model: gorm.Model{
          ID: uint(1),
          CreatedAt: now,
          UpdatedAt: now,
          DeletedAt: gorm.DeletedAt{},
        },
        Title:    "Test Book 1",
        Content:  "This is a long enough content string to pass validation... 1",
        AuthorID: uint(123),
        },
        {
        Model: gorm.Model{
          ID: uint(2),
          CreatedAt: now,
          UpdatedAt: now,
          DeletedAt: gorm.DeletedAt{},
        },
        Title:    "Test Book 2",
        Content:  "This is a long enough content string to pass validation... 2",
        AuthorID: uint(124),
            },
        }
    


    query := regexp.QuoteMeta(`SELECT * FROM "books"`)
    
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

    books, err := repo.GetAll()

    assert.NoError(t, err)
    assert.Len(t, books, 2)
    

    assert.Equal(t, uint(1), books[0].Model.ID)
    assert.Equal(t, "Test Book 1", books[0].Title)
    assert.Equal(t, "This is a long enough content string to pass validation... 1", books[0].Content)
    assert.Equal(t, uint64(123), books[0].AuthorID)
    

    assert.Equal(t, uint(2), books[1].ID)
    assert.Equal(t, "Test Book 2", books[1].Title)
    assert.Equal(t, "This is a long enough content string to pass validation... 2", books[1].Content)
    assert.Equal(t, uint64(124), books[1].AuthorID)

    assert.NoError(t, mock.ExpectationsWereMet())
}

// func TestBookRepo_GetByID(t *testing.T) {
//     db, mock, err := sqlmock.New()
//     assert.NoError(t, err)
//     defer db.Close()

//     gormDB, err := gorm.Open(postgres.New(postgres.Config{
//         Conn: db,
//     }), &gorm.Config{})
//     assert.NoError(t, err)

//     repo := repositories.NewGormBookRepo(gormDB)

//     now := time.Now().UTC()
//     testBook := models.Book{
//         ID:        1,
//         Title:     "Book get by id 1",
// 				Content:   "Test id 1",
//         CreatedAt: now,
//         UpdatedAt: now,
//         Author: models.Author{
//             Firstname: "Alan",
//             Lastname:  "Doe",
//             Birthday:  now,
// 						CreatedAt: now,
// 						UpdatedAt: now,
//           },
//         }
    


//     query := regexp.QuoteMeta(`SELECT * FROM "books"`)
    
//     row := sqlmock.NewRows([]string{
//         "id", "title", "content", "created_at", "updated_at",
//         "author_firstname", "author_lastname", "author_birthday", "author_created_at", "author_updated_at",
//     })
    
//     row.AddRow(
//         testBook.ID,
//         testBook.Title,
// 				testBook.Content,
//         testBook.CreatedAt,
//         testBook.UpdatedAt,
//         testBook.Author.Firstname,
//         testBook.Author.Lastname,
//         testBook.Author.Birthday,
// 				testBook.Author.CreatedAt,
// 				testBook.Author.UpdatedAt,
//     )
    

//     mock.ExpectQuery(query).WillReturnRows(row)

//     book, err := repo.GetByID(1)

//     assert.NoError(t, err)
    
//     // Проверяем первую книгу
//     assert.Equal(t, 1, book.ID)
//     assert.Equal(t, "Book get by id 1", book.Title)
//     assert.Equal(t, "Alan", book.Author.Firstname)
//     assert.Equal(t, "Doe", book.Author.Lastname)
  

//     assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestBookRepo_Update(t *testing.T) {
//     db, mock, err := sqlmock.New()
//     assert.NoError(t, err)
//     defer db.Close()

//     gormDB, err := gorm.Open(postgres.New(postgres.Config{
//         Conn: db,
//     }), &gorm.Config{})
//     assert.NoError(t, err)

//     repo := repositories.NewGormBookRepo(gormDB)

//     now := time.Now().UTC()
// 		testBooks := []*models.Book{
//         {
//             ID:        1,
//             Title:     "The Go Programming Language",
// 						Content:   "Test id 1",
//             CreatedAt: now,
//             UpdatedAt: now,
//             Author: models.Author{
//                 Firstname: "Alan",
//                 Lastname:  "Doe",
//                 Birthday:  now,
// 								CreatedAt: now,
// 								UpdatedAt: now,
//             },
//         },
//         {
//             ID:        2,
//             Title:     "Clean Code",
// 						Content:   "Test id 2",
//             CreatedAt: now,
//             UpdatedAt: now,
//             Author: models.Author{
//                 Firstname: "Robert",
//                 Lastname:  "Martin",
//                 Birthday:  now,
// 								CreatedAt: now,
// 								UpdatedAt: now,
//             },
//         },
//     }
    
// 		mock.ExpectBegin()

//     query := regexp.QuoteMeta(`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT $2`)
    
//     row := sqlmock.NewRows([]string{
//         "id", "title", "content", "created_at", "updated_at",
//         "author_firstname", "author_lastname", "author_birthday", "author_created_at", "author_updated_at",
//     })
    
//     row.AddRow(
//         testBooks[0].ID,
//         testBooks[0].Title,
// 				testBooks[0].Content,
//         testBooks[0].CreatedAt,
//         testBooks[0].UpdatedAt,
//         testBooks[0].Author.Firstname,
//         testBooks[0].Author.Lastname,
//         testBooks[0].Author.Birthday,
// 				testBooks[0].Author.CreatedAt,
// 				testBooks[0].Author.UpdatedAt,
//     )
	

//     mock.ExpectQuery(query).WithArgs(1, 1).WillReturnRows(row)

// 		updateQuery := regexp.QuoteMeta(
//         `UPDATE "books" SET "title"=$1,"content"=$2,"updated_at"=$3,"author_firstname"=$4,"author_lastname"=$5,"author_birthday"=$6,"author_updated_at"=$7 WHERE "id" = $8`,
//     )
//     mock.ExpectExec(updateQuery).
//         WithArgs(
//             testBooks[1].Title,
//             testBooks[1].Content,
// 						sqlmock.AnyArg(),
//             testBooks[1].Author.Firstname,
// 						testBooks[1].Author.Lastname,
// 						testBooks[1].Author.Birthday,
// 						sqlmock.AnyArg(),
//             1,
//         ).
//         WillReturnResult(sqlmock.NewResult(1, 1))

// 		mock.ExpectCommit()
		
//     err = repo.Update(testBooks[1], 1)
		

//     assert.NoError(t, err)
  

//     assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestBookRepo_Create(t *testing.T) {
//   db, mock, err := sqlmock.New()
//   assert.NoError(t, err)
//   defer db.Close()

//   gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
//   assert.NoError(t, err)

//   repo := repositories.NewGormBookRepo(gormDB)

//   now := time.Now().UTC()
//   testBook := &models.Book{
//     ID: 1,
//     Title: "Test book. Create",
//     Content: "Test test test create",
//     CreatedAt: now,
//     UpdatedAt: now,
//     Author: models.Author{
//       Firstname: "Alexz",
//       Lastname: "Zaa",
//       Birthday: now,
//       CreatedAt: now,
//       UpdatedAt: now,
//     },
//   }
//   mock.ExpectBegin()

//   query := regexp.QuoteMeta(`INSERT INTO "books" ("title","content","created_at","updated_at","author_firstname","author_lastname","author_birthday","author_created_at","author_updated_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)

//   row := sqlmock.NewRows([]string{
//         "id", "title", "content", "created_at", "updated_at",
//         "author_firstname", "author_lastname", "author_birthday", "author_created_at", "author_updated_at",
//     })
//   row.AddRow(
//     testBook.ID,
//     testBook.Title,
//     testBook.Content,
//     testBook.CreatedAt,
//     testBook.UpdatedAt,
//     testBook.Author.Firstname,
//     testBook.Author.Lastname,
//     testBook.Author.Birthday,
//     testBook.Author.CreatedAt,
//     testBook.Author.UpdatedAt,
//   )

//   mock.ExpectQuery(query).WillReturnRows(row)

//   mock.ExpectCommit()

//   err = repo.Create(testBook)

//   assert.NoError(t, err)
//   assert.NoError(t, mock.ExpectationsWereMet())
// }