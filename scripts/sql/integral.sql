create table recharge_key
(
    id           int auto_increment
        primary key,
    recharge_key int                                not null comment '密钥',
    status       int                                not null comment '状态，0代表未使用，1代表已使用，2代表已失效',
    use_account  varchar(255)                       null comment '使用账号',
    created_time datetime default CURRENT_TIMESTAMP not null,
    updated_time datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    constraint recharge_key_recharge_key_uindex
        unique (recharge_key)
)
    comment '充值密钥';