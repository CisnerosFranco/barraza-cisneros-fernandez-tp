
CREATE OR REPLACE FUNCTION procesarconsumo() RETURNS void AS $$
	DECLARE
	
		v record;
		aux record;
		cont int := 1;
	BEGIN
		
		for v in (select * from consumo) loop
			select autorizacioncompra(v.nrotarjeta, v.codseguridad, v.nrocomercio, v.monto) into aux;
			raise notice 'cont : %',cont;
			cont := cont + 1;
		end loop;
		
	END;
	
	$$ language plpgsql;








