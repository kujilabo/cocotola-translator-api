create table `custom_translation` (
 `id` integer primary key autoincrement
,`version` int not null default 1
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp
,`text` varchar(30) not null
,`pos` int not null
,`lang` varchar(2) not null
,`translated` varchar(100) not null
,unique(`text`, `pos`, `lang`)
);
