CREATE TABLE Users(
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    username text,
    Hash text,
    email text UNIQUE
);

-- @block
CREATE TABLE Posts(
    post_id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_name VARCHAR(255) NOT NULL,
    post_content TEXT,
    post_date time,  -- 'YYYY-MM-DD HH:MM' should be in the brackets 
    user_id INTEGER,
    category_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (category_id) REFERENCES Categories(category_id)
);

-- @block   
CREATE TABLE Comments(
    comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
    comment_content TEXT,
    comment_date time, -- 'YYYY-MM-DD HH:MM' should be in the brackets
    user_id integer
    post_id integer
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (post_id) REFERENCES Posts(post_id)
);

-- @block
CREATE TABLE LikesDislikes(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    like_dislike INT(1),
    FOREIGN KEY (id) REFERENCES Users(user_id),
    FOREIGN KEY (id) REFERENCES Posts(id),
    FOREIGN KEY (id) REFERENCES Comments(id)
);

-- @block    
CREATE TABLE Categories(
    category_id INTEGER PRIMARY KEY AUTOINCREMENT,
    category text
)
--need to add category stuff
