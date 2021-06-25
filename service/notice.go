package service

/*func InsertNoticesByRepo(repo *model.TRepo, n *notice.Notice) (rterr error) {
	return InsertNoticesByRepoId(repo.Id, n)
}
func InsertNoticesByRepoOpenid(openid string, n *notice.Notice) (rterr error) {
	db := comm.DBMain.GetDB()
	repo := &model.TRepo{}
	get, err := db.Where("openid =? ", openid).Where("deleted != 1").Get(repo)
	if err != nil {
		return err
	}
	if !get {
		return nil
	}
	return InsertNoticesByRepoId(repo.Id, n)
}
func InsertNoticesByRepoId(id string, n *notice.Notice) (rterr error) {
	db := comm.DBMain.GetDB()
	turs := make([]*model.TUserRepo, 0)
	err := db.Where("repo_id = ? ", id).Find(&turs)
	if err != nil {
		return err
	}
	message, err := InsertTMessage(n)
	if err != nil {
		return err
	}
	for _, tur := range turs {
		_, err = InsertTUserMsg(message, strconv.FormatInt(tur.UserId, 10))
		if err != nil {
			core.Log.Errorf("InsertTUserMsg db err: %v", err)
			continue
		}
	}
	return nil
}
func InsertTUserMsg(m *model.TMessage, sender string) (*model.TUserMsg, error) {
	db := comm.DBMain.GetDB()
	tum := &model.TUserMsg{
		Mid:     m.Xid,
		Uid:     sender,
		Created: time.Now(),
		Status:  0,
		Deleted: 0,
	}
	_, err := db.Insert(tum)
	if err != nil {
		return nil, err
	}
	return tum, nil
}
func InsertTMessage(n *notice.Notice) (*model.TMessage, error) {
	db := comm.DBMain.GetDB()
	msg := &model.TMessage{
		Xid:     utils.NewXid(),
		Uid:     strconv.FormatInt(n.Receiver, 10),
		Title:   n.Title,
		Content: n.Content,
		Types:   n.Types,
		Created: time.Now(),
		Infos:   n.Infos,
		Url:     n.Url,
	}
	_, err := db.Insert(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}*/
