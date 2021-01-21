package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
	"os"
	"os/exec"
)

type cliente struct {
	nrocliente                  int
	nombre, apellido, domicilio string
	telefono                    string
}

type tarjeta struct {
	nrotarjeta   string
	nrocliente   int
	validadesde  string
	validahasta  string
	codseguridad string
	limitecompra float64
	estado       string
}

type compra struct {
	nrooperacion int
	nrotarjeta   string
	nrocomercio  int
	fecha        time.Time
	monto        float64
	pagado       bool
}

type consumo struct {
	nrotarjeta   string
	codseguridad string
	nrocomercio  int
	monto        float64
}

func main() {
	var eleccionUsuario string
	var menu = true
	
	for menu == true {
	
		eleccionUsuario = ""

		fmt.Printf(`
Bienvenido al SUPERMENU de Base de Datos
----------------------------------------
1) Crear la base de datos
2) Crear las tablas
3) Inicializar los datos
4) Cargar operaciones
5) Probar la tabla consumo
6) Borrar las PK y FK
7) Salir

Su eleccion: `)

		fmt.Scanf("%s", &eleccionUsuario)

		if eleccionUsuario == "1" {
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()

			go crearBDD()
			fmt.Printf("\nBase de Datos creada\n")
			time.Sleep(2 * time.Second)
		}

		if eleccionUsuario == "2" {
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()
			go crearTablas()
			fmt.Printf("\nTablas creadas\n\n")
			fmt.Printf("PK y FK inicializadas\n\n")
			time.Sleep(2 * time.Second)
		}

		if eleccionUsuario == "3" {
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()
			go inicializarDatos()
			fmt.Printf("\nDatos ingresados\n")
			time.Sleep(2 * time.Second)
		}

		if eleccionUsuario == "4" {
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()
			go autorizarCompra()
			fmt.Printf("\nAutorizaciones de compras inicializadas\n")
			
			go generarResumen()
			fmt.Printf("\nGeneraciones de resumenes activadas\n")
			
			go alertaRechazo()
			fmt.Printf("\nAdiciones en tabla rechazo por cada alerta activadas\n")

			go alertaCompraMinuto()
			fmt.Printf("\nAlertas por compras hace un minuto inicializadas\n")

			go alertaCompraCincoMinutos()
			fmt.Printf("\nAlertas por compras hace 5 minutos inicializadas\n")

			go rechazoTarjeta()
			fmt.Printf("\nAlertas por rechazo de tarjeta inicializadas\n\n")
			
			time.Sleep(2 * time.Second)
		}

		if eleccionUsuario == "5" {
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()
			go probarConsumo()
			fmt.Printf("\nTabla de consumos testeada\n\n")
			time.Sleep(2 * time.Second)
		}

		if eleccionUsuario == "6" {
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()
			go borrarPKyFK()
			fmt.Printf("PK y FK borradas\n")
			time.Sleep(2 * time.Second)
		}

		if eleccionUsuario == "7" {
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()
			menu = false
			fmt.Printf("\nAdiós!\n")
		}
	}
}

func crearBDD() {
	/*Accedo al usuario postgres*/
	db, err := sql.Open("postgres", "user = postgres host = localhost dbname = postgres sslmode = disable")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	/*Borro si existe la base de datos*/
	_, err = db.Exec(`drop database if exists tp_fbc`)

	if err != nil {
		log.Fatal(err)
	}

	/*Creo la base de datos*/
	_, err = db.Exec(`create database tp_fbc`)

	if err != nil {
		log.Fatal(err)
	}
}

func crearTablas() {
	db, err := sql.Open("postgres", "user = postgres host = localhost dbname = tp_fbc sslmode = disable")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	/*Creo las tablas*/
	_, err = db.Exec(`CREATE TABLE cliente (nrocliente int, nombre text, apellido text, domicilio text, telefono char(12));
					  CREATE TABLE tarjeta (nrotarjeta char(16), nrocliente int, validadesde char(6), validahasta char(6), codseguridad char(4), limitecompra decimal(8,2), estado char(10));
					  CREATE TABLE comercio (nrocomercio int, nombre text, domicilio text, codigopostal char(8), telefono char(12));
					  CREATE TABLE compra (nrooperacion int, nrotarjeta char(16), nrocomercio int, fecha timestamp, monto decimal(7,2), pagado boolean);                                                  
					  CREATE TABLE rechazo (nrorechazo int, nrotarjeta char(16), nrocomercio int, fecha timestamp, monto decimal(7,2), motivo text);
					  CREATE TABLE cierre (año int, mes int, terminacion int, fechainicio date, fechacierre date, fechavto date);                                                      
                      CREATE TABLE cabecera (nroresumen int, nombre text, apellido text, domicilio text, nrotarjeta char(16), desde date, hasta date, total decimal(8,2));                                    
					  CREATE TABLE detalle (nroresumen int, nrolinea int, fecha date, nombrecomercio text, monto decimal(7,2));                                
					  CREATE TABLE alerta (nroalerta int, nrotarjeta char(16), fecha timestamp, nrorechazo int, codalerta int, descripcion text);
					  CREATE TABLE consumo (nrotarjeta char(16), codseguridad char(4), nrocomercio int, monto decimal(7,2));`)

	if err != nil {
		log.Fatal(err)
	}
	
	/*INSERTO LOS PK's Y FK's*/
	_, err = db.Exec(`
					alter table cliente add constraint cliente_pk primary key(nrocliente);
					alter table tarjeta add constraint tarjeta_pk primary key(nrotarjeta);
					alter table comercio add constraint comercio_pk primary key(nrocomercio);
					alter table compra add constraint compra_pk primary key(nrooperacion);
					alter table rechazo add constraint rechazo_pk primary key(nrorechazo);
					alter table cierre add constraint cierre_pk primary key(año,mes,terminacion);
					alter table cabecera add constraint cabecera_pk primary key(nroresumen);
					alter table detalle add constraint detalle_pk primary key(nroresumen,nrolinea);
					alter table alerta add constraint alerta_pk primary key(nroalerta);
	
					alter table tarjeta add constraint tarjeta_nrocliente_fk foreign key(nrocliente) references cliente(nrocliente);
					alter table compra add constraint compra_nrotarjeta_fk foreign key(nrotarjeta) references tarjeta(nrotarjeta);
					alter table compra add constraint compra_nrocomercio_fk foreign key (nrocomercio) references comercio(nrocomercio);
					alter table rechazo add constraint rechazo_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					alter table rechazo add constraint rechazo_nrocomercio_fk foreign key (nrocomercio) references comercio(nrocomercio);
					alter table cabecera add constraint cabecera_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					alter table alerta add constraint alerta_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					alter table alerta add constraint alerta_nrorechazo_fk foreign key (nrorechazo) references rechazo(nrorechazo);	
					

					`)

	if err != nil {
		log.Fatal(err)
	}
}

func inicializarDatos() {
	db, err := sql.Open("postgres", "user = postgres host = localhost dbname = tp_fbc sslmode = disable")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	/*Inserto los datos en las tablas*/
	_, err = db.Exec(`
		

INSERT INTO cliente values(0,'Leanne Graham', 'Bret', 'Kulas Light 556', '1-770-736-80');
INSERT INTO cliente values(1,'Ervin Howell', 'Antonette', 'san martin 344', '1010-692-659');
INSERT INTO cliente values(2,'Clementine Bauch', 'Samantha', 'Douglas Extensionn 847', '1-463-123-44');
INSERT INTO cliente values(3,'Patricia Lebsack', 'Karianne', 'Hoeger Mall 692', '493-170-96');
INSERT INTO cliente values(4,'Chelsey Dietrich', 'Kamren', 'Skiles Walks 017', '254-954-1289');
INSERT INTO cliente values(5,'Mrs. Dennis Schulist', 'eopoldo_Corkery', 'Norberto Crossing 950', '1-477-935-84');
INSERT INTO cliente values(6,'Kurtis Weissnat', 'Elwyn Skiles', 'Rex Trail 280', '210.067.61');
INSERT INTO cliente values(7,'Nicholas Runolfsdottir', 'Maxime', 'Ellsworth Summit 729', '586.493.69');
INSERT INTO cliente values(8,'Glenna Reichert', 'Delphine', 'Dayna Park 49', '775-976-67');
INSERT INTO cliente values(9,'Clementina DuBuque', 'Stanton', 'Kattie Turnpike 198', '024-648-38');
INSERT INTO cliente values(10,'Louse', 'Lane', 'metropolis 500', '1-012-123');
INSERT INTO cliente values(11,'Marta', 'Kent', 'metopolis 344', '1-132-232');
INSERT INTO cliente values(12,'Marta', 'Wane', 'gotham 300', '12-336-456');
INSERT INTO cliente values(13,'Jhon', 'Jamison', 'manhatan 20', '12-122-323');
INSERT INTO cliente values(14,'Su', 'Storm', 'manhatan 123', '32-232-535');
INSERT INTO cliente values(15,'Rick', 'Richart', 'manhatan', '44-343-555');
INSERT INTO cliente values(16,'Victor', 'Bom Dom', 'ladveria 01', '00-000-000');
INSERT INTO cliente values(17,'Jerome', 'Valesca', 'gothan 100', '02-111-212');
INSERT INTO cliente values(18,'Oswald', 'Cabblepot', 'gothan 01', '12-132-121');
INSERT INTO cliente values(19,'Edward', 'Nigma', 'gothan 15', '32-231-232');



INSERT INTO tarjeta values('000-000-001', 0, '201201', '202501', '01', 100, 'vigente');
INSERT INTO tarjeta values('000-000-010', 1, '201712', '202511', '02', 100, 'vigente');
INSERT INTO tarjeta values('000-000-011', 2, '201811', '203011', '03', 100, 'vigente');
INSERT INTO tarjeta values('000-000-100', 3, '201908', '204012', '04', 100, 'vigente');
INSERT INTO tarjeta values('000-000-101', 4, '201707', '202512', '05', 100, 'vigente');
INSERT INTO tarjeta values('000-000-110', 5, '201706', '202212', '06', 100, 'vigente');
INSERT INTO tarjeta values('000-000-111', 6, '201701', '202302', '07', 100, 'vigente');
INSERT INTO tarjeta values('000-001-000', 7, '201701', '202303', '08', 100, 'vigente');
INSERT INTO tarjeta values('000-001-001', 8, '201702', '202405', '09', 100, 'vigente');
INSERT INTO tarjeta values('000-001-010', 9, '201703', '202302', '10', 100, 'vigente');
INSERT INTO tarjeta values('000-001-011', 10, '201704', '202411', '11', 100, 'vigente');
INSERT INTO tarjeta values('000-001-100', 11, '201704', '202511', '12', 100, 'vigente');
INSERT INTO tarjeta values('000-001-101', 12, '201705', '202111', '13', 100, 'vigente');
INSERT INTO tarjeta values('000-001-110', 13, '201705', '202201', '14', 100, 'vigente');
INSERT INTO tarjeta values('000-001-111', 14, '201801', '202003', '15', 100, 'vigente');
INSERT INTO tarjeta values('000-010-000', 15, '201503', '202104', '16', 100, 'vigente');
INSERT INTO tarjeta values('000-010-001', 16, '201605', '202207', '17', 100, 'vigente');
INSERT INTO tarjeta values('000-010-010', 17, '201512', '202308', '18', 100, 'vigente');
INSERT INTO tarjeta values('000-010-011', 18, '201311', '202509', '19', 100, 'vigente');
INSERT INTO tarjeta values('000-010-100', 19, '201411', '202509', '20', 100, 'vigente');
INSERT INTO tarjeta values('000-010-101', 18, '201312', '201901', '21', 100, 'anulada');
INSERT INTO tarjeta values('000-010-110', 19, '201410', '201902', '22', 100, 'anulada');




INSERT INTO comercio values(001, 'Kwekii Mark', 'ricardo rojas 001', '100', '12-120-001');
INSERT INTO comercio values(010, 'La Ostra', 'ricardo rojas 010', '100', '12-120-010');
INSERT INTO comercio values(011, 'MOE', 'ricardo rojas 011', '100', '12-120-011');
INSERT INTO comercio values(100, 'Anteiku', 'ricardo rojas 100', '101', '12-120-100');
INSERT INTO comercio values(101, 'Hichiraku', 'ricardo rojas 101', '101', '12-120-101');
INSERT INTO comercio values(110, 'MacDonals', 'ricardo rojas 110', '110', '12-120-110');
INSERT INTO comercio values(111, 'Burguer King', 'ricardo rojas 111', '110', '12-120-111');
INSERT INTO comercio values(1000, 'Pizzetas', 'ricardo rojas 1000', '010', '12-121-000');
INSERT INTO comercio values(1001, 'Sandwiches y mas', 'ricardo rojas 010', '100', '12-121-001');
INSERT INTO comercio values(1010, 'Verduras', 'ricardo rojas 1010', '011', '12-121-010');
INSERT INTO comercio values(1011, 'El Tinto', 'ricardo rojas 1011', '111', '12-121-011');
INSERT INTO comercio values(1100, 'Crustacio Cascarudo', 'ricardo rojas 110', '100', '12-121-100');
INSERT INTO comercio values(1101, 'La Hora Feliz', 'ricardo rojas 1101', '1001', '12-121-101');
INSERT INTO comercio values(1110, 'la Hora Infeliz', 'ricardo rojas 1110', '100', '12-121-110');
INSERT INTO comercio values(1111, 'Pixels', 'ricardo rojas 1111', '1000', '12-121-111');
INSERT INTO comercio values(10000, 'Clasic Cars', 'ricardo rojas 10000', '110', '12-110-000');
INSERT INTO comercio values(10001, 'Inmuebles', 'ricardo rojas 10001', '100', '12-110-000');
INSERT INTO comercio values(10010, 'Electronicos', 'ricardo rojas 10010', '1010', '12-110-001');
INSERT INTO comercio values(10011, 'Relojes', 'ricardo rojas 10011', '1011', '12-110-010');
INSERT INTO comercio values(10100, 'Quesitos', 'ricardo rojas 10100', '1111', '12-110-011');




-- CIERRES

INSERT INTO cierre values(2020, 1, 0, '2020-01-01', '2020-02-01', '2020-02-15');
INSERT INTO cierre values(2020, 2, 0, '2020-02-01', '2020-03-01', '2020-03-15');
INSERT INTO cierre values(2020, 3, 0, '2020-03-01', '2020-04-01', '2020-04-15');
INSERT INTO cierre values(2020, 4, 0, '2020-04-01', '2020-05-01', '2020-05-15');
INSERT INTO cierre values(2020, 5, 0, '2020-05-01', '2020-06-01', '2020-06-15');
INSERT INTO cierre values(2020, 6, 0, '2020-06-01', '2020-07-01', '2020-07-15');
INSERT INTO cierre values(2020, 7, 0, '2020-07-01', '2020-08-01', '2020-08-15');
INSERT INTO cierre values(2020, 8, 0, '2020-08-01', '2020-09-01', '2020-09-15');
INSERT INTO cierre values(2020, 9, 0, '2020-09-01', '2020-10-01', '2020-10-15');
INSERT INTO cierre values(2020, 10, 0, '2020-10-01', '2020-11-01', '2020-11-15');
INSERT INTO cierre values(2020, 11, 0, '2020-11-01', '2020-12-01', '2020-12-15');
INSERT INTO cierre values(2020, 12, 0, '2020-12-01', '2021-01-01', '2021-01-15');

INSERT INTO cierre values(2020, 1, 1, '2020-01-02', '2020-02-02', '2020-02-15');
INSERT INTO cierre values(2020, 2, 1, '2020-02-02', '2020-03-02', '2020-03-15');
INSERT INTO cierre values(2020, 3, 1, '2020-03-02', '2020-04-02', '2020-04-15');
INSERT INTO cierre values(2020, 4, 1, '2020-04-02', '2020-05-02', '2020-05-15');
INSERT INTO cierre values(2020, 5, 1, '2020-05-02', '2020-06-02', '2020-06-15');
INSERT INTO cierre values(2020, 6, 1, '2020-06-02', '2020-07-02', '2020-07-15');
INSERT INTO cierre values(2020, 7, 1, '2020-07-02', '2020-08-02', '2020-08-15');
INSERT INTO cierre values(2020, 8, 1, '2020-08-02', '2020-09-02', '2020-09-15');
INSERT INTO cierre values(2020, 9, 1, '2020-09-02', '2020-10-02', '2020-10-15');
INSERT INTO cierre values(2020, 10, 1, '2020-10-02', '2020-11-02', '2020-11-15');
INSERT INTO cierre values(2020, 11, 1, '2020-11-02', '2020-12-02', '2020-12-15');
INSERT INTO cierre values(2020, 12, 1, '2020-12-02', '2021-01-02', '2021-01-15');

INSERT INTO cierre values(2020, 1, 2, '2020-01-03', '2020-02-03', '2020-02-15');
INSERT INTO cierre values(2020, 2, 2, '2020-02-03', '2020-03-03', '2020-03-15');
INSERT INTO cierre values(2020, 3, 2, '2020-03-03', '2020-04-03', '2020-04-15');
INSERT INTO cierre values(2020, 4, 2, '2020-04-03', '2020-05-03', '2020-05-15');
INSERT INTO cierre values(2020, 5, 2, '2020-05-03', '2020-06-03', '2020-06-15');
INSERT INTO cierre values(2020, 6, 2, '2020-06-03', '2020-07-03', '2020-07-15');
INSERT INTO cierre values(2020, 7, 2, '2020-07-03', '2020-08-03', '2020-08-15');
INSERT INTO cierre values(2020, 8, 2, '2020-08-03', '2020-09-03', '2020-09-15');
INSERT INTO cierre values(2020, 9, 2, '2020-09-03', '2020-10-03', '2020-10-15');
INSERT INTO cierre values(2020, 10, 2, '2020-10-03', '2020-11-03', '2020-11-15');
INSERT INTO cierre values(2020, 11, 2, '2020-11-03', '2020-12-03', '2020-12-15');
INSERT INTO cierre values(2020, 12, 2, '2020-12-03', '2021-01-03', '2021-01-15');

INSERT INTO cierre values(2020, 1, 3, '2020-01-04', '2020-02-04', '2020-02-15');
INSERT INTO cierre values(2020, 2, 3, '2020-02-04', '2020-03-04', '2020-03-15');
INSERT INTO cierre values(2020, 3, 3, '2020-03-04', '2020-04-04', '2020-04-15');
INSERT INTO cierre values(2020, 4, 3, '2020-04-04', '2020-05-04', '2020-05-15');
INSERT INTO cierre values(2020, 5, 3, '2020-05-04', '2020-06-04', '2020-06-15');
INSERT INTO cierre values(2020, 6, 3, '2020-06-04', '2020-07-04', '2020-07-15');
INSERT INTO cierre values(2020, 7, 3, '2020-07-04', '2020-08-04', '2020-08-15');
INSERT INTO cierre values(2020, 8, 3, '2020-08-04', '2020-09-04', '2020-09-15');
INSERT INTO cierre values(2020, 9, 3, '2020-09-04', '2020-10-04', '2020-10-15');
INSERT INTO cierre values(2020, 10, 3, '2020-10-04', '2020-11-04', '2020-11-15');
INSERT INTO cierre values(2020, 11, 3, '2020-11-04', '2020-12-04', '2020-12-15');
INSERT INTO cierre values(2020, 12, 3, '2020-12-04', '2021-01-04', '2021-01-15');

INSERT INTO cierre values(2020, 1, 4, '2020-01-05', '2020-02-05', '2020-02-15');
INSERT INTO cierre values(2020, 2, 4, '2020-02-05', '2020-03-05', '2020-03-15');
INSERT INTO cierre values(2020, 3, 4, '2020-03-05', '2020-04-05', '2020-04-15');
INSERT INTO cierre values(2020, 4, 4, '2020-04-05', '2020-05-05', '2020-05-15');
INSERT INTO cierre values(2020, 5, 4, '2020-05-05', '2020-06-05', '2020-06-15');
INSERT INTO cierre values(2020, 6, 4, '2020-06-05', '2020-07-05', '2020-07-15');
INSERT INTO cierre values(2020, 7, 4, '2020-07-05', '2020-08-05', '2020-08-15');
INSERT INTO cierre values(2020, 8, 4, '2020-08-05', '2020-09-05', '2020-09-15');
INSERT INTO cierre values(2020, 9, 4, '2020-09-05', '2020-10-05', '2020-10-15');
INSERT INTO cierre values(2020, 10, 4, '2020-10-05', '2020-11-05', '2020-11-15');
INSERT INTO cierre values(2020, 11, 4, '2020-11-05', '2020-12-05', '2020-12-15');
INSERT INTO cierre values(2020, 12, 4, '2020-12-05', '2021-01-05', '2021-01-15');

INSERT INTO cierre values(2020, 1, 5, '2020-01-06', '2020-02-06', '2020-02-15');
INSERT INTO cierre values(2020, 2, 5, '2020-02-06', '2020-03-06', '2020-03-15');
INSERT INTO cierre values(2020, 3, 5, '2020-03-06', '2020-04-06', '2020-04-15');
INSERT INTO cierre values(2020, 4, 5, '2020-04-06', '2020-05-06', '2020-05-15');
INSERT INTO cierre values(2020, 5, 5, '2020-05-06', '2020-06-06', '2020-06-15');
INSERT INTO cierre values(2020, 6, 5, '2020-06-06', '2020-07-06', '2020-07-15');
INSERT INTO cierre values(2020, 7, 5, '2020-07-06', '2020-08-06', '2020-08-15');
INSERT INTO cierre values(2020, 8, 5, '2020-08-06', '2020-09-06', '2020-09-15');
INSERT INTO cierre values(2020, 9, 5, '2020-09-06', '2020-10-06', '2020-10-15');
INSERT INTO cierre values(2020, 10, 5, '2020-10-06', '2020-11-06', '2020-11-15');
INSERT INTO cierre values(2020, 11, 5, '2020-11-06', '2020-12-06', '2020-12-15');
INSERT INTO cierre values(2020, 12, 5, '2020-12-06', '2021-01-06', '2021-01-15');

INSERT INTO cierre values(2020, 1, 6, '2020-01-07', '2020-02-07', '2020-02-15');
INSERT INTO cierre values(2020, 2, 6, '2020-02-07', '2020-03-07', '2020-03-15');
INSERT INTO cierre values(2020, 3, 6, '2020-03-07', '2020-04-07', '2020-04-15');
INSERT INTO cierre values(2020, 4, 6, '2020-04-07', '2020-05-07', '2020-05-15');
INSERT INTO cierre values(2020, 5, 6, '2020-05-07', '2020-06-07', '2020-06-15');
INSERT INTO cierre values(2020, 6, 6, '2020-06-07', '2020-07-07', '2020-07-15');
INSERT INTO cierre values(2020, 7, 6, '2020-07-07', '2020-08-07', '2020-08-15');
INSERT INTO cierre values(2020, 8, 6, '2020-08-07', '2020-09-07', '2020-09-15');
INSERT INTO cierre values(2020, 9, 6, '2020-09-07', '2020-10-07', '2020-10-15');
INSERT INTO cierre values(2020, 10, 6, '2020-10-07', '2020-11-07', '2020-11-15');
INSERT INTO cierre values(2020, 11, 6, '2020-11-07', '2020-12-07', '2020-12-15');
INSERT INTO cierre values(2020, 12, 6, '2020-12-07', '2021-01-07', '2021-01-15');

INSERT INTO cierre values(2020, 1, 7, '2020-01-08', '2020-02-08', '2020-02-15');
INSERT INTO cierre values(2020, 2, 7, '2020-02-08', '2020-03-08', '2020-03-15');
INSERT INTO cierre values(2020, 3, 7, '2020-03-08', '2020-04-08', '2020-04-15');
INSERT INTO cierre values(2020, 4, 7, '2020-04-08', '2020-05-08', '2020-05-15');
INSERT INTO cierre values(2020, 5, 7, '2020-05-08', '2020-06-08', '2020-06-15');
INSERT INTO cierre values(2020, 6, 7, '2020-06-08', '2020-07-08', '2020-07-15');
INSERT INTO cierre values(2020, 7, 7, '2020-07-08', '2020-08-08', '2020-08-15');
INSERT INTO cierre values(2020, 8, 7, '2020-08-08', '2020-09-08', '2020-09-15');
INSERT INTO cierre values(2020, 9, 7, '2020-09-08', '2020-10-08', '2020-10-15');
INSERT INTO cierre values(2020, 10, 7, '2020-10-08', '2020-11-08', '2020-11-15');
INSERT INTO cierre values(2020, 11, 7, '2020-11-08', '2020-12-08', '2020-12-15');
INSERT INTO cierre values(2020, 12, 7, '2020-12-08', '2021-01-08', '2021-01-15');


INSERT INTO cierre values(2020, 1, 8, '2020-01-09', '2020-02-09', '2020-02-15');
INSERT INTO cierre values(2020, 2, 8, '2020-02-09', '2020-03-09', '2020-03-15');
INSERT INTO cierre values(2020, 3, 8, '2020-03-09', '2020-04-09', '2020-04-15');
INSERT INTO cierre values(2020, 4, 8, '2020-04-09', '2020-05-09', '2020-05-15');
INSERT INTO cierre values(2020, 5, 8, '2020-05-09', '2020-06-09', '2020-06-15');
INSERT INTO cierre values(2020, 6, 8, '2020-06-09', '2020-07-09', '2020-07-15');
INSERT INTO cierre values(2020, 7, 8, '2020-07-09', '2020-08-09', '2020-08-15');
INSERT INTO cierre values(2020, 8, 8, '2020-08-09', '2020-09-09', '2020-09-15');
INSERT INTO cierre values(2020, 9, 8, '2020-09-09', '2020-10-09', '2020-10-15');
INSERT INTO cierre values(2020, 10, 8, '2020-10-09', '2020-11-09', '2020-11-15');
INSERT INTO cierre values(2020, 11, 8, '2020-11-09', '2020-12-09', '2020-12-15');
INSERT INTO cierre values(2020, 12, 8, '2020-12-09', '2021-01-09', '2021-01-15');


INSERT INTO cierre values(2020, 1, 9, '2020-01-10', '2020-02-10', '2020-02-15');
INSERT INTO cierre values(2020, 2, 9, '2020-02-10', '2020-03-10', '2020-03-15');
INSERT INTO cierre values(2020, 3, 9, '2020-03-10', '2020-04-10', '2020-04-15');
INSERT INTO cierre values(2020, 4, 9, '2020-04-10', '2020-05-10', '2020-05-15');
INSERT INTO cierre values(2020, 5, 9, '2020-05-10', '2020-06-10', '2020-06-15');
INSERT INTO cierre values(2020, 6, 9, '2020-06-10', '2020-07-10', '2020-07-15');
INSERT INTO cierre values(2020, 7, 9, '2020-07-10', '2020-08-10', '2020-08-15');
INSERT INTO cierre values(2020, 8, 9, '2020-08-10', '2020-09-10', '2020-09-15');
INSERT INTO cierre values(2020, 9, 9, '2020-09-10', '2020-10-10', '2020-10-15');
INSERT INTO cierre values(2020, 10, 9, '2020-10-10', '2020-11-10', '2020-11-15');
INSERT INTO cierre values(2020, 11, 9, '2020-11-10', '2020-12-10', '2020-12-15');
INSERT INTO cierre values(2020, 12, 9, '2020-12-10', '2021-01-10', '2021-01-15');








-- datos de consumo para las pruebas

INSERT INTO consumo VALUES ('000-000-001', '01', 001, 70); --pasa 
INSERT INTO consumo VALUES ('000-000-001', '01', 111, 10); --pasa  --alerta 5: compras 5m/codigopostal
INSERT INTO consumo VALUES ('000-000-010', '02', 011, 30); --pasa
INSERT INTO consumo VALUES ('000-000-011', '03', 101, 90); --pasa
INSERT INTO consumo VALUES ('000-000-100', '04', 011, 20); --pasa 
INSERT INTO consumo VALUES ('000-000-101', '05', 100, 60); --pasa
INSERT INTO consumo VALUES ('000-001-110', '14', 1111, 80); --pasa
INSERT INTO consumo VALUES ('000-000-111', '07', 001, 70);--pasa
INSERT INTO consumo VALUES ('000-001-000', '08', 001, 40); --pasa
INSERT INTO consumo VALUES ('000-001-001', '09', 101, 40); --pasa 
INSERT INTO consumo VALUES ('000-001-010', '10', 1001, 40);--pasa --alerta 1: compras 1 minuto
INSERT INTO consumo VALUES ('000-001-011', '11', 1101, 40);--rechazo por ecceso
INSERT INTO consumo VALUES ('000-001-100', '12', 10100, 40);--pasa --alerta 1: compras 1 minuto
INSERT INTO consumo VALUES ('000-001-101', '13', 10010, 40);--pasa --alerta 1: compras 1 minuto
INSERT INTO consumo VALUES ('000-001-110', '14', 10100, 40);--rechazo por ecceso
INSERT INTO consumo VALUES ('000-000-111', '07', 111, 20);  --pasa --alerta 5: 5m/codigopostal
INSERT INTO consumo VALUES ('000-010-010', '18', 111, 40);  --pasa
INSERT INTO consumo VALUES ('000-010-001', '17', 001, 40);  --pasa
INSERT INTO consumo VALUES ('000-010-010', '18', 001, 80);  --rechazo por ecceso
INSERT INTO consumo VALUES ('000-000-011', '03', 001, 40);  --rechazo por ecceso alerta 32: queda suspendida
INSERT INTO consumo VALUES ('000-010-100', '00', 001, 40);  --rechazo, codigo invalido
INSERT INTO consumo VALUES ('000-010-101', '21', 001, 40); --tarjeta vencia
INSERT INTO consumo VALUES ('000-010-110', '22', 001, 40); --tarjeta vencida
	
	`)

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`select * from cliente`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var a cliente
	for rows.Next() {
		if err := rows.Scan(&a.nrocliente, &a.nombre, &a.apellido, &a.domicilio, &a.telefono); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v %v %v %v %v\n", a.nrocliente, a.nombre, a.apellido, a.domicilio, a.telefono)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	
}

func borrarPKyFK() {

	/*Inicializo Base de Datos*/
	db, err := sql.Open("postgres", "user = postgres host = localhost dbname = tp_fbc sslmode = disable")

	if err != nil {
		log.Fatal(err);
	}

	defer db.Close()
	
    _, err = db.Exec(`
					alter table tarjeta drop constraint tarjeta_nrocliente_fk;
					alter table compra drop constraint compra_nrotarjeta_fk;
					alter table compra drop constraint compra_nrocomercio_fk;
					alter table rechazo drop constraint rechazo_nrotarjeta_fk;
					alter table rechazo drop constraint rechazo_nrocomercio_fk;
					alter table cabecera drop constraint cabecera_nrotarjeta_fk;
					alter table alerta drop constraint alerta_nrotarjeta_fk;
					alter table alerta drop constraint alerta_nrorechazo_fk;
									
					alter table alerta drop constraint alerta_pk;
					alter table rechazo drop constraint rechazo_pk;
					alter table detalle drop constraint detalle_pk;
					alter table cabecera drop constraint cabecera_pk;
					alter table cierre drop constraint cierre_pk;
					alter table compra drop constraint compra_pk;
					alter table comercio drop constraint comercio_pk;
					alter table tarjeta drop constraint tarjeta_pk;
					alter table cliente drop constraint cliente_pk;
				   `)
		if err != nil {
		log.Fatal(err);
	}

}

func probarConsumo() {
	/*Inicializo Base de Datos*/
	db, err := sql.Open("postgres", "user = postgres host = localhost dbname = tp_fbc sslmode = disable")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err2 := db.Exec(`
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
`)

	if err2 != nil {
		log.Fatal(err2)
	}

	_, err3 := db.Exec(`select procesarconsumo();`)

	if err3 != nil {
		log.Fatal(err3)
	}
}

func autorizarCompra() {
	/*Inicializo Base de Datos*/
	db, err := sql.Open("postgres", "user = postgres host = localhost dbname = tp_fbc sslmode = disable")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err2 := db.Exec(`
		
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


	`)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func generarResumen() {
	/*Accedo al usuario postgres*/
	db, err := sql.Open("postgres", "user = postgres host = localhost dbname = postgres sslmode = disable")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err2 := db.Exec(`
		
CREATE OR REPLACE FUNCTION generarresumen(_nrocliente int, _año int, _mes int) returns void as $$
	DECLARE
		cliente_encontrado record;
		nro_resumen int;
		cierre_cliente record;
		total_resumen decimal(8,2); 
		compra_aux record;
		
		nombre_comercio text;
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



`)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func alertaRechazo() {
	/*Inicializo Base de Datos*/
	db, err := sql.Open("postgres", "user = postgres host = localhost dbname = tp_fbc sslmode = disable")

	if err != nil {
		log.Fatal(err)
	}

	_, err2 := db.Exec(`
	drop trigger if exists alerta_ii on rechazo;
	
	CREATE OR REPLACE FUNCTION insertaralertarechazo() RETURNS trigger AS $$
	DECLARE
		nro_alerta int;

	BEGIN
		nro_alerta := (select getnroalerta());
		
		insert into alerta values(nro_alerta, new.nrotarjeta, new.fecha, new.nrorechazo, 0, new.motivo);
		return new;
		
	END;
	$$ language plpgsql;
	
	CREATE TRIGGER alerta_ii after INSERT ON rechazo
	FOR EACH ROW EXECUTE PROCEDURE insertaralertarechazo();
	`)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func alertaCompraMinuto() {
	/*Inicializo Base de Datos*/
	db, err := sql.Open("postgres", "user = postgres host = localhost dbname = tp_fbc sslmode = disable")

	if err != nil {
		log.Fatal(err)
	}

	_, err2 := db.Exec(`
	drop trigger if exists compraminuto_ia on compra;
	
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
	`)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func alertaCompraCincoMinutos() {
	/*Inicializo Base de Datos*/
	db, err := sql.Open("postgres", "user = postgres host = localhost dbname = tp_fbc sslmode = disable")

	if err != nil {
		log.Fatal(err)
	}

	_, err2 := db.Exec(`
	drop trigger if exists alertacompra5min on compra;
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
	for each row execute procedure compra5min();`)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func rechazoTarjeta() {
	/*Inicializo Base de Datos*/
	db, err := sql.Open("postgres", "user = postgres host = localhost dbname = tp_fbc sslmode = disable")

	if err != nil {
		log.Fatal(err)
	}

	_, err2 := db.Exec(`
	drop trigger if exists alertalimiterechazo on rechazo;
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

	create trigger alertalimiterechazo after insert on rechazo
	for each row execute procedure rechazolimite();
	
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
	`)

	if err2 != nil {
		log.Fatal(err2)
	}
}
