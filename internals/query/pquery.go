package query

//USER query
const (
	GetUserByID = `SELECT id, name, token, email, password FROM users
			WHERE id = $1`

	GetUserPasswordByEmail = `SELECT password FROM users
			WHERE email = $1`

	CheckUserByEmail = `SELECT exists(SELECT id FROM users WHERE email=$1)`

	GetAllUser = `SELECT * FROM users`

	AddUser = `INSERT INTO users (name, token, email, password)
				VALUES($1,$2,$3,$4) RETURNING id`

	DeleteUser = `DELETE FROM users WHERE id = $1`
)

//POST query
const (
	GetPostByID = `SELECT id, owner, url, title, time, anons, data FROM posts
			WHERE id = $1`

	GetAllPost = `SELECT * FROM posts
			ORDER BY id DESC`

	Get10Post = `SELECT * FROM posts LIMIT 10 OFFSET $1`

	AddPost = `INSERT INTO posts (owner, url, title, time, anons, data)
				VALUES($1,$2,$3,$4,$5,$6) RETURNING id`

	DeletePost = `DELETE FROM posts WHERE id = $1`
)

//COMMENT query
const (
	GetCommentByID = `SELECT id, owner, url, title, time, anons, data FROM comments
			WHERE id = $1`

	GetCommentByPostID = `SELECT id, owner, url, title, time, anons, data FROM comments
			WHERE postid =$1`

	AddComment = `INSERT INTO comment (postid, owner, name, time, data)
				VALUES($1,$2,$3,$4,$5) RETURNING id`

	DeleteComment = `DELETE FROM comments WHERE id = $1`
)

//TOKEN query
const (
	GetToken = `SELECT token, userid, time FROM tokens
		WHERE token = $1`

	GetTokenID = `SELECT userid FROM tokens 
		WHERE token = $1`

	AddToken = `INSERT INTO tokens (token, userid, time)
		VALUES($1,$2,$3)`

	DeleteToken = `DELETE FROM tokens WHERE token = $1`

	SetTokenUserID = `UPDATE tokens SET userid = $1 
		WHERE token = $2`

	UpdateToken = `UPDATE tokens SET token = $1, time = $2
		WHERE userid = $3`
)

const (
	CreateUserTable = `CREATE TABLE IF NOT EXISTS users (
				id serial PRIMARY KEY,
				name varchar(40) NOT NULL,
				token text,
				email varchar(254) UNIQUE NOT NULL,
				password varchar(256) NOT NULL,
				FOREIGN KEY (token) REFERENCES tokens (token) ON UPDATE CASCADE
				)`
	CreatePostTable = `CREATE TABLE IF NOT EXISTS posts (
				id serial PRIMARY KEY,
				owner integer,
				url text NOT NULL,
				title varchar(256) NOT NULL,
				time text NOT NULL,
				anons varchar(143) NOT NULL,
				data text NOT NULL,
				FOREIGN KEY (owner) REFERENCES users (id) ON UPDATE CASCADE
				)`
	CreateCommentTable = `CREATE TABLE IF NOT EXISTS comments (
				id serial PRIMARY KEY,
				postid integer,
				owner integer,
				name varchar(40),
				time text NOT NULL,
				data text NOT NULL,
				FOREIGN KEY (postid) REFERENCES posts (id),
				FOREIGN KEY (owner) REFERENCES users (id)
				)`
	CreateTokenTable = `CREATE TABLE IF NOT EXISTS tokens (
				token text PRIMARY KEY,
				userid integer UNIQUE NOT NULL,
				time bigint NOT NULL
				)`
)
