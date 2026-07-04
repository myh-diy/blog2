package database

import "log"

func InitFTS() {
	if DB == nil {
		return
	}
	DB.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS posts_fts USING fts5(
			title,
			content_md,
			content='posts',
			content_rowid='id'
		)
	`)
	DB.Exec(`
		CREATE TRIGGER IF NOT EXISTS posts_ai AFTER INSERT ON posts BEGIN
			INSERT INTO posts_fts(rowid, title, content_md) VALUES (new.id, new.title, new.content_md);
		END
	`)
	DB.Exec(`
		CREATE TRIGGER IF NOT EXISTS posts_ad AFTER DELETE ON posts BEGIN
			INSERT INTO posts_fts(posts_fts, rowid, title, content_md) VALUES('delete', old.id, old.title, old.content_md);
		END
	`)
	DB.Exec(`
		CREATE TRIGGER IF NOT EXISTS posts_au AFTER UPDATE ON posts BEGIN
			INSERT INTO posts_fts(posts_fts, rowid, title, content_md) VALUES('delete', old.id, old.title, old.content_md);
			INSERT INTO posts_fts(rowid, title, content_md) VALUES (new.id, new.title, new.content_md);
		END
	`)
	log.Println("FTS5 index initialized")
}
