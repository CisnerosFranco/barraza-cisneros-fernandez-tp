\c db_fbc;

/*
	nrorechazo int, --pk
	nrotarjeta char(16), 
	nrocomercio int, --fk
	fecha timestamp,
	nonto decimal(7,2),
	motivo text
*/
CREATE OR REPLACE FUNCTION autorizacioncompra(nrotarjeta char(16), codseguridad char(6), nrocomercio int, monto decimal(7,2) ) returns boolean AS $$ 
	DECLARE
	
		nro_rechazo int;
		aux record;
		
	BEGIN
		-- conseguimos el siguiente nrorechazo
		select nrorechazo into nro_rechazo from rechazo order by nrorechazo desc limit 1;
		if found then
			nro_rechazo := nro_rechazo + 1;
		else
			nro_rechazo := 1;
		end if;
		
		------------------------------------
	
		if (select validar_vigencia_tarjeta(nrotarjeta)) = false then
			insert into rechazo values(nro_rechazo, nrotarjeta, nrocomercio, (select current_timestamp), monto, 'tarjeta no valida o no vigente');
			return false;
		end if;
		
		if (select validar_cod_seguridad(nrotarjeta, codseguridad)) = false then
			insert into rechazo values(nro_rechazo, nrotarjeta, nrocomercio, (select current_timestamp), monto, 'codigo de seguridad invalido');
			return false;
		end if;
		
		if (select validar_monto(nrotarjeta, nrocomercio, monto)) = false then
			insert into rechazo values(nro_rechazo, nrotarjeta, nrocomercio, (select current_timestamp), monto, 'se supera el limite de la tarjeta');
			return false;
		end if;
		
		if (select tarjeta_vencida(nrotarjeta)) = false then
			insert into rechazo values(nro_rechazo, nrotarjeta, nrocomercio, (select current_timestamp), monto, 'plazo de vigencia expirado');
			return false;
		end if;
		
		if (select tarjeta_suspendida(nrotarjeta)) = false then
			insert into rechazo values(nro_rechazo, nrotarjeta, nrocomercio, (select current_timestamp), monto, 'la tarjeta se encuentra suspendida');
			return false;
		end if;
		
		
		select agregar_compra(nrotarjeta, codseguridad, nrocomercio, monto) into aux;
		return true;
		
	END;
	$$ language plpgsql;



-- ------------ funciones auxiliares - -----------------------------------------



-- agrega un nuevo registro a la tabla compra.
CREATE OR REPLACE FUNCTION agregar_compra(nrotarjeta char(16), codseguridad char(6), nrocomercio int, monto decimal(7,2)) RETURNS void AS $$
	DECLARE
		ultimaoperacion int;
	
	BEGIN
		
		select nrooperacion into ultimaoperacion from compra order by nrooperacion desc limit 1;
		
		if found then
			insert into compra values((ultimaoperacion + 1), nrotarjeta, nrocomercio, (select current_timestamp), monto, false);
		else
			insert into compra values(0, nrotarjeta, nrocomercio, (select current_timestamp), monto, false);
		end if;
	END;
	$$ LANGUAGE plpgsql;


-- te dice si la targeta existe y si esta en Estado Vigente
CREATE OR REPLACE FUNCTION validar_vigencia_tarjeta(_nrotarjeta text) returns boolean AS $$
	DECLARE
		vigencia text;
	
	BEGIN
		select estado into vigencia from tarjeta where nrotarjeta = _nrotarjeta;
		
		if found AND vigencia = 'vigente' then
			return true;
		end if;
		
		return false;
	
	END;
	$$ language plpgsql;
	

-- valida que el codigo coincida con el de la targeta
CREATE OR REPLACE FUNCTION validar_cod_seguridad(_nrotarjeta text, codigo text) returns boolean AS $$
	DECLARE
		valor text;
	
	BEGIN
		select codseguridad into valor from tarjeta where nrotarjeta = _nrotarjeta;
		
		if found AND valor = codigo then
			return true;
		end if;
		return false;
	
	END;
	$$ language plpgsql;


-- te dice si el la suma de las compras sin pagar relacionadas a esta targeta + el monto actual 
-- no superan el limite de la targeta
CREATE OR REPLACE FUNCTION validar_monto(_nrotarjeta char(16), _nrocomercio int, _monto decimal(7,2)) returns boolean AS $$
	DECLARE
		suma decimal;
		limitetarjeta decimal;
		
	BEGIN
		select sum(compra.monto) into suma from compra where nrotarjeta =_nrotarjeta AND pagado = false;
		select limitecompra into limitetarjeta from tarjeta where nrotarjeta = _nrotarjeta;
		
		if found then
			suma := suma + _monto;
		else
			suma := monto;
			
		end if;
		
		if limitetarjeta < suma then
			raise notice 'limite superado % : %  : %', limitetarjeta, suma, _nrotarjeta;
			return false;
		end if;

		return true;
		
	END;
	$$ language plpgsql;
	

-- te dice si la targeta no se encuentra vencida
CREATE OR REPLACE FUNCTION tarjeta_vencida(_nrotarjeta char(16)) returns boolean as $$
	DECLARE
		actualfecha char(6) = (select extract(year from current_timestamp) || '' ||  extract(month from current_timestamp));
		fechatarjeta char(6);
		
	BEGIN
		select validahasta into fechatarjeta from tarjeta where nrotarjeta = _nrotarjeta;
		
		if fechatarjeta < actualfecha then
			return false;
		
		end if;
		return true;
	END;
	$$ language plpgsql;



 -- te dice si la targeta no se encuentra suspendida
CREATE OR REPLACE FUNCTION tarjeta_suspendida(_nrotarjeta char(16)) returns boolean as $$
	DECLARE
		estadotarjeta text;
		
	BEGIN
		select estado into estadotarjeta from tarjeta where nrotarjeta = _nrotarjeta;
		
		if estadotarjeta = 'suspendida' then
			return false;
		end if;
		
		return true;
	END;
	$$ language plpgsql;

















