package adaptor

import "testing"

func TestPlay(t *testing.T) {
	player := PlayAdaptor{}
	player.play("mp3", "晴天")
	player.play("wma", "东风破")
	player.play("mp4", "七里香")
}
