BEGIN;

-- ENUMs
CREATE TYPE role_type AS ENUM ('admin', 'sub-admin', 'user');
CREATE TYPE address_label AS ENUM ('home', 'office', 'gym', 'other');
CREATE TYPE food_type AS ENUM ('veg', 'non-veg');
CREATE TYPE course_category AS ENUM ('starter', 'main-course', 'dessert', 'breakfast', 'beverage');
CREATE TYPE order_status AS ENUM ('placed', 'accepted', 'preparing', 'out_for_delivery', 'delivered', 'cancelled');
CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed');

-- Users Table
CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       name TEXT NOT NULL,
                       email TEXT UNIQUE NOT NULL,
                       password TEXT NOT NULL
);

-- User Roles Table
CREATE TABLE user_roles (
                            id BIGSERIAL PRIMARY KEY,
                            user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                            role role_type NOT NULL,
                            UNIQUE(user_id, role)
);

-- Addresses Table
CREATE TABLE addresses (
                           id BIGSERIAL PRIMARY KEY,
                           user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           label address_label NOT NULL,
                           address_line TEXT NOT NULL,
                           city TEXT NOT NULL,
                           latitude DOUBLE PRECISION NOT NULL,
                           longitude DOUBLE PRECISION NOT NULL,
                           is_primary BOOLEAN DEFAULT FALSE,
                           created_at TIMESTAMP DEFAULT NOW(),
                           UNIQUE(user_id, label)
);

-- Restaurants Table
CREATE TABLE restaurants (
                             id BIGSERIAL PRIMARY KEY,
                             name TEXT NOT NULL,
                             address TEXT NOT NULL,
                             city TEXT NOT NULL,
                             latitude DOUBLE PRECISION NOT NULL,
                             longitude DOUBLE PRECISION NOT NULL,
                             rating SMALLINT CHECK (rating BETWEEN 1 AND 5),
                             is_active BOOLEAN DEFAULT TRUE,
                             created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                             created_at TIMESTAMP DEFAULT NOW(),
                             updated_at TIMESTAMP DEFAULT NOW(),
                             UNIQUE(created_by, name)
);

-- Menus Table
CREATE TABLE menus (
                       id BIGSERIAL PRIMARY KEY,
                       name TEXT NOT NULL,
                       description TEXT,
                       price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
                       is_available BOOLEAN DEFAULT TRUE,
                       food_type food_type NOT NULL,
                       category course_category NOT NULL,
                       restaurant_id BIGINT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
                       created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       created_at TIMESTAMP DEFAULT NOW(),
                       updated_at TIMESTAMP DEFAULT NOW(),
                       UNIQUE(restaurant_id, name)
);

-- Orders Table
CREATE TABLE orders (
                        id BIGSERIAL PRIMARY KEY,
                        user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                        restaurant_id BIGINT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
                        address_id BIGINT NOT NULL REFERENCES addresses(id) ON DELETE SET NULL,
                        status order_status NOT NULL DEFAULT 'placed',
                        total_amount NUMERIC(10,2) NOT NULL CHECK (total_amount >= 0),
                        payment_status payment_status DEFAULT 'pending',
                        payment_method TEXT,
                        paid_at TIMESTAMP,
                        created_at TIMESTAMP DEFAULT NOW(),
                        updated_at TIMESTAMP DEFAULT NOW()
);

-- Order Items Table
CREATE TABLE order_items (
                             id BIGSERIAL PRIMARY KEY,
                             order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
                             menu_id BIGINT NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
                             quantity INT NOT NULL CHECK (quantity > 0),
                             price NUMERIC(10,2) NOT NULL CHECK (price >= 0)
);

-- Payments Table (for mock/hardcoded payments)
CREATE TABLE payments (
                          id BIGSERIAL PRIMARY KEY,
                          order_id BIGINT NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
                          payment_method TEXT NOT NULL,
                          payment_status payment_status NOT NULL DEFAULT 'pending',
                          amount NUMERIC(10,2) NOT NULL CHECK (amount >= 0),
                          paid_at TIMESTAMP DEFAULT NOW(),
                          created_at TIMESTAMP DEFAULT NOW(),
                          updated_at TIMESTAMP DEFAULT NOW()
);

COMMIT;
