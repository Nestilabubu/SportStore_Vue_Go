-- Создание таблиц
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    image_url TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    sizes TEXT NOT NULL, -- храним как строку через запятую, например "M,L,XL"
    material VARCHAR(100),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE favorites (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, product_id)
);

CREATE TABLE cart (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    size VARCHAR(10) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, product_id, size)
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    total_price INTEGER NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    size VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для ускорения
CREATE INDEX idx_favorites_user ON favorites(user_id);
CREATE INDEX idx_cart_user ON cart(user_id);
CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_order_items_order ON order_items(order_id);

-- Вставка товаров
INSERT INTO products (id, title, price, image_url, category, sizes, material, description) VALUES
(1, 'Спортивный костюм мужской Adidas', 7999, 'https://i.ebayimg.com/images/g/DRwAAOSw1q9l6c3m/s-l500.jpg', 'мужской', 'M,L,XL,XXL', 'полиэстер', 'Стильный спортивный костюм для активного отдыха'),
(2, 'Спортивный костюм женский Nike', 6999, 'https://i.pinimg.com/564x/a0/be/6b/a0be6bf9078e73990c0ed67e5b003b75.jpg', 'женский', 'XS,S,M,L', 'хлопок', 'Удобный костюм для фитнеса и йоги'),
(3, 'Детский спортивный костюм Puma', 3499, 'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_450,h_450/global/584859/01/fnd/EEA/fmt/png/Minicats-Crew-Babies''-Jogger', 'детский', '110,120,130,140', 'полиэстер', 'Яркий костюм для активных детей'),
(4, 'Мужской спортивный костюм Under Armour', 8999, 'https://goods-photos.static1-sima-land.com/items/6303013/5/1600.jpg?v=1633877733', 'мужской', 'S,M,L,XL', 'нейлон', 'Профессиональный костюм для тренировок'),
(5, 'Женский спортивный костюм Reebok', 5999, 'https://cdn1.ozone.ru/s3/multimedia-1-6/7354618026.jpg', 'женский', 'XS,S,M,L', 'флис', 'Теплый костюм для зимних пробежек'),
(6, 'Детский спортивный костюм Adidas', 4599, 'https://a.lmcdn.ru/product/A/D/AD093EBIALT6_10183752_1_v1.jpg', 'детский', '110,120,140,150,160', 'хлопок', 'Комфортный костюм для повседневной носки'),
(7, 'Мужской спортивный костюм Nike Tech', 9999, 'https://a.allegroimg.com/original/1134eb/4576cb874f7f9eb8fb849faca023/Dres-meski-sportowy-Nike-NSW-TECH-FLEECE-PANTS-HOODIE-BLACK', 'мужской', 'M,L,XL,XXL', 'техническая ткань', 'Инновационный костюм с технологией Dri-FIT'),
(8, 'Спортивный костюм мужской Nike (Хит 2025)', 3790, 'https://a.lmcdn.ru/product/N/I/NI464EMACK72_1.jpg', 'мужской', 'S,M,L,XL,2XL', 'полиэстер', 'Популярная модель 2025 года, качество и цена огонь!'),
(9, 'Спортивный костюм мужской Adidas Essentials', 6999, 'https://avatars.mds.yandex.net/get-mpic/5251502/2a00000194e0e64ebd46c28a79d0e76be175/orig', 'мужской', 'S,M,L,XL,XXL,3XL', 'полиэстер', 'Классический комплект из куртки и брюк с тремя полосками. Материал Primegreen.'),
(10, 'Спортивный костюм мужской Adidas Originals Firebird', 8999, 'https://m.media-amazon.com/images/I/91spqUOcK2L._AC_UY1100_.jpg', 'мужской', 'S,M,L,XL,2XL', 'таслан', 'Винтажный дизайн в стиле Y2K. Куртка и брюки в комплекте.'),
(11, 'Спортивный костюм мужской Under Armour Challenger', 10999, 'https://a.lmcdn.ru/product/U/N/UN001EBMKGN6_13137947_1_v1.jpg', 'мужской', 'S,M,L,XL', 'синтетическая ткань', 'Идеальное решение для любой тренировки. Легкий, теплый и хорошо тянется.'),
(12, 'Спортивный костюм женский Nike', 8999, 'https://cdn.sportmaster.ru/upload/mdm/media_content/resize/75c/768_1024_e6ae/87823600299.jpg', 'женский', 'XS,S,M,L', 'хлопковый флис', 'Стильный костюм для фитнеса и повседневной носки.'),
(13, 'Спортивный костюм мужской Nike Tech Fleece', 12499, 'https://a.lmcdn.ru/product/R/T/RTLACE462001_18595114_1_v1.jpg', 'мужской', 'S,M,L,XL', 'технический флис', 'Комплект из толстовки с капюшоном на молнии и джоггеров. Инновационная ткань.'),
(14, 'Спортивный костюм детский Adidas', 5499, 'https://m.media-amazon.com/images/I/81dF+HR0G3L._AC_SY606._SX._UX._SY._UY_.jpg', 'детский', '110,120,130,140,150,160', 'хлопок', 'Яркий и удобный костюм для активных детей. Эластичный пояс для индивидуальной посадки.'),
(15, 'Спортивный костюм мужской Under Armour', 7999, 'https://n.cdn.cdek.shopping/images/shopping/cIv7i2W0hAadh6z5.jpg?v=1', 'мужской', 'S,M,L,XL', 'флис', 'Базовый спортивный костюм. Ткань не скатывается и не теряет форму.'),
(16, 'Спортивный костюм мужской Nike Club', 6599, 'https://cdn.sportmaster.ru/upload/mdm/media_content/resize/7d8/768_1024_48f1/51057360299.jpg', 'мужской', 'M,L,XL,XXL', 'полиэстер', 'Костюм с логотипом клуба. Комплект из толстовки с капюшоном и брюк.'),
(17, 'Спортивный костюм детский Adidas (пастельный)', 4999, 'https://a.lmcdn.ru/product/R/T/RTLAEP899901_30130997_1_v1_2x.jpg', 'детский', '110,120,140,150,160', 'хлопковый трикотаж', 'Нежный пастельный цвет. Мягкая ткань для максимального комфорта ребенка.');

-- Сброс последовательности id после явной вставки
SELECT setval('products_id_seq', (SELECT MAX(id) FROM products));