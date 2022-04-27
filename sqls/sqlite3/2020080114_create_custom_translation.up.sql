create table `custom_translation` (
 `text` varchar(30) not null
,`pos` int not null
,`lang2` varchar(2) not null
,`translated` varchar(100) not null
,`disabled` tinyint(1) not null
,primary key(`text`, `pos`, `lang2`)
);
