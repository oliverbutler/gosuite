-- Create authors table
CREATE TABLE IF NOT EXISTS authors (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  bio TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample authors
INSERT INTO
  authors (name, bio)
VALUES
  ('John Doe', 'A passionate writer and blogger.'),
  (
    'Jane Smith',
    'Loves to write about technology and innovation.'
  ),
  (
    'Emily Johnson',
    'Freelance writer and nature lover.'
  );

-- Create posts table
CREATE TABLE IF NOT EXISTS posts (
  id INT AUTO_INCREMENT PRIMARY KEY,
  author_id INT,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (author_id) REFERENCES authors (id) ON DELETE CASCADE
);

-- Insert sample posts
INSERT INTO
  posts (author_id, title, content)
VALUES
  (
    1,
    'The Joy of Writing',
    'Writing is a journey of discovery...'
  ),
  (
    2,
    'Tech Trends 2024',
    'This year''s technology trends are all about...'
  ),
  (
    1,
    'A Writer''s Best Friend',
    'For many writers, isolation is a fact of life...'
  );

-- Create comments table
CREATE TABLE IF NOT EXISTS comments (
  id INT AUTO_INCREMENT PRIMARY KEY,
  post_id INT,
  author_name VARCHAR(255),
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);

-- Insert sample comments
INSERT INTO
  comments (post_id, author_name, content)
VALUES
  (1, 'Alice', 'Absolutely love this!'),
  (
    2,
    'Bob',
    'Very insightful, looking forward to more.'
  ),
  (
    1,
    'Charlie',
    'Great perspective, thanks for sharing.'
  );
