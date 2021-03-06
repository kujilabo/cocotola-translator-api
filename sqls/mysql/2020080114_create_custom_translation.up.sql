create table `custom_translation` (
--  `id` int auto_increment
 `version` int not null default 1
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp on update current_timestamp
,`text` varchar(30) character set ascii not null
,`pos` int not null
,`lang2` varchar(2) character set ascii not null
,`translated` varchar(100) not null
,`disabled` tinyint(1) not null default 0
-- ,primary key('id')
,primary key(`text`, `pos`, `lang2`)
);
