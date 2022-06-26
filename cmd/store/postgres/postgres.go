package store

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/pararti/forump/internals/entity"
	"github.com/pararti/forump/internals/query"
)

type DataBase struct {
	DB *sql.DB
}

func NewDB(config *entity.PSQLConfig) (*DataBase, error) {
	psqlConn := "host=" + config.Host + " port=" + config.Port + " user=" + config.User + " password=" + config.Password + " dbname=" + config.DBName + " sslmode=" + config.SSLMode
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	d := &DataBase{DB: db}
	_, err = d.DB.Exec(`SET client_encoding TO 'utf8'`)
	if err != nil {
		return nil, err
	}
	err = d.CreateTable()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *DataBase) CreateTable() error {
	_, err := d.DB.Exec(query.CreateTokenTable)
	if err != nil {
		return err
	}
	_, err = d.DB.Exec(query.CreateUserTable)
	if err != nil {
		return err
	}
	_, err = d.DB.Exec(query.CreatePostTable)
	if err != nil {
		return err
	}
	_, err = d.DB.Exec(query.CreateCommentTable)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) GetUserByID(id uint32) (*entity.User, error) {
	row := d.DB.QueryRow(query.GetUserByID, id)
	user := &entity.User{}
	err := row.Scan(&user.Id, &user.Name, &user.RefreshToken, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, sql.ErrNoRows
		} else {
			return user, err
		}
	}
	return user, nil
}

func (d *DataBase) GetUserPasswordByEmail(email string) (string, error) {
	var passwd string
	err := d.DB.QueryRow(query.GetUserPasswordByEmail, email).Scan(&passwd)
	if err != nil {
		return "", err
	}
	return passwd, nil
}

func (d *DataBase) AddUser(user *entity.User) (uint32, error) {
	var id uint32
	err := d.DB.QueryRow(query.AddUser, user.Name, user.RefreshToken, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *DataBase) DeleteUser(id uint32) error {
	_, err := d.DB.Exec(query.DeleteUser, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) CheckUserByEmail(email string) (bool, error) {
	var b bool
	err := d.DB.QueryRow(query.CheckUserByEmail, email).Scan(&b)
	if err != nil {
		return false, err
	}
	if b {
		return true, nil
	}
	return false, nil
}

func (d *DataBase) GetPostByID(id uint32) (*entity.Post, error) {
	row := d.DB.QueryRow(query.GetPostByID, id)
	post := &entity.Post{}
	err := row.Scan(&post.Id, &post.Owner, &post.URL, &post.Title, &post.Time, &post.Anons, &post.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return post, sql.ErrNoRows
		} else {
			return post, err
		}
	}
	return post, nil
}

func (d *DataBase) Get10Post(offset int) ([]*entity.Post, error) {
	posts := make([]*entity.Post, 0, 10)
	rows, err := d.DB.Query(query.Get10Post, offset*10)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		p := &entity.Post{}
		err := rows.Scan(&p.Id, &p.Owner, &p.URL, &p.Title, &p.Time, &p.Anons, &p.Data)
		if err != nil {
			return posts, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (d *DataBase) GetAllPost() ([]*entity.Post, error) {
	posts := make([]*entity.Post, 0, 10)
	rows, err := d.DB.Query(query.GetAllPost)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		p := &entity.Post{}
		err := rows.Scan(&p.Id, &p.Owner, &p.URL, &p.Title, &p.Time, &p.Anons, &p.Data)
		if err != nil {
			return posts, err
		}
		posts = append(posts, p)
	}
	return posts, nil

}

func (d *DataBase) AddPost(post *entity.Post) (uint32, error) {
	var id uint32
	post.URL = "/post/"
	post.Time = time.Now().Format("2 Jan 2006 в 15:04")
	if len(post.Data) > 140 {
		post.Anons = post.Data[0:139] + "..."
	} else {
		post.Anons = post.Data
	}
	err := d.DB.QueryRow(query.AddPost, post.Owner, post.URL, post.Title, post.Time, post.Anons, post.Data).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *DataBase) DeletePost(id uint32) error {
	_, err := d.DB.Exec(query.DeletePost, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) GetCommentByID(id uint32) (*entity.Comment, error) {
	row := d.DB.QueryRow(query.GetCommentByID, id)
	comment := &entity.Comment{}
	err := row.Scan(&comment.Id, &comment.PostId, &comment.Owner, &comment.Name, &comment.Time, &comment.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return comment, sql.ErrNoRows
		} else {
			return comment, err
		}
	}
	return comment, nil
}

func (d *DataBase) GetCommentByPostID(id uint32) (*entity.Comment, error) {
	row := d.DB.QueryRow(query.GetCommentByPostID, id)
	comment := &entity.Comment{}
	err := row.Scan(&comment.Id, &comment.PostId, &comment.Owner, &comment.Name, &comment.Time, &comment.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return comment, sql.ErrNoRows
		} else {
			return comment, err
		}
	}
	return comment, nil
}

func (d *DataBase) AddComment(comment *entity.Comment) (uint32, error) {
	var id uint32
	err := d.DB.QueryRow(query.AddComment, comment.Id, comment.PostId, comment.Owner, comment.Name, comment.Time, comment.Data).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *DataBase) DeleteComment(id uint32) error {
	_, err := d.DB.Exec(query.DeleteComment, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) GetToken(tk string) (*entity.Token, error) {
	row := d.DB.QueryRow(query.GetToken, tk)
	token := &entity.Token{}
	err := row.Scan(token.Token, token.UserId, token.Time)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}
	return token, nil
}

func (d *DataBase) AddToken(tk string, uid uint32) error {
	token := &entity.Token{
		Token:  tk,
		UserId: uid,
		Time:   time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	_, err := d.DB.Query(query.AddToken, token.Token, token.UserId, token.Time)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) DeleteToken(token string) error {
	_, err := d.DB.Exec(query.DeleteToken, token)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) UpdateToken(token string, uid uint32) error {
	t := time.Now().Add(24 * 30 * time.Hour).Unix()
	_, err := d.DB.Exec(query.UpdateToken, token, t, uid)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) SetTokenUserID(token string, uid uint32) error {
	_, err := d.DB.Exec(query.SetTokenUserID, uid, token)
	if err != nil {
		return err
	}
	return nil
}
