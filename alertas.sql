

drop trigger alerta_ii on rechazo;
drop trigger if exists compraminuto_ia on compra;;
drop trigger if exists alertalimiterechazo on rechazo;
drop trigger if exists alertacompra5min on compra;

CREATE OR REPLACE FUNCTION insertaralertarechazo() RETURNS trigger AS $$
	DECLARE
		nro_alerta int;

	BEGIN
		nro_alerta := (select getnroalerta());
		
		insert into alerta values(nro_alerta, new.nrotarjeta, new.fecha, new.nrorechazo, 0, new.motivo);
		return new;
	END;
	$$ language plpgsql;
	

-- inser alerta - rechazo
CREATE TRIGGER alerta_ii after INSERT ON rechazo
FOR EACH ROW EXECUTE PROCEDURE insertaralertarechazo();




CREATE OR REPLACE FUNCTION compraminuto() RETURNS trigger AS $$
	DECLARE
		nro_alerta int;

	BEGIN
		nro_alerta := (select getnroalerta());

		if (existecompra(new.nrocomercio, new.nrooperacion, new.nrotarjeta)) = true then 
			insert into alerta(nroalerta, nrotarjeta, fecha, codalerta, descripcion)
			 values(nro_alerta, new.nrotarjeta, new.fecha, 1, '2 compras en menos de 1 minuto');
		end if;
		return new;
	END;
	$$ language plpgsql;


CREATE TRIGGER compraminuto_ia AFTER INSERT on compra
FOR EACH ROW EXECUTE PROCEDURE compraminuto();
 
 

	
	
	

/*Si una tarjeta registra dos compras en un lapso menor de 5 minutos en comercios
con diferentes códigos postales.*/

create or replace function compra5min() returns trigger as $$

declare
	cantidad int;
	nro_alerta int;

begin
	nro_alerta := (select getnroalerta());

	select count(distinct co.codigopostal) into cantidad from compra c, comercio co where c.nrocomercio = co.nrocomercio and c.nrotarjeta = new.nrotarjeta
	and c.fecha >= (new.fecha - interval '5 minute');

	if found  and cantidad >= 2 then
        insert into alerta values (nro_alerta, new.nrotarjeta, current_timestamp, null, 5, 'Se registro 2 compras en menos de 5 min, en comercios con diferente codigo postal');
	end if;
	return new;
end;

$$ language plpgsql;


/*Trigger*/
create trigger alertacompra5min
after insert on compra
for each row execute procedure compra5min();




/*Si una tarjeta registra dos rechazos por exceso de límite en el mismo día, la tarjeta
tiene que ser suspendida preventivamente, y se debe grabar una alerta asociada a
este cambio de estado.*/

create or replace function rechazolimite () returns trigger as $$

declare
	cant_rechazos int;
	nro_alerta int;

begin
	select getnroalerta() into nro_alerta;
	
    select count(*) into cant_rechazos from rechazo where nrotarjeta = new.nrotarjeta and
	 fecha::date = new.fecha::date and motivo ='se supera el limite de la tarjeta';

	if found and cant_rechazos >= 2 then
		update tarjeta set estado = 'suspendida' where nrotarjeta = new.nrotarjeta;
		insert into alerta values (nro_alerta, new.nrotarjeta, current_timestamp, new.nrorechazo, 32, 'limite de rechazos');
	end if;
	return new;
end;

$$ language plpgsql;


/*Trigger*/
create trigger alertalimiterechazo after insert on rechazo
for each row execute procedure rechazolimite();




-- esta funcion auxiliar te devuelve el siguiente nro de alerta.
CREATE OR REPLACE FUNCTION getnroalerta() returns int as $$
	declare
		nro_alerta int;
		
	begin
		select nroalerta into nro_alerta from alerta order by nroalerta desc limit 1 ;
		if found then 
			nro_alerta := nro_alerta + 1;
			return nro_alerta;
		else 
			return 1;
		end if;
	
	end;
	
	$$ language plpgsql;


 
-- te dice si existe una compra dentro del minuto
CREATE OR REPLACE FUNCTION existecompra(_nrocomercio int, _operacion int, _tarjeta text) RETURNS boolean AS $$
	DECLARE
		_fecha timestamp;
		_aux record;
		_comercio record;
	BEGIN
		
		select * into _comercio from comercio where nrocomercio = _nrocomercio;
		select fecha into _fecha from compra where nrooperacion = _operacion;
		
		select * into _aux from compra c, comercio o where c.nrocomercio != _nrocomercio AND c.nrotarjeta = _tarjeta 
		AND o.codigopostal = _comercio.codigopostal AND fecha >= (_fecha - interval '1 minute');

		if found then 
			return true;
		end if;
		return false;
	END;	
	$$ language plpgsql;







