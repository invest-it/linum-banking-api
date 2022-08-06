BEGIN;

SET client_encoding = "UTF8";

CREATE TABLE UserRequisitions (
    RequisitionId varchar(36) PRIMARY KEY,
    RequisitionReference varchar(256),
    UserId varchar(36) NOT NULL,
    Approved boolean NOT NULL
);

COMMIT;