-- +goose Up
-- The lease contract a landlord is prepared to sign, published up front so a
-- renter reads the terms (deposit, minimum term, utilities, pets) before
-- calling instead of discovering them at the viewing. Kept apart from the
-- `documents` array, which holds plans and title papers: this one is a single
-- file and earns its own button on the listing page.
--
-- Rent only, by design. A sale contract is notarial and negotiated at the deal,
-- so publishing a draft would inform nobody and would read as a public offer.
-- What a buyer wants — the title papers — already fits `documents`.
ALTER TABLE listings ADD COLUMN IF NOT EXISTS contract_url TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE listings DROP COLUMN IF EXISTS contract_url;
