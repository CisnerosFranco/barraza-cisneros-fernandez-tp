package main

import (
    "encoding/json"
    bolt "github.com/coreos/bbolt"
    "log"
    "strconv"
)

type Cliente struct {
    Nrocliente int
	Nombre string
	Apellido string 
	Domicilio string
	Telefono string
}

type Tarjeta struct {
   Nrotarjeta string
	Nrocliente string
	Validadesde string
	Validahasta string
	Codseguridad string
	Limitecompra string
	Estado string 
}
type Comercio struct {
    Nrocomercio int
	Nombre string
	Domicilio string
	Codigopostal string
	Telefono string
}
type Compra struct {
    Nrooperacion int
	Nrotarjeta string
	Nrocomercio string
	Fecha string
	Monto string
	Pagado string
}

func CreateUpdate(db *bolt.DB, bucketName string, key []byte, val []byte) error {
    // abre transacción de escritura
    tx, err := db.Begin(true)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))

    err = b.Put(key, val)
    if err != nil {
        return err
    }

    // cierra transacción
    if err := tx.Commit(); err != nil {
        return err
    }

    return nil
}

func main() {
    
    datosClientes()
    datosTarjetas()
    datosComercios()
    datosCompras()
}


func datosClientes() {
    db, err := bolt.Open("tp_fbc.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    
   leanne:= Cliente{0,"Leanne Graham", "Bret", "Kulas Light 556", "1-770-736-80"}
   ervin:= Cliente{1,"Ervin Howell", "Antonette", "san martin 344", "1010-692-659"}
   clementine:= Cliente{2,"Clementine Bauch", "Samantha", "Douglas Extensionn 847", "1-463-123-44"}
   
   data, err := json.Marshal(leanne)
    if err != nil {
        log.Fatal(err)
    }
   data2, err := json.Marshal(ervin)
    if err != nil {
        log.Fatal(err)
    }
    data3, err := json.Marshal(clementine)
    if err != nil {
        log.Fatal(err)
    }

    CreateUpdate(db, "cliente", []byte(strconv.Itoa(leanne.Nrocliente)), data)

    CreateUpdate(db, "cliente", []byte(strconv.Itoa(ervin.Nrocliente)), data2)

    CreateUpdate(db, "cliente", []byte(strconv.Itoa(clementine.Nrocliente)), data3)

}

func datosTarjetas() {
    db, err := bolt.Open("tp_fbc.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

   tarjeta:= Tarjeta{"000-000-001", "0", "201201", "202501", "01", "100", "vigente"};
   tarjeta2:= Tarjeta{"000-000-010", "1", "201712", "202511", "02", "100", "vigente"};
   tarjeta3:= Tarjeta{"000-000-011", "2", "201811", "203011", "03", "100", "vigente"};

   data, err := json.Marshal(tarjeta)
    if err != nil {
        log.Fatal(err)
    }
   data2, err := json.Marshal(tarjeta2)
    if err != nil {
        log.Fatal(err)
    }
    data3, err := json.Marshal(tarjeta3)
    if err != nil {
        log.Fatal(err)
    }
    
    CreateUpdate(db, "tarjeta", []byte(tarjeta.Nrotarjeta), data)

    CreateUpdate(db, "tarjeta", []byte(tarjeta2.Nrotarjeta), data2)

    CreateUpdate(db, "tarjeta", []byte(tarjeta3.Nrotarjeta), data3)
 
}

func datosComercios() {
    db, err := bolt.Open("tp_fbc.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    
   mark:= Comercio{001, "Kwekii Mark", "ricardo rojas 001", "100", "12-120-001"};
   ostra:= Comercio {010, "La Ostra", "ricardo rojas 010", "100", "12-120-010"};
   moe:= Comercio {011, "MOE", "ricardo rojas 011", "100", "12-120-011"};

   data, err := json.Marshal(mark)
    if err != nil {
        log.Fatal(err)
    }
   data2, err := json.Marshal(ostra)
    if err != nil {
        log.Fatal(err)
    }
    data3, err := json.Marshal(moe)
    if err != nil {
        log.Fatal(err)
    }

    CreateUpdate(db, "comercio", []byte(strconv.Itoa(mark.Nrocomercio)), data)
    
    CreateUpdate(db, "comercio", []byte(strconv.Itoa(ostra.Nrocomercio)), data2)

    CreateUpdate(db, "comercio", []byte(strconv.Itoa(moe.Nrocomercio)), data3)

}

func datosCompras() {
    db, err := bolt.Open("tp_fbc.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

   compra1:= Compra {1,"000-010-001", "001", "2020-01-10", "40","true"};
   compra2:= Compra {2,"000-010-010", "001", "2020-01-10", "40","false"};
   compra3:= Compra {3,"000-010-011", "001", "2020-01-10", "40","true"};
   data, err := json.Marshal(compra1)
    if err != nil {
        log.Fatal(err)
    }
   data2, err := json.Marshal(compra2)
    if err != nil {
        log.Fatal(err)
    }
    data3, err := json.Marshal(compra3)
    if err != nil {
        log.Fatal(err)
    }

    CreateUpdate(db, "compra", []byte(strconv.Itoa(compra1.Nrooperacion)), data)

    CreateUpdate(db, "compra", []byte(strconv.Itoa(compra2.Nrooperacion)), data2)

    CreateUpdate(db, "compra", []byte(strconv.Itoa(compra3.Nrooperacion)), data3)

}
