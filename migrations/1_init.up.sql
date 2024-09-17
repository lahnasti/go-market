CREATE TABLE IF NOT EXISTS users
    (
        id serial PRIMARY KEY,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE
    );
    
CREATE TABLE IF NOT EXISTS products 
    (
        uid serial PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT NOT NULL,
        price double precision NOT NULL,
        delete BOOLEAN NOT NULL DEFAULT false,
        quantity INT NOT NULL
    );

CREATE TABLE IF NOT EXISTS purchases
    (
        uid serial PRIMARY KEY,
        user_id INT NOT NULL,
        product_id INT NOT NULL,
        quantity INT NOT NULL,
        purchase_date TIMESTAMP NOT NULL DEFAULT NOW(),
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (product_id) REFERENCES products(uid)
    );