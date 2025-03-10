package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
)

type PostDetails struct {
	ID             int    `db:"id"`
	Title          string `db:"title"`
	Thumbnail      string `db:"thumbnail"`
	Chapter_Number int    `db:"chapter_number"`
	Video          string `db:"video"`
}

func GetChapters(ctx context.Context, conn *pgx.Conn) ([]PostDetails, error) {
	query := `SELECT id, title, thumbnail, chapter_number, video FROM chapters;` //Explicitly list columns

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var posts []PostDetails
	for rows.Next() {
		var post PostDetails
		err := rows.Scan(
			&post.ID, &post.Title, &post.Thumbnail, &post.Chapter_Number, &post.Video,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return posts, nil
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %s\n", handler.Filename)
	fmt.Printf("File Size: %d\n", handler.Size)
	fmt.Printf("MIME Header: %v\n", handler.Header)

	dst, err := os.Create("./uploads/" + handler.Filename)

	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File uploaded successfully"))

}
