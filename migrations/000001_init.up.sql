CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users ( 
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS items (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  raw_name VARCHAR(255) NOT NULL,
  quantity INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  invoice_url TEXT,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  delivered_at TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  item_id uuid CONSTRAINT order_item_item_fk REFERENCES items (id),
  order_id uuid CONSTRAINT order_item_order_fk REFERENCES orders (id),
  quantity INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS schools (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  name VARCHAR (100) NOT NULL UNIQUE,
  address VARCHAR (255),
  phone_number VARCHAR (11),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS school_orders (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  school_id uuid CONSTRAINT school_order_school_fk REFERENCES schools (id),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  delivered_at TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS school_items (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  item_id uuid CONSTRAINT school_item_item_fk REFERENCES items (id),
  quantity INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS school_order_items (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  item_id uuid CONSTRAINT school_order_item_item_fk REFERENCES items (id),
  order_id uuid CONSTRAINT school_order_item_order_fk REFERENCES school_orders (id),
  quantity INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP
);