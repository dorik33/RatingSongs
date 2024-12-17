CREATE TABLE song (
                      id SERIAL PRIMARY KEY,
                      title VARCHAR(255) NOT NULL,
                      artist VARCHAR(255) NOT NULL,
                      album VARCHAR(255),
                      release_date DATE,
                      created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE rating (
                        id SERIAL PRIMARY KEY,
                        song_id INT NOT NULL,
                        rating INT CHECK (rating >= 1 AND rating <= 100),
                        description TEXT,
                        created_at TIMESTAMP DEFAULT NOW(),
                        FOREIGN KEY (song_id) REFERENCES song(id) ON DELETE CASCADE
);

INSERT INTO song(title, artist, album, release_date) VALUES ('rave mod', 'cmh', null, '2018-06-10');
INSERT INTO rating(song_id, rating) VALUES (1, 100);

INSERT INTO song(title, artist, album, release_date) VALUES ($1, $2, $3, $4) RETURNING id;

SELECT SELECT id, title, artist, album, release_date FROM song FROM song;

