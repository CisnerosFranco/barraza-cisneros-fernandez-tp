

CREATE OR REPLACE FUNCTION generarresumen(_nrocliente int, _año int, _mes int) returns void as $$
	DECLARE
		cliente_encontrado record;
		nro_resumen int;
		cierre_cliente record;
		total_resumen cabecera.total%type := 0;
		compra_aux record;
		
		nombre_comercio comercio.nombre%type;
		nrolinea int := 1;
		
		tarjetacliente record;
	
	BEGIN
		
		select * into cliente_encontrado from cliente where nrocliente = _nrocliente;
		
		
		if not found then
			raise notice 'CLIENTE NO ENCONTRADO!.';
		
		else
			for tarjetacliente in (select * from tarjeta where nrocliente = _nrocliente) loop
				nro_resumen := getnroresumen();
				select * into cierre_cliente from cierre c where c.año = _año and c.mes = _mes and c.terminacion = substring(tarjetacliente.nrotarjeta, 11, 1)::int;
			
				insert into cabecera values (nro_resumen, cliente_encontrado.nombre, cliente_encontrado.apellido, cliente_encontrado.domicilio, tarjetacliente.nrotarjeta, cierre_cliente.fechainicio, cierre_cliente.fechacierre, cierre_cliente.fechavto, 0);

				for compra_aux in (select * from compra where nrotarjeta = tarjetacliente.nrotarjeta and 
				fecha::date >= (cierre_cliente.fechainicio)::date and fecha::date <= (cierre_cliente.fechacierre)::date and pagado = false) loop
					
					nombre_comercio := (select nombre from comercio where nrocomercio = compra_aux.nrocomercio);
							
					insert into detalle values (nro_resumen, nrolinea, compra_aux.fecha, nombre_comercio, compra_aux.monto);
					total_resumen := total_resumen + compra_aux.monto;
					nrolinea = nrolinea + 1;
					
					update compra set pagado = true where nrooperacion = compra_aux.nrooperacion;
			
				end loop;
				
				update cabecera set total = total_resumen where nrotarjeta = tarjetacliente.nrotarjeta and desde = cierre_cliente.fechainicio and hasta = cierre_cliente.fechacierre;
				
			end loop;
		
		end if;
		
	END;
	$$ language plpgsql;
	
	


CREATE  OR REPLACE FUNCTION getnroresumen() RETURNS integer AS $$
	DECLARE
		
	   nro_resumen int;

	BEGIN
		select nroresumen into nro_resumen from cabecera order by nroresumen desc limit 1; 
		
		if not found then
			return 1;
		end if;
		
		nro_resumen := nro_resumen + 1;
		return nro_resumen;
	END;
	$$ language plpgsql;
















