create table 
  if not exists 
  letters
(
  id serial primary key,
  letter    text not null,
  email     varchar(255),
  date      timestamp,
  isActual  bool
)