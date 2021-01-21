
CREATE DATABASE db_fbc;

\c db_fbc;



-- -------------------------------------------------------------


CREATE TABLE cliente(
	nrocliente int, --pk
	nombre text,
	apellido text, 
	domicilio text,
	telefono char(12)
);

CREATE TABLE tarjeta(
	nrotarjeta char(16), --pk
	nrocliente int,   --fk
	validadesde char(6),
	validahasta char(6),
	codseguridad char(4),
	limitecompra decimal(8,2),
	estado char(10)      -- 'vigente', 'suspendida', 'anulada'
);

CREATE TABLE comercio(
	nrocomercio int,--pk
	nombre text,
	domicilio text,
	codigopostal char(8),
	telefono char(12)
);

CREATE TABLE compra(
	nrooperacion int, -- pk
	nrotarjeta char(16), --fk
	nrocomercio int, --fk
	fecha timestamp,
	monto decimal(7,2),
	pagado boolean
);

CREATE TABLE rechazo(
	nrorechazo int, --pk
	nrotarjeta char(16), --no fk para poder registrar tarjetas no existentes
	nrocomercio int, --fk
	fecha timestamp,
	nonto decimal(7,2),
	motivo text
);

CREATE TABLE cierre(
	año int, -- pk,
	mes int, -- pk
	terminacion int, --pk
	fechainicio date,
	fechacierre date,
	fechavto date
);

CREATE TABLE cabecera(
	nroresumen int, --pk
	nombre text,
	apellido text,
	domicilio text,
	nrotarjeta char(16), --fk
	desde date,
	hasta date,
	vence date,
	total decimal(8,2)
);

CREATE TABLE detalle(
	nroresumen int, --pk
	nrolinea int, -- pk
	fecha date,
	nombrecomercio text,
	monto decimal(7,2)
);

CREATE TABLE alerta(
	nroalerta int, --pk
	nrotarjeta char(16), --fk
	fecha timestamp,
	nrorechazo int, --fk
	codalerta int, -- 0:rechazo, 1:compra 1min, 5:compra 5min, 32:límite
	descripcion text
);



-- Esta tabla no es parte del modelo de datos, pero se incluye para
-- poder probar las funciones.

CREATE TABLE consumo(
	nrotarjeta char(16), --fk
	codseguridad char(4),
	nrocomercio int, -- fk
	monto decimal(7,2)
);




-- agredamos las PKs

alter table cliente add constraint cliente_pk primary key (nrocliente);
alter table tarjeta add constraint tarjeta_pk primary key (nrotarjeta);
alter table comercio add constraint comercio_pk primary key (nrocomercio);
alter table compra add constraint compra_pk primary key (nrooperacion);
alter table rechazo add constraint rechazo_pk primary key (nrorechazo);
alter table cierre add constraint cierre_pk primary key (año, mes, terminacion);
alter table cabecera add constraint cabecera_pk primary key (nroresumen);
alter table detalle add constraint detalle_pk primary key (nroresumen, nrolinea);
alter table alerta add constraint alerta_pk primary key (nroalerta);




-- FKs
alter table tarjeta add constraint nrocliente_fk foreign key (nrocliente) references cliente (nrocliente);
alter table compra add constraint nrotarjeta_fk foreign key (nrotarjeta) references tarjeta (nrotarjeta);
alter table compra add constraint nrocomercio_fk foreign key (nrocomercio) references comercio (nrocomercio);
alter table rechazo add constraint nrocomercio_fk foreign key (nrocomercio) references comercio (nrocomercio);
alter table cabecera add constraint nrotarjeta_fk foreign key (nrotarjeta) references tarjeta (nrotarjeta);
alter table alerta add constraint nrorechazo_fk foreign key (nrorechazo) references rechazo (nrorechazo);

/*
alter table cliente drop constraint cliente_pk;
alter table tarjeta drop constraint tarjeta_pk;
alter table comercio drop constraint comercio_pk;
alter table compra drop constraint compra_pk;
alter table rechazo drop constraint rechazo_pk;
alter table cierre drop constraint cierre_pk ;
alter table cabecera drop constraint cabecera_pk;
alter table detalle drop constraint detalle_pk;
alter table alerta drop constraint alerta_pk;

alter table tarjeta drop constraint nrocliente_fk;
alter table compra drop constraint nrotarjeta_fk;
alter table compra drop constraint nrocomercio_fk;
alter table rechazo drop constraint nrocomercio_fk;
alter table cabecera drop constraint nrotarjeta_fk;
alter table alerta drop constraint nrorechazo_fk;

*/






















