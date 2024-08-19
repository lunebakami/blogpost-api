package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lunebakami/futtodos-api/config"
	"github.com/lunebakami/futtodos-api/models"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type PostStorage struct {
	db *sql.DB
}

func NewPostStorage(cfg *config.Config) (*PostStorage, error) {
	url := fmt.Sprintf("%s?authToken=%s", cfg.TursoURL, cfg.TursoToken)

	db, err := sql.Open("libsql", url)

	if err != nil {
		return nil, fmt.Errorf("failed to open db %s: %w", url, err)
	}

	storage := &PostStorage{db: db}

	err = storage.initTable()

	if err != nil {
		db.Close()
		return nil, err
	}

	return storage, nil
}

func (s *PostStorage) initTable() error {
	_, err := s.db.ExecContext(context.Background(), `
      CREATE TABLE IF NOT EXISTS posts (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          title TEXT NOT NULL,
          content TEXT NOT NULL
      )
  `)

	return err
}

func (s *PostStorage) Create(post *models.BlogPost) error {
	result, err := s.db.ExecContext(context.Background(),
		"INSERT INTO posts (title, content) VALUES (?, ?)",
		post.Title,
		post.Content,
	)

	if err != nil {
		return err
	}

	post.ID, err = result.LastInsertId()

	return err
}

func (s *PostStorage) GetAll() ([]models.BlogPost, error) {
	rows, err := s.db.QueryContext(context.Background(),
		"SELECT * FROM posts",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.BlogPost

	for rows.Next() {
		var post models.BlogPost
		err := rows.Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func (s *PostStorage) GetByID(id int64) (models.BlogPost, error) {
	var post models.BlogPost
	err := s.db.QueryRowContext(
    context.Background(), 
    "SELECT id, title, content FROM posts WHERE id = ?", 
    id,
  ).Scan(&post.ID, &post.Title, &post.Content)

  if err != nil {
    if err != sql.ErrNoRows {
      return models.BlogPost{}, errors.New("post not found")
    }
    return models.BlogPost{}, err
  }

  return post, nil
}

func (s *PostStorage) Update(id int64, updatedPost *models.BlogPost) error {
  result, err := s.db.ExecContext(
    context.Background(),
    "UPDATE posts SET title = ?, content = ? WHERE id = ?",
    updatedPost.Title, updatedPost.Content, updatedPost.ID,
  )

  if err != nil {
    return err
  }

  rowsAffected, err := result.RowsAffected()
  if err != nil {
    return err
  }

  if rowsAffected == 0 {
    return errors.New("post not found")
  }

  updatedPost.ID = id

  return nil
}

func (s *PostStorage) Delete(id int64) error {
  result, err := s.db.ExecContext(
    context.Background(),
    "DELETE FROM posts WHERE id = ?",
    id,
  )

  if err != nil {
    return err
  }

  rowsAffected, err := result.RowsAffected()
  if err != nil {
    return err
  }

  if rowsAffected == 0 {
    return errors.New("post not found")
  }

  return nil
}

func (s *PostStorage) Close() error {
  return s.db.Close()
}
