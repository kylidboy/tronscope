CREATE TABLE `block` (
    `number` bigint(10) not null primary key,
    `witnessid` bigint(10) not null,
    `timestamp` int not null,
    `parenthash` varchar(64) not null,
    `witnessaddr` varchar(64) not null,
    `witnesssig` varchar(128) not null
    index `idx_timestamp` (`timestamp`)
) default charset=utf8;

CREATE TABLE `tx` (
    `id` bigint(10) not null auto_increment primary key,
    `refblocknumber` bigint(10) not null,
    `expiration` int not null,
    `timestamp` int not null,
    `type` varchar(64) not null,
    `refblockhash` varchar(64) not null,
    `scripts` text null,
    `date` text null,
    `sigs` text null,
    index `idx_refblocknumber` (`refblocknumber`)
) default charset=utf8;