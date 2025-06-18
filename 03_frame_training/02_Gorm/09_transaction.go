package main

import (
	"errors"
	"gorm.io/gorm"
)

func main() {

}

func GormExecuteSQLWithTransaction() {
	nativeDb.Transaction(func(tx *gorm.DB) error {
		/*如果没有err则提交事务*/
		if err := tx.Create(&User{Name: "Tom", Age: 18}).Error; err != nil {
			return err
		}
		if err := tx.Create(&Employee{Name: "Jerry", Gender: 1}).Error; err != nil {
			return err
		}
		return nil
	})
}

func GormExecuteSQLWithNestTransaction() {
	nativeDb.Transaction(func(tx *gorm.DB) error {
		tx.Create(&Employee{Name: "Tom", Gender: 1})
		tx.Transaction(func(tx1 *gorm.DB) error {
			if err := tx1.Create(&User{Name: "Tom", Age: 18}).Error; err != nil {
				return errors.New("only rollback tx1  transaction")
			}
			return nil
		})
		tx.Transaction(func(tx2 *gorm.DB) error {
			if err := tx2.Create(&User{Name: "Tom", Age: 18}).Error; err != nil {
				return err
			}
			return nil
		})
		return nil
	})
}

func GormExecuteSQLWithOperationTransaction() error {
	tx := nativeDb.Begin()
	defer func() {
		/*这里只会处理panic,数据库的err是非panic的所以不会被捕获*/
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&User{Name: "Tom", Age: 18}).Error; err != nil {
		tx.Rollback()
		return err
	}
	/*这里提交出现错误直接返回错误*/
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func GormExecuteSQLWithSetSavePoint() {
	/*这里应该会提交第一个事务*/
	tx := nativeDb.Begin()
	tx.Create(&User{Name: "Tom", Age: 18})
	tx.SavePoint("c1")
	tx.Create(&User{Name: "Tom", Age: 18})
	tx.RollbackTo("c1")
	tx.Commit()
}
