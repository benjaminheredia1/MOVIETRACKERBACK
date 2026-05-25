CREATE TABLE ITEMS (
    id SERIAL PRIMARY KEY,
    tmdb_id INT NOT NULL,
    adult BOOLEAN DEFAULT FALSE,
    backdrop_path TEXT,
    name VARCHAR(255),
    original_name VARCHAR(255),
    overview TEXT,
    poster_path TEXT,
    media_type VARCHAR(10),
    original_language VARCHAR(5),
    popularity DECIMAL(10, 4),
    first_air_date DATE,
    softcore BOOLEAN,
    genre_ids TEXT, -- Almacena los IDs de género como una cadena separada por comas
    origin_country TEXT, -- Almacena los países de origen como una cadena separada por comas
    vote_average DECIMAL(3, 1),
    vote_count INT,
    -- Opciones de seguimiento para el usuario
    list_id INT NULL,  
    status VARCHAR(10) DEFAULT 'pending',
    comentary_user TEXT,
    calification_user FLOAT,
    watched_at TIMESTAMP,
    added_at TIMESTAMP DEFAULT NOW()
);




