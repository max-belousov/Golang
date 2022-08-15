# NetSqlGolangApp  
This is the first pet project on Golang.  
In this app there is client.exe and server.exe. In client you can choose options, server will get the result back from the database.  

I mastered:  
- Goroutines in Golang
- Golang syntax
- importing and using packages
- PostreSql in Golang

To run this project you need create the database with this SQL:

CREATE TABLE IF NOT EXISTS public.clients  
(  
id integer NOT NULL DEFAULT nextval('clients_id_seq'::regclass),  
firstname character varying(50) COLLATE pg_catalog."default" NOT NULL,  
secondname character varying(50) COLLATE pg_catalog."default" NOT NULL,  
phone character varying(15) COLLATE pg_catalog."default" NOT NULL,  
email character varying(50) COLLATE pg_catalog."default" NOT NULL,  
CONSTRAINT clients_pkey PRIMARY KEY (id)  
)
