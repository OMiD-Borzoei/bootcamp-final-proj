package repositories

import (
	"Project/config"
	"Project/models"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func stamp(code string) (string, string) {
	code = fmt.Sprintf("%s-%d", code, time.Now().UnixNano())
	return code, code + "-T"
}

// Create 2 SLs, 1 DL, 1 Voucher with 2 VIs:
// vi-1: SL1, DL, 100, 0
// vi-2: SL2, nil, 0, 100
func scenario1(db *gorm.DB) (sl1_id, sl2_id, dl_id, v_id uint) {
	vRepo := NewVoucherRepository(db)
	slRepo := NewSLRepository(db)
	dlRepo := NewDLRepository(db)

	dlRepo.db.AutoMigrate(&models.DL{}, &models.SL{}, &models.Voucher{}, &models.Voucheritem{})

	code, title := stamp("S010")
	sl1_id, _ = slRepo.Create(code, title, true)

	code, title = stamp("S011")
	sl2_id, _ = slRepo.Create(code, title, false)

	code, title = stamp("D010")
	dl_id, _ = dlRepo.Create(code, title)

	vi1 := models.Voucheritem{
		SLID:   sl1_id,
		DLID:   &dl_id,
		Debit:  100,
		Credit: 0,
	}
	vi2 := models.Voucheritem{
		SLID:   sl2_id,
		DLID:   nil,
		Debit:  0,
		Credit: 100,
	}

	number, _ := stamp("V010")
	v := &models.Voucher{Number: number}
	v.Items = append(v.Items, vi1)
	v.Items = append(v.Items, vi2)

	v_id, _ = vRepo.Create(v)
	return sl1_id, sl2_id, dl_id, v_id
}

func revert_scenario1(db *gorm.DB, sl1_id, sl2_id, dl_id, v_id uint) {
	vRepo := NewVoucherRepository(db)
	slRepo := NewSLRepository(db)
	dlRepo := NewDLRepository(db)

	vRepo.Delete(v_id)
	slRepo.Delete(sl1_id)
	slRepo.Delete(sl2_id)
	dlRepo.Delete(dl_id)
}

func TestAllDLmethods(t *testing.T) {

	// Set up database connection
	godotenv.Load()
	db, err := config.SetupDB()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	repo := NewDLRepository(db)

	repo.db.AutoMigrate(&models.DL{})

	var code, title string

	t.Run("DL CRUD", func(t *testing.T) {
		t.Run("Create", func(t *testing.T) {

			t.Run("Should not Add DL with empty code", func(t *testing.T) {
				code, title = "", fmt.Sprintf("%d", time.Now().UnixNano())
				_, err := repo.Create(code, title)
				assert.NotNil(t, err)
			})

			t.Run("Should not Add DL with empty title", func(t *testing.T) {
				code, title = fmt.Sprintf("%d", time.Now().UnixNano()), ""
				_, err := repo.Create(code, title)
				assert.NotNil(t, err)
			})

			t.Run("Should not Add DL with a code that has more that 64 chars", func(t *testing.T) {
				code, title = fmt.Sprintf("%d", time.Now().UnixNano()), "acceptable title"
				for range 65 {
					code += "a"
				}
				_, err := repo.Create(code, title)
				assert.NotNil(t, err)
			})

			t.Run("Should not Add DL with a title that has more that 64 chars", func(t *testing.T) {
				code, title = "acceptable code", fmt.Sprintf("%d", time.Now().UnixNano())
				for range 65 {
					title += "a"
				}
				_, err := repo.Create(code, title)
				assert.NotNil(t, err)
			})

			t.Run("Should be Able to add a DL with code that has 64 persian chars", func(t *testing.T) {
				code, title = stamp("")
				for range 64 - len(code) {
					code += "م"
				}
				_, err := repo.Create(code, title)
				assert.Nil(t, err)

				// Check for correct insertion in database
				dl := &models.DL{}
				repo.db.Model(&models.DL{}).Last(dl)
				assert.Equal(t, dl.Code, code)
				assert.Equal(t, dl.Title, title)
			})

			t.Run("Should be Able to add a DL with title that has 64 persian chars", func(t *testing.T) {
				code, title = stamp("")
				for range 64 - len(title) {
					title += "م"
				}
				_, err := repo.Create(code, title)
				assert.Nil(t, err)

				// Check for correct insertion in database
				dl := &models.DL{}
				repo.db.Model(&models.DL{}).Last(dl)
				assert.Equal(t, dl.Code, code)
				assert.Equal(t, dl.Title, title)

			})

			t.Run("Shouldn't Add a DL with repetative code", func(t *testing.T) {
				code, title = stamp("D001")
				_, err := repo.Create(code, title)
				assert.Nil(t, err)

				_, err = repo.Create(code, title+"new")
				assert.NotNil(t, err)
			})

			t.Run("Shouldn't Add a DL with repetative title", func(t *testing.T) {
				code, title = stamp("D002")
				_, err := repo.Create(code, title)
				assert.Nil(t, err)

				_, err = repo.Create(code+"new", title)
				assert.NotNil(t, err)
			})

		})

		t.Run("Read", func(t *testing.T) {
			dl := &models.DL{}
			if err := repo.db.Model(&models.DL{}).First(dl).Error; err != nil {
				fmt.Println("DB issue, aborting...")
				return
			}

			readdl, err := repo.Read(dl.ID)
			assert.Nil(t, err)
			assert.Equal(t, dl, readdl)
		})

		t.Run("Update", func(t *testing.T) {
			dl := &models.DL{}
			if err := repo.db.Model(&models.DL{}).Last(dl).Error; err != nil {
				fmt.Println("DB issue, aborting...")
				return
			}

			t.Run("Shouldnt Update DL Code to empty string", func(t *testing.T) {
				newdl := &models.DL{
					Code:  "",
					Title: "Valid" + fmt.Sprintf("%d", time.Now().UnixNano()),
				}

				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update DL Title to empty string", func(t *testing.T) {
				newdl := &models.DL{
					Code:  "Valid" + fmt.Sprintf("%d", time.Now().UnixNano()),
					Title: "",
				}

				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update DL code to string with more than 64 chars", func(t *testing.T) {
				code, title = stamp("")
				for range 65 - len(code) {
					code += "o"
				}
				newdl := &models.DL{
					Code:  code,
					Title: title,
				}

				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update DL title to string with more than 64 chars", func(t *testing.T) {
				code, title = stamp("")
				for range 65 - len(title) {
					title += "o"
				}
				newdl := &models.DL{
					Code:  code,
					Title: title,
				}

				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update DL code to sth repetetive", func(t *testing.T) {
				firstdl := &models.DL{}
				repo.db.Model(&models.DL{}).First(firstdl)

				newdl := &models.DL{
					Code:  firstdl.Code,
					Title: "sth valid",
				}

				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update DL titel to sth repetetive", func(t *testing.T) {
				firstdl := &models.DL{}
				repo.db.Model(&models.DL{}).First(firstdl)

				newdl := &models.DL{
					Code:  "sth valid",
					Title: firstdl.Title,
				}

				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)
			})

			t.Run("should update to a valid state", func(t *testing.T) {
				code, title = stamp("D003-Updated")
				newdl := &models.DL{
					Code:  code,
					Title: title,
				}
				err := repo.Update(dl.ID, newdl)
				assert.Nil(t, err)

				updateddl := &models.DL{}
				repo.db.Model(&models.DL{}).First(updateddl, "id = ?", dl.ID)

				assert.Equal(t, code, updateddl.Code)
				assert.Equal(t, title, updateddl.Title)

			})

			t.Run("should update a DL that is referenced", func(t *testing.T) {
				sl1_id, sl2_id, dl_id, v_id := scenario1(repo.db)

				code, title = stamp("")
				validdl := &models.DL{
					Code:  code,
					Title: title,
				}

				assert.Nil(t, repo.Update(dl_id, validdl))

				dl, _ := repo.Read(dl_id)
				assert.Equal(t, dl.Version, uint(1))

				revert_scenario1(repo.db, sl1_id, sl2_id, dl_id, v_id)

			})

			t.Run("Row Versioning", func(t *testing.T) {
				code, title = stamp("D003-Updated-Twice")
				newdl := &models.DL{
					Code:  code,
					Title: title,
				}
				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)

				newdl.Version = 1
				err = repo.Update(dl.ID, newdl)
				assert.Nil(t, err)

				code, title = stamp("D003-Updated-Triple")
				newdl = &models.DL{
					Code:  code,
					Title: title,
				}
				err = repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)

				newdl.Version = 2
				err = repo.Update(dl.ID, newdl)
				assert.Nil(t, err)

				code, title = stamp("D003-Updated-4th")
				newdl = &models.DL{
					Code:  code,
					Title: title,
				}
				err = repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)

			})

		})

		t.Run("Delete", func(t *testing.T) {

			t.Run("Should Delete the last 4 DL's", func(t *testing.T) {

				DL := &models.DL{}
				if err := repo.db.Model(&models.DL{}).Last(DL).Error; err != nil {
					fmt.Println("DB issue, aborting...", err)
					return
				}
				err := repo.Delete(DL.ID)
				assert.Nil(t, err)
				_, err = repo.Read(DL.ID)
				assert.NotNil(t, err)

				DL = &models.DL{}
				if err := repo.db.Model(&models.DL{}).Last(DL).Error; err != nil {
					fmt.Println("DB issue, aborting...", err)
					return
				}

				err = repo.Delete(DL.ID)
				assert.Nil(t, err)
				_, err = repo.Read(DL.ID)
				assert.NotNil(t, err)

				DL = &models.DL{}
				if err := repo.db.Model(&models.DL{}).Last(DL).Error; err != nil {
					fmt.Println("DB issue, aborting...", err)
					return
				}

				err = repo.Delete(DL.ID)
				assert.Nil(t, err)
				_, err = repo.Read(DL.ID)
				assert.NotNil(t, err)

				DL = &models.DL{}
				if err := repo.db.Model(&models.DL{}).Last(DL).Error; err != nil {
					fmt.Println("DB issue, aborting...", err)
					return
				}

				err = repo.Delete(DL.ID)
				assert.Nil(t, err)
				_, err = repo.Read(DL.ID)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Delete a DL that is referenced by a vi", func(t *testing.T) {
				sl1_id, sl2_id, dl_id, v_id := scenario1(repo.db)

				assert.NotNil(t, repo.Delete(dl_id))

				revert_scenario1(repo.db, sl1_id, sl2_id, dl_id, v_id)
			})
		})
	})
}
