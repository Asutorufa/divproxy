package divproxyinit

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

func AddOneRule(str string, RulePath string) error {
	f, err := os.OpenFile(RulePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer func() {
		if err = f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	if _, err = f.WriteString(str + "\n"); err != nil {
		return err
	}
	return nil
}

func DeleteOneRule(str string, RulePath string) error {
	if str == "" {
		return nil
	}
	f, err := os.Open(RulePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	fBak, err := os.Create(RulePath + ".bak")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer func() {
		_ = fBak.Close()
	}()

	ff := bufio.NewReader(f)
	for {
		line, _, err := ff.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
				break
			}
		}
		if bytes.HasPrefix(line, []byte(str+" ")) {
			continue
		}
		if _, err = fBak.Write(append(line, 10)); err != nil {
			fmt.Println(err)
		}
	}
	if err = f.Close(); err != nil {
		return err
	}
	if err = fBak.Close(); err != nil {
		return err
	}
	if err = os.Rename(RulePath, RulePath+".bak2"); err != nil {
		return err
	}
	if err = os.Rename(RulePath+".bak", RulePath); err != nil {
		return err
	}
	return nil
}
