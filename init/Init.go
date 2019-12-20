package divproxyinit

import (
	"divproxy/config"
	"os"
)

// PathExists 判断目录是否存在返回布尔类型
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func InitConfig() error {
	if !PathExists(GetConfigPath()) {
		if err := config.InitJSON(GetConfigPath()); err != nil {
            return err
		}
	}

    if !PathExists(GetRuleFilePath()){
        f,err := os.Create(GetRuleFilePath())
        if err != nil{
            return err
        }
        if _,err := f.WriteString("www.example.com test"); err != nil{
            return err
        }
        if err := f.Close(); err != nil{
            return err
        }
    }
    return nil
}
