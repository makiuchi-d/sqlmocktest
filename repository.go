package sqlmocktest

import(
	"database/sql"

	"github.com/jmoiron/sqlx"
)

const schema = `
CREATE TABLE user (
  id integer aut_increment,
  name string,
  primary key (id)
)`

type User struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{db}
}

func (r *Repo) Init() error {
	_, err := r.db.Exec(schema)
	return err
}

func (r *Repo) Insert(u *User) (sql.Result, error) {
	return r.db.Exec("insert into user values (?, ?)", u.Id, u.Name)
}

func (r *Repo) GetUser(id int) (*User, error) {
	var u User
	err := r.db.Get(&u, "SELECT * FROM user WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, err
}
