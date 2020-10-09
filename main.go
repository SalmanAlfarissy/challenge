package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// Buku struct (Model)
type Buku struct {
	IDbuku        string `json:"id_buku"`
	Kodebuku      string `form:"kode_buku" json:"kode_buku"`
	Judulbuku     string `form:"judul_buku" json:"judul_buku"`
	Penulisbuku   string `form:"penulis_buku" json:"penulis_buku"`
	Penerbitbuku  string `form:"penerbit_buku" json:"penerbit_buku"`
	Tahunpenerbit int    `form:"tahun_penerbit" json:"tahun_penerbit"`
	Stok          int    `form:"stok" json:"stok"`
}

func getBuku(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tbbuku []Buku

	sql := `SELECT id_buku,kode_buku,judul_buku,penulis_buku,penerbit_buku,tahun_penerbit, stok FROM tb_buku`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var tbbukus Buku
		err := result.Scan(&tbbukus.IDbuku, &tbbukus.Kodebuku, &tbbukus.Judulbuku,
			&tbbukus.Penulisbuku, &tbbukus.Penerbitbuku, &tbbukus.Tahunpenerbit, &tbbukus.Stok)

		if err != nil {
			panic(err.Error())
		}
		tbbuku = append(tbbuku, tbbukus)
	}

	json.NewEncoder(w).Encode(tbbuku)

}

func createBuku(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		IDbuku := r.FormValue("id_buku")
		Kodebuku := r.FormValue("kode_buku")
		Judulbuku := r.FormValue("judul_buku")
		Penulisbuku := r.FormValue("penulis_buku")
		Penerbitbuku := r.FormValue("penerbit_buku")
		Tahunpenerbit := r.FormValue("tahun_penerbit")
		Stok := r.FormValue("stok")

		stmt, err := db.Prepare("INSERT INTO tb_buku (id_buku,kode_buku,judul_buku,penulis_buku,penerbit_buku,tahun_penerbit,stok) VALUES (?,?,?,?,?,?,?)")
		_, err = stmt.Exec(IDbuku, Kodebuku, Judulbuku, Penulisbuku, Penerbitbuku, Tahunpenerbit, Stok)

		if err != nil {
			fmt.Fprintf(w, "Data Duplikat")
		} else {
			fmt.Fprintf(w, "Data Create")
		}

	}

}

func updateBuku(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		IDbuku := r.FormValue("id_buku")
		Kodebuku := r.FormValue("kode_buku")
		Judulbuku := r.FormValue("judul_buku")
		Penulisbuku := r.FormValue("penulis_buku")
		Penerbitbuku := r.FormValue("penerbit_buku")
		Tahunpenerbit := r.FormValue("tahun_penerbit")
		Stok := r.FormValue("stok")

		stmt, err := db.Prepare("UPDATE tb_buku set kode_buku = ?, judul_buku = ?, penulis_buku = ?, penerbit_buku = ?, tahun_penerbit = ?, stok = ? WHERE id_buku = ?")

		if err != nil {
			fmt.Fprintf(w, "statement error!!")
		} else {
			result, _ := stmt.Exec(Kodebuku, Judulbuku, Penulisbuku, Penerbitbuku, Tahunpenerbit, Stok, IDbuku)
			affect, _ := result.RowsAffected()

			if affect == 0 {
				fmt.Fprintf(w, "Data tidak ditemukan,Update gagal!!")
			} else {
				fmt.Fprintf(w, "Buku dengan IDbuku = %s berhasil di update", IDbuku)
			}

		}

	}

}

func deleteBuku(w http.ResponseWriter, r *http.Request) {

	if r.Method == "DELETE" {
		IDbuku := r.FormValue("id_buku")
		Kodebuku := r.FormValue("kode_buku")

		stmt, err := db.Prepare("DELETE FROM tb_buku WHERE id_buku =? AND kode_buku =?")

		if err != nil {
			fmt.Fprintf(w, "statement error!!")
		} else {

			result, _ := stmt.Exec(IDbuku, Kodebuku)
			affect, _ := result.RowsAffected()

			if affect == 0 {
				fmt.Fprintf(w, "delete gagal data tidak di temukan!!")
			} else {
				fmt.Fprintf(w, "Buku dengan IDbuku = %s dan Kodebuku = %s berhasil di delete", IDbuku, Kodebuku)

			}

		}

	}

}

func getSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tbbuku []Buku

	IDbuku := r.FormValue("id_buku")
	Kodebuku := r.FormValue("kode_buku")

	sql := `SELECT id_buku,kode_buku,judul_buku,penulis_buku,penerbit_buku,tahun_penerbit, stok
			FROM tb_buku WHERE id_buku = ? AND kode_buku = ?`
	result, err := db.Query(sql, IDbuku, Kodebuku)

	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var tbbukus Buku
	for result.Next() {

		err := result.Scan(&tbbukus.IDbuku, &tbbukus.Kodebuku, &tbbukus.Judulbuku, &tbbukus.Penulisbuku, &tbbukus.Penerbitbuku, &tbbukus.Tahunpenerbit, &tbbukus.Stok)
		if err != nil {
			panic(err.Error())
		}

		tbbuku = append(tbbuku, tbbukus)

	}
	json.NewEncoder(w).Encode(tbbuku)
}

func main() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_pustaka")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/tb_buku", getBuku).Methods("GET")
	r.HandleFunc("/tb_buku", createBuku).Methods("POST")
	r.HandleFunc("/tb_buku", updateBuku).Methods("PUT")
	r.HandleFunc("/tb_buku", deleteBuku).Methods("DELETE")

	r.HandleFunc("/getsearch", getSearch).Methods("POST")

	fmt.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
