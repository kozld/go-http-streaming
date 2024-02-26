package testdata

import (
	"bufio"
	"io"
	"os"
)

const (
	N = 1000
)

func Generate(filepath string) (*os.File, error) {
	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(file)

	for i := 0; i < N; i += 1 {
		writer.WriteString(`
			«Во всякой гениальной или новой человеческой мысли, 
			или просто даже во всякой серьезной человеческой мы
			сли, зарождающейся в чьей-нибудь голове, всегда ост
			ается нечто такое, чего никак нельзя передать други
			м людям, хотя бы вы исписали целые томы и растолков
			ывали вашу мысль тридцать пять лет; всегда останетс
			я нечто, что ни за что не захочет выйти из-под ваше
			го черепа и останется при вас навеки». — Федор Дост
			оевский, «Идиот»`)
	}

	_, err = file.Seek(0, io.SeekStart)
	return file, err
}
