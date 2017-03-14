package main

import (
	// "log"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/kataras/iris.v6"
)

type User struct {
	gorm.Model
	Email    string `sql:"size:255" unique_index`
	Password string `sql:"size:255"`
	Verified bool   `sql:"not null"`
}

func (u *User) getDefaultWebsite() *Website {
	w := &Website{OwnerID: u.ID, Default: true}
	db.Where(w).First(w)
	return w
}

func (u *User) redirectToDefaultWebsite(ctx *iris.Context) {
	w := u.getDefaultWebsite()
	var redirectUrl string
	if w.ID == 0 {
		redirectUrl = conf.AppUrl + APP_PATH + "/websites/new"
	} else {
		redirectUrl = conf.AppUrl + APP_PATH + "/" + strconv.Itoa(int(w.ID))
	}
	ctx.Redirect(redirectUrl)
	return
}

type Website struct {
	gorm.Model
	Owner   *User
	OwnerID uint   `sql:"index"`
	Name    string `sql:"size:255"`
	Url     string `sql:"size:255"`
	Default bool   `sql:"not null"`
}

func (w *Website) CountPageViews() int {
	count := 0
	var pvs []*PageView
	weekAgo := time.Now().Truncate(time.Hour).Add(-time.Hour*time.Duration(time.Now().Hour())).AddDate(0, 0, -7)
	db.Where("website_id = ? AND created_at BETWEEN ? and ?", w.ID, weekAgo, time.Now()).Find(&pvs).Count(&count)
	return count
}

func (w *Website) CountUsers() int {
	// count := 0
	counter := map[uint]bool{}
	countPvs := 0
	var pvs []*PageView
	weekAgo := time.Now().Truncate(time.Hour).Add(-time.Hour*time.Duration(time.Now().Hour())).AddDate(0, 0, -7)
	db.Where("website_id = ? AND created_at BETWEEN ? and ?", w.ID, weekAgo, time.Now()).Find(&pvs).Count(&countPvs)
	for _, pv := range pvs {
		counter[pv.VisitorID] = true
	}
	return len(counter)
}

func (w *Website) CountVisits() int {
	count := 0
	// var pvs []*Visitor
	// weekAgo := time.Now().Truncate(time.Hour).Add(-time.Hour*time.Duration(time.Now().Hour())).AddDate(0, 0, -7)
	// db.Where("website_id = ? AND created_at BETWEEN ? and ?", w.ID, weekAgo, time.Now()).Find(&pvs).Count(&count)
	return count
}

type Visitor struct {
	gorm.Model
	IpAddress  string `sql:"size:255"`
	Resolution string `sql:"size:255"`
	Language   string `sql:"size:255"`
}

type Page struct {
	gorm.Model
	Hostname string `sql:"size:255"`
	Path     string `sql:"size:255"`
	Title    string `sql:"size:255"`
}

type PageView struct {
	gorm.Model
	Website   *Website
	WebsiteID uint `sql:"index"`
	Visitor   *Visitor
	VisitorID uint `sql:"index"`
	Page      *Page
	PageID    uint `sql:"index"`
}
