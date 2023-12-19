create table balance_type
(
    balance_type_cd varchar(64)  not null
        primary key,
    comments        varchar(255) not null
)
    charset = latin1;

create table account_type
(
    account_type_cd varchar(64)  not null
        primary key,
    name            varchar(255) not null,
    balance_type_cd varchar(64)  not null,
    sequence_no     int          null,
    constraint fk_account_balance_type
        foreign key (balance_type_cd) references balance_type (balance_type_cd)
)
    comment 'account type (expense, income, assets etc)' charset = latin1;

create table account
(
    account_id      varchar(64)  not null
        primary key,
    code            varchar(64)  not null,
    name            varchar(255) not null,
    account_type_cd varchar(64)  not null,
    company_id      varchar(64)  not null,
    constraint uq_account
        unique (code) comment 'account code should be unique',
    constraint fk_account_type
        foreign key (account_type_cd) references account_type (account_type_cd)
)
    charset = latin1;

create table sequence_values
(
    seq_val int not null
)
    comment 'no primary key, just one record to keep values for sequences' charset = latin1;

create table user
(
    user_id     varchar(64)  not null
        primary key,
    provided_id varchar(64)  not null,
    login       varchar(64)  not null,
    name        varchar(255) not null
)
    charset = latin1;

create table company
(
    company_id varchar(64)   not null
        primary key,
    user_id    varchar(64)   not null,
    name       varchar(255)  null,
    month_end  int default 1 not null,
    day_end    int default 1 not null,
    constraint uq_company_name
        unique (name),
    constraint fk_company_user
        foreign key (user_id) references user (user_id)
)
    charset = latin1;

create table session
(
    session_id  varchar(64) not null
        primary key,
    user_id     varchar(64) not null,
    token       varchar(64) not null,
    expired_ind varchar(1)  not null,
    expiry_dt   datetime    not null,
    constraint uq_session_token
        unique (token) comment 'token should be unique 100 percent',
    constraint fk_session_user
        foreign key (user_id) references user (user_id)
)
    charset = latin1;

create table transaction
(
    transaction_id   varchar(64)  not null
        primary key,
    company_id       varchar(65)  not null,
    transaction_date varchar(32)  not null,
    sequence         int          not null,
    comments         varchar(255) null,
    constraint fk_transaction_company
        foreign key (company_id) references company (company_id)
)
    charset = latin1;

create table transaction_position
(
    transaction_position_id varchar(64)    not null
        primary key,
    transaction_id          varchar(64)    not null,
    account_id              varchar(64)    not null,
    sequence                mediumtext     not null,
    amount_debit            decimal(14, 2) null,
    amount_credit           decimal(14, 2) null,
    comments                varchar(255)   null,
    constraint fk_pos_account
        foreign key (account_id) references account (account_id),
    constraint fk_pos_transaction
        foreign key (transaction_id) references transaction (transaction_id)
)
    charset = latin1;

