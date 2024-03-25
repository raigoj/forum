package main

type Users struct {
	Post_id  int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Posts struct {
	Post_id      int64  `json:"post_id"`
	Post_name    string `json:"post_name"`
	Post_content string `json:"post_content"`
	Post_date    string `json:"post_date"`
	User_id      int64  `json:"user_id"`
	Category_id  int64  `json:"category_id"`
	// Comments      []*Comments
	// LikesDislikes []*LikesDislikes
}

type Comments struct {
	Comment_id      int64  `json:"comment_id"`
	Comment_content string `json:"comment_content"`
	Comment_date    string `json:"comment_date"`
	User_id         int64  `json:"user_id"`
	Post_id         int64  `json:"post_id"`
	// LikesDislikes   []*LikesDislikes
}

type LikesDislikes struct {
	Likedislike_id int64 `json:"likedislike_id"`
	Like_dislike   int64 `json:"like_dislike"`
	User_id        int64 `json:"user_id"`
	Post_id        int64 `json:"post_id"`
	Comment_id     int64 `json:"comment_id"`
}

type Categories struct {
	Category_id int64  `json:"category_id"`
	Category    string `json:"category"`
}

// type PostsComments struct {
// 	post     Posts
// 	posts    []Posts
// 	comment  Comments
// 	comments []Comments
// }
