package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Handle struct {
}

type Req struct {
	Id              string      `json:"id"`
	Project         string      `json:"project"`
	ProjectName     string      `json:"project_name"`
	ProjectSlug     string      `json:"project_slug"`
	Logger          interface{} `json:"logger"`
	Level           string      `json:"level"`
	Culprit         string      `json:"culprit"`
	Message         string      `json:"message"`
	Url             string      `json:"url"`
	TriggeringRules []string    `json:"triggering_rules"`
	Event           struct {
		EventId     string   `json:"event_id"`
		Level       string   `json:"level"`
		Version     string   `json:"version"`
		Type        string   `json:"type"`
		Fingerprint []string `json:"fingerprint"`
		Culprit     string   `json:"culprit"`
		Logentry    struct {
			Formatted string `json:"formatted"`
		} `json:"logentry"`
		Logger    string            `json:"logger"`
		Modules   map[string]string `json:"modules"`
		Platform  string            `json:"platform"`
		Timestamp float64           `json:"timestamp"`
		Received  float64           `json:"received"`
		Release   string            `json:"release"`
		User      struct {
			Id        string `json:"id"`
			IpAddress string `json:"ip_address"`
		} `json:"user"`
		Contexts struct {
			Device struct {
				Arch   string `json:"arch"`
				NumCpu int    `json:"num_cpu"`
				Type   string `json:"type"`
			} `json:"device"`
			Os struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"os"`
			Runtime struct {
				Name          string `json:"name"`
				Version       string `json:"version"`
				GoMaxprocs    int    `json:"go_maxprocs"`
				GoNumcgocalls int    `json:"go_numcgocalls"`
				GoNumroutines int    `json:"go_numroutines"`
				Type          string `json:"type"`
			} `json:"runtime"`
		} `json:"contexts"`
		Tags [][]string `json:"tags"`
		Sdk  struct {
			Name         string   `json:"name"`
			Version      string   `json:"version"`
			Integrations []string `json:"integrations"`
			Packages     []struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"packages"`
		} `json:"sdk"`
		Errors []struct {
			Type   string `json:"type"`
			Name   string `json:"name"`
			Reason string `json:"reason"`
		} `json:"errors"`
		KeyId          string `json:"key_id"`
		Project        int    `json:"project"`
		GroupingConfig struct {
			Enhancements string `json:"enhancements"`
			Id           string `json:"id"`
		} `json:"grouping_config"`
		Metrics struct {
			BytesIngestedEvent int `json:"bytes.ingested.event"`
			BytesStoredEvent   int `json:"bytes.stored.event"`
		} `json:"_metrics"`
		Ref        int      `json:"_ref"`
		RefVersion int      `json:"_ref_version"`
		Hashes     []string `json:"hashes"`
		Metadata   struct {
			Title string `json:"title"`
		} `json:"metadata"`
		NodestoreInsert float64     `json:"nodestore_insert"`
		Title           string      `json:"title"`
		Location        string      `json:"location"`
		Meta            interface{} `json:"_meta"`
		Id              string      `json:"id"`
	} `json:"event"`
}

var (
	Token *string
)

func init() {

}

func ServeHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cover := r.Body
		if cover == nil {
			w.Write(returnMsg("body is null"))
			return
		}
		body, err := io.ReadAll(cover)
		var req *Req
		err = json.Unmarshal(body, &req)
		if err != nil {
			w.Write(returnMsg("json unmarshal fail:%v", err))
			return
		}
		HandleHook(r, req)
		w.Write([]byte("ok"))
	}
}

func HandleHook(r *http.Request, req *Req) error {
	param := map[string]string{
		"project_name": req.ProjectName,
		"culprit":      req.Culprit,
		"level":        req.Level,
		"message":      req.Message,
		"title":        req.Event.Title,
		"location":     req.Event.Location,
		"url":          req.Url,
	}
	query := r.URL.Query()
	token := query.Get("token")
	fmt.Println("token", query.Get("token"))
	content := "- 项目：@{project_name}\n- level：@{level}\n## 内容：\n```\n@{title}\n@{culprit}\n@{location}\n```\n- [查看更多](@{url})\n"
	content = ReplaceParam(content, param)
	fmt.Println("content", content)
	NewTalkRobot(token).Markdown("error", content).Send(false)
	return nil
}

func ReplaceParam(str string, arr map[string]string) string {
	re, _ := regexp.Compile("@{(\\S+?)}")
	rep := re.ReplaceAllStringFunc(str, func(s string) string {
		key := s[2 : len(s)-1]
		if arr[key] != "" {
			return arr[key]
		}
		return s
	})
	return rep
}

func returnMsg(format string, a ...interface{}) []byte {
	return []byte(fmt.Sprintf(format, a...))
}
