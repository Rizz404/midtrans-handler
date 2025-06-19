package utils

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// * EpochTime adalah custom type yang membungkus time.Time
// * untuk menangani timestamp integer (milliseconds since epoch) dari JSON.
type EpochTime time.Time

// * UnmarshalJSON adalah method yang membuat EpochTime memenuhi interface json.Unmarshaler.
// * Ini adalah jantung dari solusi kita.
func (t *EpochTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	// * Coba parsing sebagai integer (milliseconds) terlebih dahulu.
	ms, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		// * Jika gagal, coba parsing sebagai string format waktu standar (RFC3339).
		// * Ini membuat tipe data kita lebih robust, bisa menerima kedua format.
		parsedTime, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return err
		}
		*t = EpochTime(parsedTime)
		return nil
	}

	// * Konversi milliseconds ke object time.Time
	// * Go's time.Unix butuh seconds dan nanoseconds, jadi kita konversi.
	*t = EpochTime(time.Unix(ms/1000, (ms%1000)*1000000))
	return nil
}

// * MarshalJSON adalah method yang membuat EpochTime memenuhi interface json.Marshaler.
// * Ini mengontrol bagaimana waktu akan di-encode KEMBALI ke JSON.
// * Kita akan konsisten mengembalikannya sebagai milliseconds.
func (t EpochTime) MarshalJSON() ([]byte, error) {
	// * Casting `t` kembali ke `time.Time` untuk mengakses methodnya.
	ms := time.Time(t).UnixNano() / int64(time.Millisecond)
	return json.Marshal(ms)
}

// * Time mengembalikan underlying time.Time object.
func (t EpochTime) Time() time.Time {
	return time.Time(t)
}
