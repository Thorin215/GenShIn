CREATE TABLE sentences (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sentence TEXT NOT NULL,
    label BOOLEAN NOT NULL,
    dataset_id VARCHAR(255) NOT NULL,
    sentence_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
