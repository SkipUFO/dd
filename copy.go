package dd

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// FastCopy copies fast and without output to console
func FastCopy(source string, destination string, offset int64, limit int64) error {
	// Проверяем доступность source на чтение и права на запись в destination
	src, err := os.Open(source)
	defer src.Close()
	if err != nil {
		return err
	}

	dst, err := os.Create(destination)
	defer dst.Close()
	if err != nil {
		return err
	}

	// Сдвигаемся по файлу
	if _, err := src.Seek(offset, 0); err != nil {
		return err
	}

	if limit == 0 {
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
	} else {
		if _, err := io.CopyN(dst, src, limit); err != nil {
			return err
		}
	}

	return nil
}

// Copy copies slower than FastCopy, but prints progress
func Copy(source string, destination string, offset int64, limit int64, bufferSize int64) error {
	// Проверяем доступность source на чтение и права на запись в destination
	src, err := os.Open(source)
	defer src.Close()
	if err != nil {
		return err
	}

	dst, err := os.Create(destination)
	defer dst.Close()
	if err != nil {
		return err
	}

	// Сдвигаемся по файлу
	if _, err := src.Seek(offset, 0); err != nil {
		return err
	}

	fi, err := src.Stat()
	if err != nil {
		return err
	}

	size := fi.Size() - offset

	var reader io.Reader = src
	var writer io.Writer = dst

	return copy(&reader, &writer, limit, bufferSize, size)

	/*

		bw := bufio.NewWriter(dst)

		err = nil

		for err != io.EOF {

			buf := make([]byte, bufferSize)
			var read = 0
			// Здесь пропускаем UnexpectedEOF, т.к. ReadFull её возвращает, когда
			if read, err = io.ReadFull(src, buf); err != nil && err != io.ErrUnexpectedEOF {
				return err
			}

			// Костыль, для того, чтобы выйти из цикла, по-хорошему надо было бы использовать чтение в цикле, но использовал ReadFull
			if err == io.ErrUnexpectedEOF {
				err = io.EOF
			}

			if _, errWrite := bw.Write(buf); errWrite != nil {
				fmt.Println(2)
				return errWrite
			}
			bw.Flush()
			copied += int64(read)
			fmt.Printf("copied %d of %d bytes\n", copied, size)
		}

		if limit == 0 {
			if _, err := io.Copy(dst, src); err != nil {
				return err
			}
		} else {
			if _, err := io.CopyN(dst, src, limit); err != nil {
				return err
			}
		}
	*/
}

func copy(reader *io.Reader, writer *io.Writer, limit int64, bufferSize int64, readerSize int64) error {
	bw := bufio.NewWriter(*writer)

	var err error
	var copied int64

	for err != io.EOF && ((copied < limit) || (limit == 0)) {

		var tempSize int64
		if limit != 0 {
			if copied+bufferSize > limit {
				tempSize = copied + bufferSize - limit
			} else {
				tempSize = bufferSize
			}
		} else {
			tempSize = bufferSize
		}
		buf := make([]byte, tempSize)
		var read = 0
		// Здесь пропускаем UnexpectedEOF, т.к. ReadFull её возвращает, когда
		if read, err = io.ReadFull(*reader, buf); err != nil && err != io.ErrUnexpectedEOF {
			return err
		}

		// Костыль, для того, чтобы выйти из цикла, по-хорошему надо было бы использовать чтение в цикле, но использовал ReadFull
		if err == io.ErrUnexpectedEOF {
			err = io.EOF
		}

		if _, errWrite := bw.Write(buf); errWrite != nil {
			return errWrite
		}
		bw.Flush()
		copied += int64(read)
		fmt.Printf("copied %d of %d bytes\n", copied, readerSize)
	}

	// if limit == 0 {
	// 	if _, err := io.Copy(*writer, *reader); err != nil {
	// 		return err
	// 	}
	// } else {
	// 	fmt.Println(limit)
	// 	if _, err := io.CopyN(*writer, *reader, limit); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
