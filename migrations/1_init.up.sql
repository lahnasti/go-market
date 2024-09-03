CREATE TABLE IF NOT EXISTS products 
    (
        uid serial PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT NOT NULL,
        price double precision NOT NULL,
        delete BOOLEAN NOT NULL DEFAULT false
    );