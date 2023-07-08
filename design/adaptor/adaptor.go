package adaptor

import "fmt"

// 新接口    音乐播放
type MusicPlayer interface {
	play(fileType, fileName string)
}

// 旧接口
type ExistPlayer struct {
}

func (ExistPlayer) playMp3(fileName string) {
	fmt.Println("play mp3: ", fileName)
}

func (ExistPlayer) playWma(fileName string) {
	fmt.Println("play wma: ", fileName)
}

type PlayAdaptor struct {
	ExistPlayer
}

// 实现接口
func (p PlayAdaptor) play(fileType, fileName string) {
	switch fileType {
	case "mp3":
		p.playMp3(fileName)
	case "wma":
		p.playWma(fileName)
	default:
		fmt.Println("暂时不支持此类型的文件播放")
	}
}
