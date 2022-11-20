CREATE DATABASE xm;
\c xm;
CREATE TYPE company_type AS ENUM ('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship');
DROP TABLE IF EXISTS companies;
CREATE TABLE companies (
  id uuid NOT NULL,
  name varchar(15) NOT NULL UNIQUE,
  description varchar(300),
  amount_of_employees integer NOT NULL,
  registered boolean NOT NULL,
  company_type company_type NOT NULL
);
