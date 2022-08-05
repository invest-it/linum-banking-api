BEGIN;

SET client_encoding = "UTF8";

CREATE TABLE UserRequisitions (
    RequisitionId varchar(36) PRIMARY KEY,
    UserId varchar(36) NOT NULL
);

COMMIT;