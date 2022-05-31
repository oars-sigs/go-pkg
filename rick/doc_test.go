package rick

import "testing"

func TestPutDoc(t *testing.T) {
	cfg := &Config{
		Addr: "https://app.oars.gzsunrun.cn",
	}
	cli := New(cfg)
	err := cli.PutDoc("default", &PutDocRequestParam{
		FileId:   "4e97d98a-7c2b-4a67-815d-d6449993e0fe",
		FileName: "超级邮编(第一期)开发方案V0.6.docx",
		Tags:     []string{"namespace::212e391d-adc2-45d6-a10b-9a8ab77b7de3"},
	})
	if err != nil {
		t.Error(err)
	}
}
