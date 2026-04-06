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
    image_url TEXT,            
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) NOT NULL UNIQUE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для ускорения
CREATE INDEX idx_favorites_user ON favorites(user_id);
CREATE INDEX idx_cart_user ON cart(user_id);
CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_order_items_order ON order_items(order_id);

-- Вставка товаров
INSERT INTO products (id, title, price, image_url, category, sizes, material, description) VALUES
(1, 'Спортивный костюм мужской Adidas', 7999, 'https://i.ebayimg.com/images/g/DRwAAOSw1q9l6c3m/s-l500.jpg', 'мужской', 'XS,S,M,L,XL,XXL', 'полиэстер', 'Стильный спортивный костюм для активного отдыха'),
(2, 'Спортивный костюм женский Nike', 6999, 'https://i.pinimg.com/564x/a0/be/6b/a0be6bf9078e73990c0ed67e5b003b75.jpg', 'женский', 'XS,S,M,L,XL,XXL', 'хлопок', 'Удобный костюм для фитнеса и йоги'),
(3, 'Детский спортивный костюм Puma', 3499, 'https://catalog-cdn.detmir.st/media/kGLgHPzDmalWkMFNdy3DEmMoc-8Wwqhv-Fy8PizFM-4=.webp?preset=site_product_gallery_r450', 'детский', '110,120,130,140,150,160', 'полиэстер', 'Яркий костюм для активных детей'),
(4, 'Мужской спортивный костюм Under Armour', 8999, 'https://goods-photos.static1-sima-land.com/items/6303013/5/1600.jpg?v=1633877733', 'мужской', 'XS,S,M,L,XL,XXL', 'нейлон', 'Профессиональный костюм для тренировок'),
(5, 'Женский спортивный костюм Reebok', 5999, 'https://basket-13.wbbasket.ru/vol1982/part198253/198253741/images/big/1.webp', 'женский', 'XS,S,M,L,XL,XXL', 'флис', 'Теплый костюм для зимних пробежек'),
(6, 'Детский спортивный костюм Adidas', 4599, 'https://www.famil.ru/shopfiles/resize_cache/2175599/502936d04bdb00a8d28d8a2a82e7c2bf/iblock/665/66527657dc797dde86db66cfd51f80df/2ed2ed1acf01e03edf2a76cd0fda228e.png', 'детский', '110,120,140,150,160', 'хлопок', 'Комфортный костюм для повседневной носки'),
(7, 'Мужской спортивный костюм Nike Tech', 9999, 'https://a.allegroimg.com/original/1134eb/4576cb874f7f9eb8fb849faca023/Dres-meski-sportowy-Nike-NSW-TECH-FLEECE-PANTS-HOODIE-BLACK', 'мужской', 'XS,S,M,L,XL,XXL', 'техническая ткань', 'Инновационный костюм с технологией Dri-FIT'),
(8, 'Спортивный костюм мужской Nike', 3790, 'https://a.lmcdn.ru/product/N/I/NI464EMACK72_1.jpg', 'мужской', 'XS,S,M,L,XL,2XL', 'полиэстер', 'Популярная модель 2025 года, качество и цена огонь!'),
(9, 'Спортивный костюм мужской Adidas Essentials', 6999, 'https://basket-29.wbbasket.ru/vol5563/part556331/556331049/images/big/1.webp', 'мужской', 'S,M,L,XL,XXL,3XL', 'полиэстер', 'Классический комплект из куртки и брюк с тремя полосками. Материал Primegreen.'),
(10, 'Спортивный костюм мужской Adidas Originals Firebird', 8999, 'https://m.media-amazon.com/images/I/91spqUOcK2L._AC_UY1100_.jpg', 'мужской', 'XS,S,M,L,XL,2XL', 'таслан', 'Винтажный дизайн в стиле Y2K. Куртка и брюки в комплекте.'),
(11, 'Спортивный костюм мужской Under Armour Challenger', 10999, 'https://ir.ozone.ru/s3/multimedia-a/6498306358.jpg', 'мужской', 'XS,S,M,L,XL,XXL', 'синтетическая ткань', 'Идеальное решение для любой тренировки. Легкий, теплый и хорошо тянется.'),
(12, 'Спортивный костюм женский Nike', 8999, 'https://basket-01.wbbasket.ru/vol110/part11031/11031124/images/big/1.webp', 'женский', 'XS,S,M,L,XL,XXL', 'хлопковый флис', 'Стильный костюм для фитнеса и повседневной носки.'),
(13, 'Спортивный костюм мужской Nike Tech Fleece', 12499, 'https://avatars.mds.yandex.net/i?id=c8cc2aecd8c9f300e23ce855c5091509_l-10755475-images-thumbs&n=13', 'мужской', 'XS,S,M,L,XL,XXL', 'технический флис', 'Комплект из толстовки с капюшоном на молнии и джоггеров. Инновационная ткань.'),
(14, 'Спортивный костюм детский Adidas', 5499, 'https://m.media-amazon.com/images/I/81dF+HR0G3L._AC_SY606._SX._UX._SY._UY_.jpg', 'детский', '110,120,130,140,150,160', 'хлопок', 'Яркий и удобный костюм для активных детей. Эластичный пояс для индивидуальной посадки.'),
(15, 'Спортивный костюм мужской Under Armour', 7999, 'https://n.cdn.cdek.shopping/images/shopping/cIv7i2W0hAadh6z5.jpg?v=1', 'мужской', 'XS,S,M,L,XL,XXL', 'флис', 'Базовый спортивный костюм. Ткань не скатывается и не теряет форму.'),
(16, 'Спортивный костюм мужской Nike Club', 6599, 'https://basket-01.wbbasket.ru/vol62/part6276/6276102/images/big/1.webp', 'мужской', 'XS,S,M,L,XL,XXL', 'полиэстер', 'Костюм с логотипом клуба. Комплект из толстовки с капюшоном и брюк.'),
(17, 'Спортивный костюм детский Adidas (пастельный)', 4999, 'https://a.lmcdn.ru/product/A/D/AD002EBCDAU4_6995569_1_v1.jpg', 'детский', '110,120,140,150,160', 'хлопковый трикотаж', 'Нежный пастельный цвет. Мягкая ткань для максимального комфорта ребенка.');

-- Сброс последовательности id после явной вставки
SELECT setval('products_id_seq', (SELECT MAX(id) FROM products));