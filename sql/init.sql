CREATE TABLE user_positions(
                               id SERIAL PRIMARY KEY,
                               name VARCHAR(30) UNIQUE NOT NULL
);

CREATE TABLE users (
                       uid UUID PRIMARY KEY,
                       name VARCHAR(20),
                       surname VARCHAR(20),
                       position_id INT,
                       email VARCHAR(35),
                       phone VARCHAR(10), -- без +7/8
                       hire_date DATE,
                       last_online TIMESTAMP,
                       pfp_url TEXT,
                       is_admin BOOL DEFAULT false, --TODO change to diff service(admin panel)
                       FOREIGN KEY (position_id) REFERENCES user_positions(id) ON DELETE CASCADE
    --TODO more fields
);

CREATE TABLE idea_categories(
                                id SERIAL PRIMARY KEY,
                                name VARCHAR(30) UNIQUE NOT NULL
);


CREATE TABLE idea_statuses(
                              id SERIAL PRIMARY KEY,
                              name VARCHAR(10) UNIQUE NOT NULL --initiated, rejected, approved
);

CREATE TABLE ideas(
                      idea_uid UUID PRIMARY KEY,
                      author UUID NOT NULL,
                      creation_date DATE NOT NULL, --maybe change to TIMESTAMP type
                      status_id INT,
                      category_id INT,
                      like_count INT DEFAULT 0,
                      dislike_count INT DEFAULT 0,
                      FOREIGN KEY (status_id) REFERENCES idea_statuses(id),
                      FOREIGN KEY (category_id) REFERENCES idea_categories(id)
);

CREATE TABLE comments(
                         comment_id SERIAL PRIMARY KEY ,
                         idea_uid UUID NOT NULL,
                         author_id UUID NOT NULL,
                         timestamp TIMESTAMP NOT NULL,
                         comment_text TEXT NOT NULL,
    --TODO create reactions for comments
                         FOREIGN KEY (idea_uid) REFERENCES ideas(idea_uid) ON DELETE CASCADE,
                         FOREIGN KEY (author_id) REFERENCES users(uid)
);

CREATE TABLE replies(
                        comment_id INT NOT NULL,
                        author_id UUID NOT NULL,
                        timestamp TIMESTAMP NOT NULL,
                        reply_text TEXT NOT NULL,
                        FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE,
                        FOREIGN KEY (author_id) REFERENCES users(uid)
);

CREATE TABLE browse_history(
                        visitor_id UUID NOT NULL,
                        idea_id UUID NOT NULL
)