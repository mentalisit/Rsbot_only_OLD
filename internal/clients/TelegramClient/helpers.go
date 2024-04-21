package TelegramClient

import (
	"strconv"
	"strings"
)

func (t *Telegram) chat(chatid string) (chatId int64, ThreadID int) {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	if len(a) > 1 {
		ThreadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
	return chatId, ThreadID
}
