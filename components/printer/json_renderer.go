package printer

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/sabudaye/openproject-cli/models"
)

type JsonRenderer struct{}

type relationWorkPackageJson struct {
	Id    uint64 `json:"id"`
	Title string `json:"title,omitempty"`
}

type relationJson struct {
	Id          uint64                  `json:"id"`
	Name        string                  `json:"name"`
	Type        string                  `json:"type"`
	ReverseType string                  `json:"reverse_type,omitempty"`
	Lag         int                     `json:"lag,omitempty"`
	Description string                  `json:"description,omitempty"`
	From        relationWorkPackageJson `json:"from"`
	To          relationWorkPackageJson `json:"to"`
}

func toRelationJson(rel *models.Relation) relationJson {
	return relationJson{
		Id:          rel.Id,
		Name:        rel.Name,
		Type:        rel.Type,
		ReverseType: rel.ReverseType,
		Lag:         rel.Lag,
		Description: rel.Description,
		From:        relationWorkPackageJson{Id: rel.FromId, Title: rel.FromTitle},
		To:          relationWorkPackageJson{Id: rel.ToId, Title: rel.ToTitle},
	}
}

func (r *JsonRenderer) Relation(rel *models.Relation) {
	printJson(toRelationJson(rel))
}

func (r *JsonRenderer) Relations(rels []*models.Relation) {
	out := make([]relationJson, len(rels))
	for i, rel := range rels {
		out[i] = toRelationJson(rel)
	}
	printJson(out)
}

func (r *JsonRenderer) Attachment(a *models.Attachment) {
	printJson(struct {
		Id              uint64 `json:"id"`
		FileName        string `json:"file_name"`
		FileSize        int64  `json:"file_size"`
		Status          string `json:"status"`
		ContentType     string `json:"content_type"`
		DigestAlgorithm string `json:"digest_algorithm,omitempty"`
		DigestHash      string `json:"digest_hash,omitempty"`
		CreatedAt       string `json:"created_at"`
	}{a.Id, a.FileName, a.FileSize, a.Status, a.ContentType, a.DigestAlgo, a.DigestHash, a.CreatedAt})
}

func (r *JsonRenderer) Attachments(as []*models.Attachment) {
	type item struct {
		Id          uint64 `json:"id"`
		FileName    string `json:"file_name"`
		FileSize    int64  `json:"file_size"`
		Status      string `json:"status"`
		ContentType string `json:"content_type"`
		DigestAlgo  string `json:"digest_algorithm,omitempty"`
		DigestHash  string `json:"digest_hash,omitempty"`
		CreatedAt   string `json:"created_at"`
	}
	out := make([]item, len(as))
	for i, a := range as {
		out[i] = item{a.Id, a.FileName, a.FileSize, a.Status, a.ContentType, a.DigestAlgo, a.DigestHash, a.CreatedAt}
	}
	printJson(out)
}

func (r *JsonRenderer) Budget(b *models.Budget) {
	printJson(struct {
		Id      uint64 `json:"id"`
		Subject string `json:"subject"`
	}{b.Id, b.Subject})
}

func (r *JsonRenderer) Budgets(bs []*models.Budget) {
	type item struct {
		Id      uint64 `json:"id"`
		Subject string `json:"subject"`
	}
	out := make([]item, len(bs))
	for i, b := range bs {
		out[i] = item{b.Id, b.Subject}
	}
	printJson(out)
}

func (r *JsonRenderer) WorkPackage(wp *models.WorkPackage) {
	printJson(struct {
		Id          uint64 `json:"id"`
		DisplayId   string `json:"display_id,omitempty"`
		Subject     string `json:"subject"`
		Type        string `json:"type"`
		Status      string `json:"status"`
		Assignee    string `json:"assignee"`
		Description string `json:"description"`
	}{wp.Id, wp.DisplayId, wp.Subject, wp.Type, wp.Status, wp.Assignee, wp.Description})
}

func (r *JsonRenderer) WorkPackages(wps []*models.WorkPackage) {
	type item struct {
		Id        uint64 `json:"id"`
		DisplayId string `json:"display_id,omitempty"`
		Subject   string `json:"subject"`
		Type      string `json:"type"`
		Status    string `json:"status"`
		Assignee  string `json:"assignee"`
	}
	out := make([]item, len(wps))
	for i, wp := range wps {
		out[i] = item{wp.Id, wp.DisplayId, wp.Subject, wp.Type, wp.Status, wp.Assignee}
	}
	printJson(out)
}

func (r *JsonRenderer) Project(p *models.Project) {
	printJson(struct {
		Id         uint64 `json:"id"`
		Identifier string `json:"identifier"`
		Name       string `json:"name"`
	}{p.Id, p.Identifier, p.Name})
}

func (r *JsonRenderer) Projects(ps []*models.Project) {
	type item struct {
		Id         uint64 `json:"id"`
		Identifier string `json:"identifier"`
		Name       string `json:"name"`
	}
	out := make([]item, len(ps))
	for i, p := range ps {
		out[i] = item{p.Id, p.Identifier, p.Name}
	}
	printJson(out)
}

func (r *JsonRenderer) User(u *models.User) {
	printJson(struct {
		Id        uint64 `json:"id"`
		Name      string `json:"name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}{u.Id, u.Name, u.FirstName, u.LastName})
}

func (r *JsonRenderer) Users(us []*models.User) {
	type item struct {
		Id        uint64 `json:"id"`
		Name      string `json:"name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	out := make([]item, len(us))
	for i, u := range us {
		out[i] = item{u.Id, u.Name, u.FirstName, u.LastName}
	}
	printJson(out)
}

func (r *JsonRenderer) Types(types []*models.Type) {
	type item struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
	}
	out := make([]item, len(types))
	for i, t := range types {
		out[i] = item{t.Id, t.Name}
	}
	printJson(out)
}

func (r *JsonRenderer) Status(s *models.Status) {
	printJson(struct {
		Id        uint64 `json:"id"`
		Name      string `json:"name"`
		IsDefault bool   `json:"is_default"`
		IsClosed  bool   `json:"is_closed"`
	}{s.Id, s.Name, s.IsDefault, s.IsClosed})
}

func (r *JsonRenderer) StatusList(statuses []*models.Status) {
	type item struct {
		Id        uint64 `json:"id"`
		Name      string `json:"name"`
		IsDefault bool   `json:"is_default"`
		IsClosed  bool   `json:"is_closed"`
	}
	out := make([]item, len(statuses))
	for i, s := range statuses {
		out[i] = item{s.Id, s.Name, s.IsDefault, s.IsClosed}
	}
	printJson(out)
}

func (r *JsonRenderer) TimeEntry(t *models.TimeEntry) {
	printJson(struct {
		Id          uint64  `json:"id"`
		Comment     string  `json:"comment"`
		Project     string  `json:"project"`
		WorkPackage string  `json:"work_package"`
		SpentOn     string  `json:"spent_on"`
		Hours       float64 `json:"hours"`
		Activity    string  `json:"activity"`
		User        string  `json:"user"`
	}{t.Id, t.Comment, t.Project, t.WorkPackage, t.SpentOn.Format("2006-01-02"), t.Hours.Hours(), t.Activity, t.User})
}

func (r *JsonRenderer) TimeEntryList(entries []*models.TimeEntry) {
	type item struct {
		Id          uint64  `json:"id"`
		Comment     string  `json:"comment"`
		Project     string  `json:"project"`
		WorkPackage string  `json:"work_package"`
		SpentOn     string  `json:"spent_on"`
		Hours       float64 `json:"hours"`
		Activity    string  `json:"activity"`
		User        string  `json:"user"`
	}
	out := make([]item, len(entries))
	for i, t := range entries {
		out[i] = item{
			t.Id,
			t.Comment,
			t.Project,
			t.WorkPackage,
			t.SpentOn.Format("2006-01-02"),
			t.Hours.Hours(),
			t.Activity,
			t.User,
		}
	}
	printJson(out)
}

func (r *JsonRenderer) Notification(n *models.Notification) {
	printJson(struct {
		Id              uint64 `json:"id"`
		ResourceId      uint64 `json:"resource_id"`
		ResourceSubject string `json:"resource_subject"`
		Reason          string `json:"reason"`
		Read            bool   `json:"read"`
	}{n.Id, n.ResourceId, n.ResourceSubject, n.Reason, n.Read})
}

func (r *JsonRenderer) Notifications(ns []*models.Notification) {
	type item struct {
		Id              uint64 `json:"id"`
		ResourceId      uint64 `json:"resource_id"`
		ResourceSubject string `json:"resource_subject"`
		Reason          string `json:"reason"`
		Read            bool   `json:"read"`
	}
	out := make([]item, len(ns))
	for i, n := range ns {
		out[i] = item{n.Id, n.ResourceId, n.ResourceSubject, n.Reason, n.Read}
	}
	printJson(out)
}

func (r *JsonRenderer) Activities(activities []*models.Activity, users []*models.User) {
	type item struct {
		Id        uint64    `json:"id"`
		Comment   string    `json:"comment"`
		Details   []*string `json:"details"`
		UpdatedAt string    `json:"updated_at"`
		User      string    `json:"user"`
	}
	out := make([]item, len(activities))
	for i, a := range activities {
		userName := ""
		if a.UserId > 0 {
			idx := sort.Search(len(users), func(j int) bool { return users[j].Id >= a.UserId })
			if idx < len(users) && users[idx].Id == a.UserId {
				userName = users[idx].Name
			}
		}
		out[i] = item{a.Id, a.Comment, a.Details, a.UpdatedAt, userName}
	}
	printJson(out)
}

func (r *JsonRenderer) CustomActions(actions []*models.CustomAction) {
	type item struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
	}
	out := make([]item, len(actions))
	for i, a := range actions {
		out[i] = item{a.Id, a.Name}
	}
	printJson(out)
}

type whoamiJson struct {
	Profile string `json:"profile"`
	Server  string `json:"server"`
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
}

func (r *JsonRenderer) Whoami(profile, host string, user *models.User) {
	printJson(whoamiJson{profile, host, user.Id, user.Name})
}

func (r *JsonRenderer) WhoamiList(entries []WhoamiEntry) {
	out := make([]whoamiJson, len(entries))
	for i, entry := range entries {
		out[i] = whoamiJson{entry.Profile, entry.Host, entry.User.Id, entry.User.Name}
	}
	printJson(out)
}

func (r *JsonRenderer) Number(n int64) {
	printJson(struct {
		Total int64 `json:"total"`
	}{n})
}

func printJson(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		// Marshal the error payload too, so the failure output is still valid
		// JSON even when the error text contains quotes or newlines.
		errJson, _ := json.Marshal(struct {
			Error string `json:"error"`
		}{fmt.Sprintf("failed to serialize output: %s", err)})
		activePrinter.Println(string(errJson))
		return
	}
	activePrinter.Println(string(b))
}
