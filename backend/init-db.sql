-- Create necessary types
CREATE TYPE invoice_status AS ENUM ('PENDING', 'PAID', 'CANCELED');

-- Create invoices table
CREATE TABLE IF NOT EXISTS invoices (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(255) NOT NULL UNIQUE,
    amount NUMERIC(10,2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'USDC',
    description TEXT,
    status invoice_status NOT NULL DEFAULT 'PENDING',
    due_date TIMESTAMP NOT NULL,
    receiver_addr VARCHAR(255) NOT NULL,
    link_token VARCHAR(255) NOT NULL UNIQUE,
    sender_details JSONB NOT NULL,
    recipient_details JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_invoices_link_token ON invoices(link_token);
CREATE INDEX IF NOT EXISTS idx_invoices_invoice_number ON invoices(invoice_number);
CREATE INDEX IF NOT EXISTS idx_invoices_receiver_addr ON invoices(receiver_addr);

-- Add some example data
INSERT INTO invoices (
    invoice_number,
    amount,
    currency,
    description,
    status,
    due_date,
    receiver_addr,
    link_token,
    sender_details,
    recipient_details
) VALUES (
    'INV-001',
    100.00,
    'USDC',
    'Web development services',
    'PENDING',
    NOW() + INTERVAL '7 days',
    'GsbwXfJraMomNxBcjcNRiGKuPGsJxkVyecCw3VYP5wTZ', -- Example Solana address (replace with your own)
    'aaa8d5dc-1a8a-5bc7-9c0e-d3de82081c88',
    '{"name": "Acme Inc.", "email": "billing@acme.com", "address": "123 Business St, Anytown, AN 12345"}',
    '{"name": "Client LLC", "email": "accounts@client.com", "address": "456 Corporate Ave, Business City, BC 67890"}'
);

INSERT INTO invoices (
    invoice_number,
    amount,
    currency,
    description,
    status,
    due_date,
    receiver_addr,
    link_token,
    sender_details,
    recipient_details
) VALUES (
    'INV-002',
    250.50,
    'USDC',
    'Consulting services',
    'PENDING',
    NOW() + INTERVAL '14 days',
    'GsbwXfJraMomNxBcjcNRiGKuPGsJxkVyecCw3VYP5wTZ', -- Example Solana address (replace with your own)
    'bbb8d5dc-2a8a-5bc7-9c0e-d3de82081c99',
    '{"name": "Acme Inc.", "email": "billing@acme.com", "address": "123 Business St, Anytown, AN 12345"}',
    '{"name": "Another Corp", "email": "finance@anothercorp.com", "address": "789 Enterprise Blvd, Metropolis, MP 54321"}'
); 