package words

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"star/internal/dao"
	"star/internal/model"
	"star/internal/model/do"
	"star/internal/model/entity"
)

func Create(ctx context.Context, in *model.WordInput) error {
	if err := checkWord(ctx, 0, in); err != nil {
		return err
	}

	_, err := dao.Words.Ctx(ctx).Data(do.Words{
		Uid:                in.Uid,
		Word:               in.Word,
		Definition:         in.Definition,
		ExampleSentence:    in.ExampleSentence,
		ChineseTranslation: in.ChineseTranslation,
		Pronunciation:      in.Pronunciation,
		ProficiencyLevel:   in.ProficiencyLevel,
	}).Insert()
	if err != nil {
		return err
	}
	return nil
}

func Update(ctx context.Context, id uint, in *model.WordInput) error {
	if err := checkWord(ctx, id, in); err != nil {
		return err
	}
	_, err := dao.Words.Ctx(ctx).Data(do.Words{
		Word:               in.Word,
		Definition:         in.Definition,
		ExampleSentence:    in.ExampleSentence,
		ChineseTranslation: in.ChineseTranslation,
		Pronunciation:      in.Pronunciation,
		ProficiencyLevel:   in.ProficiencyLevel,
	}).Where("id", id).Update()
	if err != nil {
		return err
	}
	return nil
}

func List(ctx context.Context, query *model.WordQuery) (list []entity.Words, total uint, err error) {
	if query == nil {
		query = &model.WordQuery{}
	}
	// 对于查询初始值的处理
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Size == 0 {
		query.Size = 15
	}

	// 组成查询链
	db := dao.Words.Ctx(ctx)
	// 模糊查询
	if len(query.Word) != 0 {
		db = db.Where("word like ?", "%"+query.Word+"%")
	}
	db = db.Order("created_at desc, id desc").Page(query.Page, query.Size)

	data, totalInt, err := db.AllAndCount(true)
	if err != nil {
		return
	}

	list = []entity.Words{}
	_ = data.Structs(&list)
	total = uint(totalInt)

	return
}

func Detail(ctx context.Context, id uint) (word *entity.Words, err error) {
	word = &entity.Words{}
	err = dao.Words.Ctx(ctx).Where("id", id).Scan(&word)
	return
}

func Delete(ctx context.Context, id uint) (err error) {
	_, err = dao.Words.Ctx(ctx).Where("id", id).Delete()
	return
}

// checkWord 在更新时不检查自身
func checkWord(ctx context.Context, id uint, in *model.WordInput) error {
	db := dao.Words.Ctx(ctx).Where("uid", in.Uid).Where("word", in.Word)
	if id > 0 {
		db = db.WhereNot("id", id)
	}
	count, err := db.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("单词已存在")
	}
	return nil
}