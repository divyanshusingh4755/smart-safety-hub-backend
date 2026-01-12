CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    company_id UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_email ON users(email);

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_name ON roles(name);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_permissions_name ON permissions(name);

CREATE TABLE roles_permissions (
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY(role_id, permission_id)
);

CREATE INDEX idx_role_id ON roles_permissions(role_id);

CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE brands (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    logo_url TEXT,
    website_url TEXT,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    parent_id UUID REFERENCES categories(id),
    level INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE status_enum AS ENUM('DRAFT', 'ACTIVE', 'ARCHIVED');

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    seller_id UUID NOT NULL,
    brand_id UUID REFERENCES brands(id) ON DELETE SET NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    status status_enum DEFAULT 'DRAFT',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products_attributes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    attribute_key VARCHAR(100),
    attribute_value VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE products_attributes ADD CONSTRAINT product_key_unique UNIQUE (product_id, attribute_key);

CREATE TABLE product_options (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product_option_values(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    option_id UUID REFERENCES product_options(id) ON DELETE CASCADE,
    value VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(100) UNIQUE NOT NULL,
    price DECIMAL(12,2) NOT NULL,
    weight DECIMAL(8, 2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE variant_option_values (
    variant_id UUID REFERENCES product_variants(id) ON DELETE CASCADE,
    option_value_id UUID REFERENCES product_option_values(id) ON DELETE CASCADE,
    PRIMARY KEY(variant_id, option_value_id)
);

CREATE TYPE media_enum as ENUM('image', 'video', 'pdf');

CREATE TABLE product_media (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    variant_id UUID REFERENCES product_variants(id),
    url TEXT NOT NULL,
    type media_enum DEFAULT 'image',
    display_order INT DEFAULT 0
);

CREATE TABLE product_seo (
    product_id UUID PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    meta_title VARCHAR(255),
    meta_description TEXT,
    keywords JSONB DEFAULT '[]'::jsonb,
    og_image_url TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Indexes for faster searches
CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_variants_sku ON product_variants(sku);

CREATE INDEX idx_products_brand_id ON products(brand_id);
CREATE INDEX idx_products_category_id ON products(category_id);

-- Create the function
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpsql';

-- Apply it to your main table
CREATE TRIGGER update_products_modtime BEFORE UPDATE ON products FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE TRIGGER update_variants_modtime BEFORE UPDATE ON product_variants FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- Data for roles and permissions
INSERT INTO roles (id, name, description)
VALUES
(uuid_generate_v4(), 'ADMIN',  'Platform administrator with full system access'),
(uuid_generate_v4(), 'SELLER', 'Seller who manages catalog, pricing, inventory, and orders'),
(uuid_generate_v4(), 'BUYER',  'Buyer who browses products, places orders, and manages payments');


INSERT INTO permissions (id, name, description) VALUES
(uuid_generate_v4(), 'auth:login', 'Login to the system'),
(uuid_generate_v4(), 'auth:logout', 'Logout from the system'),
(uuid_generate_v4(), 'user:create', 'Create users'),
(uuid_generate_v4(), 'user:update', 'Update user details'),
(uuid_generate_v4(), 'user:view', 'View user details'),
(uuid_generate_v4(), 'user:deactivate', 'Deactivate users');

INSERT INTO permissions (id, name, description) VALUES
(uuid_generate_v4(), 'company:view', 'View company profile'),
(uuid_generate_v4(), 'company:update', 'Update company profile'),
(uuid_generate_v4(), 'company:manage_users', 'Manage company users');

INSERT INTO permissions (id, name, description) VALUES
(uuid_generate_v4(), 'catalog:view', 'View product catalog'),
(uuid_generate_v4(), 'catalog:create', 'Create products'),
(uuid_generate_v4(), 'catalog:update', 'Update products'),
(uuid_generate_v4(), 'catalog:delete', 'Delete products');

INSERT INTO permissions (id, name, description) VALUES
(uuid_generate_v4(), 'pricing:view', 'View pricing'),
(uuid_generate_v4(), 'pricing:update', 'Update pricing rules'),
(uuid_generate_v4(), 'rfq:create', 'Create RFQ'),
(uuid_generate_v4(), 'rfq:view', 'View RFQ'),
(uuid_generate_v4(), 'rfq:respond', 'Respond to RFQ');

INSERT INTO permissions (id, name, description) VALUES
(uuid_generate_v4(), 'inventory:view', 'View inventory'),
(uuid_generate_v4(), 'inventory:update', 'Update inventory levels');

INSERT INTO permissions (id, name, description) VALUES
(uuid_generate_v4(), 'order:create', 'Create order'),
(uuid_generate_v4(), 'order:view', 'View orders'),
(uuid_generate_v4(), 'order:update_status', 'Update order status'),
(uuid_generate_v4(), 'order:cancel', 'Cancel order');

INSERT INTO permissions (id, name, description) VALUES
(uuid_generate_v4(), 'payment:view', 'View payments'),
(uuid_generate_v4(), 'payment:initiate', 'Initiate payment'),
(uuid_generate_v4(), 'invoice:view', 'View invoices'),
(uuid_generate_v4(), 'invoice:download', 'Download invoices');

INSERT INTO permissions (id, name, description) VALUES
(uuid_generate_v4(), 'compliance:view', 'View compliance documents'),
(uuid_generate_v4(), 'compliance:upload', 'Upload compliance documents'),
(uuid_generate_v4(), 'document:upload', 'Upload documents'),
(uuid_generate_v4(), 'document:view', 'View documents');

INSERT INTO permissions (id, name, description) VALUES
(uuid_generate_v4(), 'report:view', 'View reports'),
(uuid_generate_v4(), 'report:export', 'Export reports'),
(uuid_generate_v4(), 'audit:view', 'View audit logs'),
(uuid_generate_v4(), 'system:configure', 'Configure system settings');

INSERT INTO roles_permissions (role_id, permission_id) SELECT '37e13c1b-cfb5-44ad-a2ac-613d8e9650b4' AS role_id, id AS permission_id FROM permissions;

INSERT INTO roles_permissions (role_id, permission_id)
SELECT 'a5c16d38-8a3a-49bd-874c-94ab690e314a', id
FROM permissions
WHERE name IN (
  -- Auth
  'auth:login', 'auth:logout',

  -- Company
  'company:view',

  -- Catalog
  'catalog:view', 'catalog:create', 'catalog:update', 'catalog:delete',

  -- Pricing
  'pricing:view', 'pricing:update',

  -- RFQ
  'rfq:view', 'rfq:respond',

  -- Inventory
  'inventory:view', 'inventory:update',

  -- Orders
  'order:view', 'order:update_status',

  -- Compliance & Documents
  'compliance:view', 'compliance:upload',
  'document:upload', 'document:view',

  -- Reports
  'report:view'
);

INSERT INTO roles_permissions (role_id, permission_id)
SELECT '1c0cc3a1-f9e6-4ba7-adae-7d8a1c3128bd', id
FROM permissions
WHERE name IN (
  -- Auth
  'auth:login', 'auth:logout',

  -- Catalog
  'catalog:view',

  -- RFQ
  'rfq:create', 'rfq:view',

  -- Orders
  'order:create', 'order:view', 'order:cancel',

  -- Payments
  'payment:view', 'payment:initiate',

  -- Invoices
  'invoice:view', 'invoice:download',

  -- Compliance & Documents
  'compliance:view',
  'document:view'
);
