package classify

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"strings"
)

type ClassifyService struct {
	app.Service

	Create      *ClassifyCreateTask
	BatchCreate *ClassifyBatchCreateTask
	Set         *ClassifySetTask
	Get         *ClassifyTask
	Remove      *ClassifyRemoveTask
	Query       *ClassifyQueryTask
}

func (S *ClassifyService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *ClassifyService) HandleClassifyCreateTask(a IClassifyApp, task *ClassifyCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CLASSIFY
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Classify{}

	v.Name = task.Name
	v.Pid = task.Pid
	v.Alias = task.Alias
	v.Oid = NewOid()
	v.Tags = task.Tags
	v.Path = "/"

	if v.Pid != 0 {

		t := ClassifyTask{}
		t.Id = v.Pid

		app.Handle(a, &t)

		if t.Result.Classify == nil {
			task.Result.Errno = ERROR_CLASSIFY_NOT_FOUND_PCLASSIFY
			task.Result.Errmsg = "Not found parent classify"
			return nil
		}

		v.Path = fmt.Sprintf("%s%d/", t.Result.Classify.Path, v.Pid)
		v.Alias = v.Alias

	}

	_, err = kk.DBInsert(db, a.GetClassifyTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_CLASSIFY
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Classify = &v

	return nil
}

func (S *ClassifyService) HandleClassifyBatchCreateTask(a IClassifyApp, task *ClassifyBatchCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CLASSIFY
		task.Result.Errmsg = err.Error()
		return nil
	}

	classifys := []Classify{}

	var path string = "/"
	var alias = task.Alias

	if task.Pid != 0 {

		t := ClassifyTask{}
		t.Id = task.Pid

		app.Handle(a, &t)

		if t.Result.Classify == nil {
			task.Result.Errno = ERROR_CLASSIFY_NOT_FOUND_PCLASSIFY
			task.Result.Errmsg = "Not found parent classify"
			return nil
		}

		path = fmt.Sprintf("%s%d/", t.Result.Classify.Path, t.Result.Classify.Id)
		alias = t.Result.Classify.Alias

	}

	tx, err := db.Begin()

	if err != nil {
		task.Result.Errno = ERROR_CLASSIFY
		task.Result.Errmsg = err.Error()
		return nil
	}

	err = func() error {

		for _, name := range strings.Split(task.Names, ",") {

			if name != "" {

				v := Classify{}

				v.Name = name
				v.Alias = alias
				v.Pid = task.Pid
				v.Path = path

				_, err = kk.DBInsert(db, a.GetClassifyTable(), a.GetPrefix(), &v)

				if err != nil {
					return err
				}
			}

		}

		return nil

	}()

	if err == nil {
		err = tx.Commit()
	}

	if err != nil {
		tx.Rollback()
		e, ok := err.(*app.Error)
		if ok {
			task.Result.Errno = e.Errno
			task.Result.Errmsg = e.Errmsg
			return nil
		} else {
			task.Result.Errno = ERROR_CLASSIFY
			task.Result.Errmsg = err.Error()
			return nil
		}
	}

	task.Result.Classifys = classifys

	return nil
}

func (S *ClassifyService) HandleClassifySetTask(a IClassifyApp, task *ClassifySetTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_CLASSIFY_NOT_FOUND_ID
		task.Result.Errmsg = "Not found classify id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CLASSIFY
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Classify{}

	rows, err := kk.DBQuery(db, a.GetClassifyTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_CLASSIFY
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_CLASSIFY
			task.Result.Errmsg = err.Error()
			return nil
		}

		keys := map[string]bool{}

		if task.Name != nil {
			v.Name = dynamic.StringValue(task.Name, v.Name)
			keys["name"] = true
		}

		if task.Pid != nil {
			var pid = dynamic.IntValue(task.Pid, 0)

			v.Pid = pid

			if pid != 0 {

				t := ClassifyTask{}
				t.Id = v.Pid

				app.Handle(a, &t)

				if t.Result.Classify == nil {
					task.Result.Errno = ERROR_CLASSIFY_NOT_FOUND_PCLASSIFY
					task.Result.Errmsg = "Not found parent classify"
					return nil
				}

				v.Path = fmt.Sprintf("%s%d/", t.Result.Classify.Path, v.Pid)
				v.Alias = t.Result.Classify.Alias

			} else {
				v.Path = "/"
				v.Alias = task.Alias
			}

			keys["pid"] = true
			keys["path"] = true
			keys["alias"] = true
		}

		if task.Tags != nil {
			v.Tags = dynamic.StringValue(task.Tags, v.Tags)
			keys["tags"] = true
		}

		_, err = kk.DBUpdateWithKeys(db, a.GetClassifyTable(), a.GetPrefix(), &v, keys)

		if err != nil {
			task.Result.Errno = ERROR_CLASSIFY
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Classify = &v

	} else {
		task.Result.Errno = ERROR_CLASSIFY_NOT_FOUND
		task.Result.Errmsg = "Not found classify"
		return nil
	}

	return nil
}

func (S *ClassifyService) HandleClassifyRemoveTask(a IClassifyApp, task *ClassifyRemoveTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CLASSIFY
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.Id != 0 {

		rows, err := db.Query(fmt.Sprintf("SELECT path FROM %s%s WHERE id=?", a.GetPrefix(), a.GetClassifyTable().Name), task.Id)

		if err != nil {
			task.Result.Errno = ERROR_CLASSIFY
			task.Result.Errmsg = err.Error()
			return nil
		}

		defer rows.Close()

		if rows.Next() {

			var path interface{} = nil

			err = rows.Scan(&path)

			if err != nil {
				task.Result.Errno = ERROR_CLASSIFY
				task.Result.Errmsg = err.Error()
				return nil
			}

			_, err = kk.DBDelete(db, a.GetClassifyTable(), a.GetPrefix(), " WHERE id=? OR path LIKE ?", task.Id, fmt.Sprintf("%s%d/%%", dynamic.StringValue(path, "/"), task.Id))

			if err != nil {
				task.Result.Errno = ERROR_CLASSIFY
				task.Result.Errmsg = err.Error()
				return nil
			}

		} else {
			task.Result.Errno = ERROR_CLASSIFY_NOT_FOUND
			task.Result.Errmsg = "Not found classify"
			return nil
		}

	} else if task.Pid != nil {

		var pid = dynamic.IntValue(task.Pid, 0)

		if pid == 0 {
			_, err = kk.DBDelete(db, a.GetClassifyTable(), a.GetPrefix(), " WHERE pid=? OR path LIKE '/%'", pid)
			if err != nil {
				task.Result.Errno = ERROR_CLASSIFY
				task.Result.Errmsg = err.Error()
				return nil
			}
		} else {

			rows, err := db.Query(fmt.Sprintf("SELECT path FROM %s%s WHERE id=?", a.GetPrefix(), a.GetClassifyTable().Name), pid)

			if err != nil {
				task.Result.Errno = ERROR_CLASSIFY
				task.Result.Errmsg = err.Error()
				return nil
			}

			defer rows.Close()

			if rows.Next() {

				var path interface{} = nil

				err = rows.Scan(&path)

				if err != nil {
					task.Result.Errno = ERROR_CLASSIFY
					task.Result.Errmsg = err.Error()
					return nil
				}

				_, err = kk.DBDelete(db, a.GetClassifyTable(), a.GetPrefix(), " WHERE path LIKE ?", fmt.Sprintf("%s%d/%%", dynamic.StringValue(path, "/"), pid))

				if err != nil {
					task.Result.Errno = ERROR_CLASSIFY
					task.Result.Errmsg = err.Error()
					return nil
				}

			} else {
				task.Result.Errno = ERROR_CLASSIFY_NOT_FOUND_PCLASSIFY
				task.Result.Errmsg = "Not found parent classify"
				return nil
			}
		}

	}

	return nil
}

func (S *ClassifyService) HandleClassifyTask(a IClassifyApp, task *ClassifyTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CLASSIFY
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Classify{}

	rows, err := kk.DBQuery(db, a.GetClassifyTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_CLASSIFY
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_CLASSIFY
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Classify = &v

	} else {
		task.Result.Errno = ERROR_CLASSIFY_NOT_FOUND
		task.Result.Errmsg = "Not found classify"
		return nil
	}

	return nil
}

func (S *ClassifyService) HandleClassifyQueryTask(a IClassifyApp, task *ClassifyQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CLASSIFY
		task.Result.Errmsg = err.Error()
		return nil
	}

	var classifys = []Classify{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(" WHERE 1")

	if task.Alias != "" {
		sql.WriteString(" AND alias=?")
		args = append(args, task.Alias)
	}

	if task.Id != 0 {
		sql.WriteString(" AND id=?")
		args = append(args, task.Id)
	}

	if task.Pid != nil {
		sql.WriteString(" AND pid=?")
		args = append(args, task.Pid)
	}

	if task.Keyword != "" {
		q := "%" + task.Keyword + "%"
		sql.WriteString(" AND (tags LIKE ? OR name LIKE ?)")
		args = append(args, q, q)
	}

	sql.WriteString(" ORDER BY oid ASC,id ASC")

	var pageIndex = task.PageIndex
	var pageSize = task.PageSize

	if pageIndex < 1 {
		pageIndex = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	if task.Counter {

		var counter = ClassifyQueryCounter{}
		counter.PageIndex = pageIndex
		counter.PageSize = pageSize
		counter.RowCount, err = kk.DBQueryCount(db, a.GetClassifyTable(), a.GetPrefix(), sql.String(), args...)

		if counter.RowCount%pageSize == 0 {
			counter.PageCount = counter.RowCount / pageSize
		} else {
			counter.PageCount = counter.RowCount/pageSize + 1
		}

		task.Result.Counter = &counter
	}

	sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

	var v = Classify{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetClassifyTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_CLASSIFY
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_CLASSIFY
			task.Result.Errmsg = err.Error()
			return nil
		}

		classifys = append(classifys, v)
	}

	task.Result.Classifys = classifys

	return nil
}
