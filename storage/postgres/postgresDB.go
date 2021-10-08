package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"src/http/pkg/model"
)

type PostgresDB struct {
	Pdb *sql.DB
	mu  sync.Mutex
}

func NewPDB(host string, port string, user string, psw string, dbname string, ssl string) (*PostgresDB, error) {
	connStr := "host=" + host + " port=" + port + " user=" + user + " password=" + psw + " dbname=" + dbname + " sslmode=" + ssl

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database %w\n", err)
	}

	database := &PostgresDB{Pdb: db}

	database.Pdb.Exec("CREATE TABLE posts (\n    id VARCHAR(40) PRIMARY KEY NOT NULL,\n    author VARCHAR(50) NOT NULL,\n    message VARCHAR(150) NOT NULL\n);")

	return database, nil
}

func (pdb *PostgresDB) Create(p model.Post) (model.Post, error) {
	pdb.mu.Lock()
	defer pdb.mu.Unlock()

	id := uuid.New()
	idStr := id.String()

	_, err := pdb.Pdb.Exec(
		"INSERT INTO posts (id, author, message) VALUES ($1, $2, $3)", idStr, p.Author, p.Message)
	if err != nil {
		return model.Post{}, errors.New("couldn't create post in database")
	}

	p.Id = id

	return p, nil
}

func (pdb *PostgresDB) Get(id uuid.UUID) (model.Post, error) {
	pdb.mu.Lock()
	defer pdb.mu.Unlock()

	idStr := id.String()

	var p model.Post

	err := pdb.Pdb.QueryRow(
		`SELECT author, message FROM posts WHERE id=$1`, idStr).Scan(&p.Author, &p.Message)
	if err != nil {
		return model.Post{}, errors.New("couldn't find post")
	}

	p.Id = id

	return p, nil
}

func (pdb *PostgresDB) GetAll() []model.Post {
	pdb.mu.Lock()
	defer pdb.mu.Unlock()

	rows, err := pdb.Pdb.Query(
		`SELECT * FROM posts`)
	if err != nil {
		return []model.Post{}
	}
	defer rows.Close()

	posts := []model.Post{}

	for rows.Next() {
		p := model.Post{}
		var pp string
		err := rows.Scan(&pp, &p.Author, &p.Message)
		if err != nil {
			return []model.Post{}
		}
		p.Id, err = uuid.Parse(pp)
		if err != nil {
			return []model.Post{}
		}
		posts = append(posts, p)
	}

	return posts
}

func (pdb *PostgresDB) Update(p model.Post) (model.Post, error) {
	pdb.mu.Lock()
	defer pdb.mu.Unlock()

	idStr := p.Id.String()

	var oldPost model.Post

	err := pdb.Pdb.QueryRow(
		`SELECT author, message FROM posts WHERE id=$1`, idStr).Scan(&oldPost.Author, &oldPost.Message)
	if err != nil {
		return model.Post{}, errors.New("couldn't find post")
	}

	if p.Author == "" {
		p.Author = oldPost.Author
	}
	if p.Message == "" {
		p.Message = oldPost.Message
	}

	_, err = pdb.Pdb.Exec(
		`UPDATE posts SET author=$1, message=$2 WHERE id=$3`, p.Author, p.Message, idStr)
	if err != nil {
		return model.Post{}, errors.New("couldn't update post")
	}

	return p, nil
}

func (pdb *PostgresDB) Delete(id uuid.UUID) error {
	pdb.mu.Lock()
	defer pdb.mu.Unlock()

	idStr := id.String()

	_, err := pdb.Pdb.Exec(
		`DELETE FROM posts where id = $1`, idStr)
	if err != nil {
		return errors.New("couldn't delete post")
	}

	return nil
}

func (pdb *PostgresDB) CreateFromFile(p model.Post) error {
	pdb.mu.Lock()
	defer pdb.mu.Unlock()

	idStr := p.Id.String()

	_, err := pdb.Pdb.Exec(
		`INSERT INTO posts (id, author, message) VALUES ($1, $2, $3)`, idStr, p.Author, p.Message)
	if err != nil {
		return errors.New("couldn't create post")
	}

	return nil
}
