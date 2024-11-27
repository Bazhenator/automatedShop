CREATE TABLE IF NOT EXISTS "users"
(
    id        bigserial PRIMARY KEY,
    login     varchar(30),
    pass_hash varchar(60),
    is_admin  bool
);

CREATE TABLE IF NOT EXISTS "expense_items"
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS "warehouses"
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(20),
    quantity INT,
    amount   INT
);

CREATE TABLE IF NOT EXISTS "charges"
(
    id              SERIAL PRIMARY KEY,
    amount          INT,
    charge_date     TIMESTAMP WITHOUT TIME ZONE,
    expense_item_id INT,
    CONSTRAINT fk_charges_expense_items
        FOREIGN KEY (expense_item_id)
            REFERENCES "expense_items" (id)
);

CREATE TABLE IF NOT EXISTS "sales"
(
    id            SERIAL PRIMARY KEY,
    amount        INT,
    quantity      INT,
    sale_date     TIMESTAMP WITHOUT TIME ZONE,
    warehouses_id INT,
    CONSTRAINT fk_sales_warehouses
        FOREIGN KEY (warehouses_id)
            REFERENCES "warehouses" (id)
);