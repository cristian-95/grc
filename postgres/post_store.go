package postgres

import (
	"fmt"

	"github.com/cristian-95/grc"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostStore struct {
	*sqlx.DB
}

func (s *PostStore) Post(id uuid.UUID) (grc.Post, error) {
	var p grc.Post
	if err := s.Get(&p, `SELECT * FROM posts WHERE id = $1`, id); err != nil {
		return grc.Post{}, fmt.Errorf("error geting post: %w", err)
	}
	return p, nil
}

func (s *PostStore) PostsByThread(threadID uuid.UUID) ([]grc.Post, error) {
	var pp []grc.Post
	query := `
	SELECT
		posts.*,
		COUNT(comments.*) AS comments_count 
	FROM posts 
	LEFT JOIN comments ON comments.post_id = posts.id
	WHERE thread_id = $1 
	GROUP BY posts.id
	ORDER BY votes DESC`
	if err := s.Select(&pp, query, threadID); err != nil {
		return []grc.Post{}, fmt.Errorf("error geting posts: %w", err)
	}
	return pp, nil
}

func (s *PostStore) Posts() ([]grc.Post, error) {
	var pp []grc.Post
	query := ` 
	SELECT
		posts.*,
		COUNT(comments.*) AS comments_count,
		threads.title AS thread_title
	FROM posts 
	JOIN comments ON comments.post_id = posts.id
	JOIN threads ON threads.id = posts.thread_id
	GROUP BY posts.id, threads.title
	ORDER BY votes DESC`
	if err := s.Select(&pp, query); err != nil {
		return []grc.Post{}, fmt.Errorf("error geting posts: %w", err)
	}
	return pp, nil
}

func (s *PostStore) CreatePost(p *grc.Post) error {
	if err := s.Get(p, `INSERT INTO posts VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		p.ID,
		p.ThreadID,
		p.Title,
		p.Content,
		p.Votes); err != nil {
		return fmt.Errorf("error creating post: %w", err)
	}
	return nil
}

func (s *PostStore) UpdatePost(p *grc.Post) error {
	if err := s.Get(p, `UPDATE posts SET thread_id = $1, title = $2, content = $3, votes = $4 WHERE id = $5 RETURNING *`,
		p.ThreadID,
		p.Title,
		p.Content,
		p.Votes,
		p.ID); err != nil {
		return fmt.Errorf("error updating post: %w", err)
	}
	return nil
}

func (s *PostStore) DeletePost(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM threads WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}
	return nil
}
