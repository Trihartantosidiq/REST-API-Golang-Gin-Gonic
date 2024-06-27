package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Harus Masukan Nomor Absen Terlebih Dahulu")
		return
	}

	// Ambil nomor absen dari argumen CLI
	nomorAbsen := os.Args[1]

	// Konversi nomor absen ke integer
	absen, err := strconv.Atoi(nomorAbsen)
	if err != nil {
		fmt.Println("Harus masukan nomor absen dengan benar")
		return
	}

	// Cari data teman berdasarkan nomor absen
	teman, found := validasiteman(absen)

	// Tampilkan hasil pencarian
	if found {
		fmt.Println("DATA TEMAN KELAS")
		fmt.Println("No absen", absen, ":")
		fmt.Println("Nama:", teman.Nama)
		fmt.Println("Alamat:", teman.Alamat)
		fmt.Println("Pekerjaan:", teman.Pekerjaan)
		fmt.Println("Alasan memilih kelas Golang:", teman.Alasan)
	} else {
		fmt.Println("Data teman dengan absen", absen, "tidak ditemukan atau tidak ada")
	}

}

type AnggotaKelas struct {
	Absen     int
	Nama      string
	Alamat    string
	Pekerjaan string
	Alasan    string
}

var Anggota = []AnggotaKelas{
	{
		1,
		"Sidiq Trihartanto",
		"Jl. Pondok Ranggon RT07/01",
		"Karyawan Swasta",
		"Karna ingin mempelajari backend dan tertarik terhadap bahasa Go",
	},
	{
		2,
		"Rinaldi Mulya",
		"Jl. Bahagia No.1 Bogor",
		"Karyawan Swasta",
		"Karna pekerjaan backend dan tertarik dengan Go",
	},
	{
		3,
		"Ristiani",
		"Jl. Bahagia No.2 Bogor",
		"Pegawai Negeri Sipil",
		"Sangat Menarik mempelajari bahasa Go",
	},
	{
		4,
		"Hendri Heryanto",
		"Jl. Bahagia No.3 Bogor",
		"Pegawai Negeri Sipil",
		"Mengisi kesibukan dan menambahkan ilmu tentang backend dan bahasa Go",
	},
	{
		5,
		"Gusti Ayu",
		"Jl. Bahagia No.4 Bogor",
		"Guru Matematika",
		"Rasa ingin tahu dengan bahasa Go yang amat tinggi",
	},
	{
		6,
		"Faza Iman",
		"Jl. Bahagia No.5 Bogor",
		"Pegawai BUMN",
		"Saya ingin mengetahui tentang bahasa Go",
	},
}

func validasiteman(absen int) (AnggotaKelas, bool) {
	for _, temanaku := range Anggota {
		if temanaku.Absen == absen {
			return temanaku, true

		}
	}
	return AnggotaKelas{}, false

}
