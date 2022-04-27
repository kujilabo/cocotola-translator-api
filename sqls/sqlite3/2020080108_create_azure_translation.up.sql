create table `azure_translation` (
 `text` varchar(30) not null
,`lang2` varchar(2) not null
,`result` text not null
,primary key(`text`, `lang2`)
);
